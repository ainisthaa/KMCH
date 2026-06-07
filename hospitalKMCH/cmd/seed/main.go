package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"lineoa-miniapp/domain"
)

// ── Room ID constants ─────────────────────────────────────────────────────────

const (
	room1ID = "room-001"
	room2ID = "room-002"
	room3ID = "room-003"
	room4ID = "room-004"
	room5ID = "room-005"
)

var seedTime = time.Now()
var queueSeq int

// ── Helpers ───────────────────────────────────────────────────────────────────

func nextQueue() string {
	queueSeq++
	return fmt.Sprintf("20260607%06d", queueSeq)
}

// lineIDFrom generates a deterministic UUID v5 from the patient ID string.
func lineIDFrom(id string) string {
	nsBytes, _ := hex.DecodeString("6ba7b8109dad11d180b400c04fd430c8") // DNS namespace
	h := sha1.New()
	h.Write(nsBytes)
	h.Write([]byte("kmch-seed-" + id))
	b := h.Sum(nil)
	b[6] = (b[6] & 0x0f) | 0x50 // version 5
	b[8] = (b[8] & 0x3f) | 0x80 // variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// normalizeID matches production logic in mentalhealthcache.NormalizeID.
func normalizeID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	cleaned := strings.ReplaceAll(strings.ReplaceAll(raw, " ", ""), "-", "")
	if f, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return fmt.Sprintf("%d", int64(math.Round(f)))
	}
	return strings.ToUpper(cleaned)
}

func allDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func mustDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic("bad date: " + s)
	}
	return t
}

func telNo(nationalID, passportID string) string {
	id := nationalID
	if id == "" {
		id = passportID
	}
	digits := ""
	for _, r := range id {
		if r >= '0' && r <= '9' {
			digits += string(r)
		}
	}
	for len(digits) < 8 {
		digits = "0" + digits
	}
	if len(digits) >= 8 {
		return "08" + digits[len(digits)-8:]
	}
	return "0800000000"
}

func ptr[T any](v T) *T { return &v }

// ── Excel eligibility map ─────────────────────────────────────────────────────

type excelRow struct {
	PrefixName  string
	FirstName   string
	LastName    string
	IDPassport  string
	AppDate     string
}

func readExcel(path string) (map[string]excelRow, error) {
	result := make(map[string]excelRow)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return result, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return result, fmt.Errorf("no sheets found in %s", path)
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil || len(rows) < 2 {
		return result, err
	}

	// Locate columns by header name (case-insensitive substring match).
	colIdx := map[string]int{
		"prefix":   -1,
		"first":    -1,
		"last":     -1,
		"id":       -1,
		"appdate":  -1,
	}
	for i, cell := range rows[0] {
		h := strings.ToLower(strings.TrimSpace(cell))
		switch {
		case strings.Contains(h, "prefix") || strings.Contains(h, "คำนำหน้า"):
			colIdx["prefix"] = i
		case strings.Contains(h, "first") || strings.Contains(h, "ชื่อ"):
			colIdx["first"] = i
		case strings.Contains(h, "last") || strings.Contains(h, "นามสกุล"):
			colIdx["last"] = i
		case strings.Contains(h, "id_passport") || strings.Contains(h, "passport") || strings.Contains(h, "เลขบัตร"):
			colIdx["id"] = i
		case strings.Contains(h, "appointment_date") || strings.Contains(h, "วันนัด"):
			colIdx["appdate"] = i
		}
	}

	idCol := colIdx["id"]
	if idCol == -1 {
		return result, fmt.Errorf("id_passport column not found in %s", path)
	}

	for _, row := range rows[1:] {
		get := func(col int) string {
			if col >= 0 && col < len(row) {
				return strings.TrimSpace(row[col])
			}
			return ""
		}
		rawID := get(idCol)
		if rawID == "" {
			continue
		}
		norm := normalizeID(rawID)
		result[norm] = excelRow{
			PrefixName: get(colIdx["prefix"]),
			FirstName:  get(colIdx["first"]),
			LastName:   get(colIdx["last"]),
			IDPassport: rawID,
			AppDate:    get(colIdx["appdate"]),
		}
	}
	return result, nil
}

// ── Seed patient definitions ──────────────────────────────────────────────────

type seedDef struct {
	Idx        int
	FirstName  string
	LastName   string
	NationalID string
	PassportID string
	Date       time.Time

	// patient_check
	PsyevalForm      bool
	IsSV             bool
	IsPaid           bool
	NeedsTransfer    bool
	NeedsPsycho      bool
	PsychoDone       bool
	TransferDone     bool
	RouteType        string

	// patient_queue (HasQueue=false → no queue record)
	HasQueue     bool
	QueueStatus  string
	QueueStation string
	RoomID       string // empty = no room

	TestCase string
}

func (p seedDef) rawID() string {
	if p.NationalID != "" {
		return p.NationalID
	}
	return p.PassportID
}

func (p seedDef) lineID() string { return lineIDFrom(p.rawID()) }

func (p seedDef) idType() string {
	if p.PassportID != "" {
		return "passport"
	}
	return "national"
}

func (p seedDef) routeLabel() string {
	if p.RouteType == "" {
		return "-"
	}
	return p.RouteType
}

func (p seedDef) queueLabel() string {
	switch {
	case !p.HasQueue:
		return "-"
	case p.QueueStatus == domain.QueueAssigned && p.RoomID != "":
		return "ASSIGNED→R" + roomNum(p.RoomID)
	default:
		return strings.ToUpper(p.QueueStatus)
	}
}

func roomNum(roomID string) string {
	switch roomID {
	case room1ID:
		return "1"
	case room2ID:
		return "2"
	case room3ID:
		return "3"
	case room4ID:
		return "4"
	case room5ID:
		return "5"
	}
	return "?"
}

// patients is the authoritative list of 50 seed patients.
//
// Route B WAITING patients (9-15): transferDone=true — they completed rights
// transfer and are now waiting for doctor consultation.
// Route C patients 16-18: psychoDone=true, COMPLETED — already through doctor.
// Route C patients 19-23: psychoDone=false — blocked at psychologist station.
var patients = []seedDef{
	// ── Case 1: Route A — no psychologist, no transfer, WAITING (8) ─────────
	{1, "พัชรพล", "พันธะไชย", "1459901255427", "", mustDate("2026-06-19"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A / dup scan-after-payment"},
	{2, "สรวรรณ", "อยู่ทอง", "1849901970599", "", mustDate("2026-06-22"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},
	{3, "ชัญญานุช", "นุชประมูล", "1209702282141", "", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},
	{4, "ปนัฎฎา", "คงแก้ว", "1102400200871", "", mustDate("2026-06-29"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},
	{5, "อาภาพร", "เหมวงษ์", "1229901232968", "", mustDate("2026-06-19"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},
	{6, "เพชรลดา", "สิทธิวัง", "1104301058751", "", mustDate("2026-06-29"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},
	{7, "กนกวรรณ", "ดอดกระโทก", "1560101649766", "", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},
	{8, "ณหทัย", "พนาพฤกษกุล", "1102400219858", "", mustDate("2026-06-22"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route A"},

	// ── Case 2: Route B — transfer required, WAITING (7) ────────────────────
	// transferDone=true: rights transfer completed before joining doctor queue
	{9, "นภสร", "พลเยี่ยม", "1103100926410", "", mustDate("2026-06-29"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B / dup scan-doctor-queue"},
	{10, "ทินภัทร", "แซ่ลี้", "1209000339827", "", mustDate("2026-06-25"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B"},
	{11, "พนธกร", "ศิริพรรค", "1719900755997", "", mustDate("2026-06-16"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B"},
	{12, "พิมพ์มาดา", "ศรีใจ", "1549900747708", "", mustDate("2026-06-23"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B"},
	{13, "นวมินทร์", "รสจันทร์", "1349901525514", "", mustDate("2026-06-24"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B"},
	{14, "นนทนัตถ์", "ดอกป่าน", "1749901191744", "", mustDate("2026-06-18"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B"},
	{15, "ธีรดา", "มาซา", "1103200138022", "", mustDate("2026-06-23"),
		true, false, true, true, false, false, true, domain.RouteB,
		true, domain.QueueWaiting, domain.QStationQueue, "",
		"Route B"},

	// ── Case 3: Route C — psychologist required ───────────────────────────────
	// Patients 16-18: psychologist done → joined + COMPLETED doctor consultation
	{16, "กฤษฏ์", "ฉิมเลี้ยง", "1250101705383", "", mustDate("2026-06-17"),
		false, false, true, false, true, true, false, domain.RouteC,
		true, domain.QueueCompleted, domain.QStationCompleted, room3ID,
		"Route C / psych done / COMPLETED"},
	{17, "กรนันท์", "เกิดสิน", "1100703947050", "", mustDate("2026-06-15"),
		false, false, true, false, true, true, false, domain.RouteC,
		true, domain.QueueCompleted, domain.QStationCompleted, room4ID,
		"Route C / psych done / COMPLETED"},
	{18, "ญาณิศา", "สีขาว", "1103000201031", "", mustDate("2026-06-19"),
		false, false, true, false, true, true, false, domain.RouteC,
		true, domain.QueueCompleted, domain.QStationCompleted, room5ID,
		"Route C / psych done / COMPLETED"},
	// Patients 19-23: psychologist NOT done → blocked, no queue entry
	{19, "สุวพัชร", "โชติพงศ์พุฒิ", "1729800371810", "", mustDate("2026-06-24"),
		false, false, true, false, true, false, false, domain.RouteC,
		false, "", "", "",
		"Route C / psych blocked"},
	{20, "ธัญวลัย", "ภักดีพันธ์", "1709901705977", "", mustDate("2026-06-15"),
		false, false, true, false, true, false, false, domain.RouteC,
		false, "", "", "",
		"Route C / psych blocked"},
	{21, "พราวพิมล", "เฉิดศรีธนิตย์", "1103704307898", "", mustDate("2026-06-19"),
		false, false, true, false, true, false, false, domain.RouteC,
		false, "", "", "",
		"Route C / psych blocked"},
	{22, "ชวิศ", "ถนัดธนูศิลป์", "1101402376828", "", mustDate("2026-06-22"),
		false, false, true, false, true, false, false, domain.RouteC,
		false, "", "", "",
		"Route C / psych blocked"},
	{23, "ณัชพล", "ฐานานุกรม", "1909803268686", "", mustDate("2026-06-15"),
		false, false, true, false, true, false, false, domain.RouteC,
		false, "", "", "",
		"Route C / psych blocked"},

	// ── Case 4: Route D — psychologist + transfer required, blocked (7) ──────
	// isSV=true triggers needsPsycho. Psychologist first in Route D, so blocked.
	{24, "ปรีญามาตย์", "โสปัญหริ", "1101501341950", "", mustDate("2026-06-15"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},
	{25, "ชนัญชิดา", "เหลาธรรม", "1103704272601", "", mustDate("2026-06-24"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},
	{26, "Thanphitcha", "Saosiri", "1339600162501", "", mustDate("2026-06-19"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},
	{27, "ฐานิกา", "แก้วแสงเรือง", "1909803195875", "", mustDate("2026-06-25"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},
	{28, "สุชานันท์", "อัตถารม", "1149600172720", "", mustDate("2026-06-16"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},
	{29, "ธนภูมิ", "ชาวกระเดียน", "1104700177583", "", mustDate("2026-06-15"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},
	{30, "นรบดี", "สีขน", "1103400150620", "", mustDate("2026-06-15"),
		true, true, true, true, true, false, false, domain.RouteD,
		false, "", "", "",
		"Route D / psych blocked"},

	// ── Case 5: Unpaid — cannot join queue (3) ────────────────────────────────
	{31, "ภัทราพร", "จิบสมานบุญ", "1729800379519", "", mustDate("2026-06-24"),
		true, false, false, false, false, false, false, "",
		false, "", "", "",
		"Unpaid"},
	{32, "ภควดี", "คุ้มตะสิน", "1669900602245", "", mustDate("2026-06-19"),
		true, false, false, false, false, false, false, "",
		false, "", "", "",
		"Unpaid"},
	{33, "ธนดล", "อ้นสิงห์มา", "1103100942806", "", mustDate("2026-06-15"),
		true, false, false, false, false, false, false, "",
		false, "", "", "",
		"Unpaid"},

	// ── Case 6: Route A ASSIGNED to rooms (12 Thai patients) ─────────────────
	// Room 1: 34, 39, 44   Room 2: 35, 40, 45   Room 3: 36, 41
	// Room 4: 37, 42       Room 5: 38, 43
	{34, "สุขพิชย์ฏา", "แจ่มจันทร์", "1309903472264", "", mustDate("2026-06-19"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room1ID,
		"Assigned→Room1"},
	{35, "นรีลักษณ์", "อนุชาติชัยกุล", "1100401402468", "", mustDate("2026-06-25"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room2ID,
		"Assigned→Room2"},
	{36, "สุดพิเศษ", "วงศ์ฝั้น", "1510101438986", "", mustDate("2026-06-19"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room3ID,
		"Assigned→Room3"},
	{37, "ศศิธร", "คำรัง", "1439600063642", "", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room4ID,
		"Assigned→Room4"},
	{38, "Ittiwat", "Chonpatatip", "1101700445135", "", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room5ID,
		"Assigned→Room5"},
	{39, "คุณากร", "บุตรครุธ", "1819900531303", "", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room1ID,
		"Assigned→Room1"},
	{40, "ศักรนันทน์", "ลือสวัสดิ์", "1909803177826", "", mustDate("2026-06-22"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room2ID,
		"Assigned→Room2"},
	{41, "ธนธร", "นุ่นมะลัง", "1302201138830", "", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room3ID,
		"Assigned→Room3"},
	{42, "สุภาพร", "สง่าแสง", "1729800373570", "", mustDate("2026-06-23"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room4ID,
		"Assigned→Room4"},
	{43, "ศศิกานต์", "ดวงดี", "1279900342476", "", mustDate("2026-06-19"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room5ID,
		"Assigned→Room5"},
	{44, "ธัญชนก", "ปิ่นทอง", "1102200256202", "", mustDate("2026-06-22"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room1ID,
		"Assigned→Room1"},
	{45, "ภูกิจ", "เกตุสวาสดิ์", "1209601466769", "", mustDate("2026-06-25"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room2ID,
		"Assigned→Room2"},

	// ── Case 7: Passport patients — force is_transfer=false (5) ─────────────
	// 46,47 → ASSIGNED  48,49 → COMPLETED  50 → SKIP
	{46, "Angel Anne", "Astillero", "", "P1506276D", mustDate("2026-06-24"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room1ID,
		"Passport / ASSIGNED→R1"},
	{47, "VUOCH SIM", "LIM", "", "N01669092", mustDate("2026-06-17"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueAssigned, domain.QStationQueue, room2ID,
		"Passport / ASSIGNED→R2"},
	{48, "ปรียนันท์", "เต็มฤทธิกุลชัย", "", "AD3221766", mustDate("2026-06-29"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueCompleted, domain.QStationCompleted, room3ID,
		"Passport / COMPLETED / dup complete"},
	{49, "Jaeyoung", "Ko", "", "M890T7108", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueCompleted, domain.QStationCompleted, room4ID,
		"Passport / COMPLETED"},
	{50, "สาริศา", "เล้าเจริญ", "", "171880010920", mustDate("2026-06-15"),
		true, false, true, false, false, false, false, domain.RouteA,
		true, domain.QueueSkip, domain.QStationQueue, "",
		"Passport / SKIP"},
}

// ── DB connection ─────────────────────────────────────────────────────────────

func connectDB() *gorm.DB {
	_ = godotenv.Load()
	host := env("SERVICE_DB_HOST", "localhost")
	port := env("SERVICE_DB_PORT", "3306")
	user := env("SERVICE_DB_USER", "appuser")
	pass := env("SERVICE_DB_PASS", "password")
	name := env("SERVICE_DB_NAME", "lineoa_miniapp")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormlogger.Default.LogMode(gormlogger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	return db
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	log.Println("KMCH seed — starting")

	db := connectDB()

	// ── 1. Event ──────────────────────────────────────────────────────────────
	seedEvent(db)

	// ── 2. Doctor Rooms ───────────────────────────────────────────────────────
	seedRooms(db)

	// ── 3. Excel eligibility map ──────────────────────────────────────────────
	excelPath := env("EXCEL_FILE_PATH", "./data/registrations_export_20260529_165557.xlsx")
	excelMap, err := readExcel(excelPath)
	if err != nil {
		log.Printf("WARN: cannot read Excel (%v) — skipping Excel validation", err)
	} else {
		log.Printf("Excel loaded: %d eligibility records", len(excelMap))
	}

	// ── 4. Patients ───────────────────────────────────────────────────────────
	var patientCount, checkCount, queueCount int
	for _, p := range patients {
		// Validate against Excel.
		normID := normalizeID(p.rawID())
		if _, found := excelMap[normID]; !found && len(excelMap) > 0 {
			log.Printf("WARN: patient #%d %s %s (id=%s) not found in Excel — seeding anyway",
				p.Idx, p.FirstName, p.LastName, p.rawID())
		}

		// patient_info
		regDate := p.Date
		patient := domain.PatientInfo{
			LineID:       p.lineID(),
			FirstName:    p.FirstName,
			LastName:     p.LastName,
			TelNo:        telNo(p.NationalID, p.PassportID),
			NationalID:   p.NationalID,
			PassportID:   p.PassportID,
			RegisterDate: &regDate,
		}
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "line_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "tel_no", "national_id", "passport_id", "register_date"}),
		}).Create(&patient).Error; err != nil {
			log.Fatalf("upsert patient #%d: %v", p.Idx, err)
		}
		patientCount++

		// patient_check
		check := domain.PatientCheck{
			LineID:            p.lineID(),
			EventID:           1,
			PsyevalForm:       p.PsyevalForm,
			IsSV:              p.IsSV,
			IsPaid:            p.IsPaid,
			NeedsTransfer:     p.NeedsTransfer,
			TransferCompleted: p.TransferDone,
			NeedsPsychologist: p.NeedsPsycho,
			PsychologistDone:  p.PsychoDone,
			RouteType:         p.RouteType,
		}
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "line_id"}, {Name: "event_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"psyeval_form", "is_sv", "is_paid", "needs_transfer",
				"transfer_completed", "needs_psychologist", "psychologist_done", "route_type",
			}),
		}).Create(&check).Error; err != nil {
			log.Fatalf("upsert check #%d: %v", p.Idx, err)
		}
		checkCount++

		// patient_queue
		if p.HasQueue {
			if ensureQueueEntry(db, p) {
				queueCount++
			}
		}
	}

	// ── 5. Update room timestamps ─────────────────────────────────────────────
	updateRoomTimestamps(db)

	// ── 6. Summary ────────────────────────────────────────────────────────────
	printSummary(db, patientCount, checkCount, queueCount)
	printTable()
}

// ── Seed helpers ──────────────────────────────────────────────────────────────

func seedEvent(db *gorm.DB) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	db.Exec(`
		INSERT INTO event_info (event_id, event_name, event_date_from, event_date_to)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE event_name=VALUES(event_name), event_date_to=VALUES(event_date_to)`,
		1, "KMCH Health Check Test Event", now, tomorrow)
	log.Println("event_info seeded")
}

func seedRooms(db *gorm.DB) {
	rooms := []domain.DoctorRoom{
		{RoomID: room1ID, RoomName: "Room 1"},
		{RoomID: room2ID, RoomName: "Room 2"},
		{RoomID: room3ID, RoomName: "Room 3"},
		{RoomID: room4ID, RoomName: "Room 4"},
		{RoomID: room5ID, RoomName: "Room 5"},
	}
	for _, r := range rooms {
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"room_name"}),
		}).Create(&r)
	}
	log.Println("doctor_room seeded (5 rooms)")
}

// ensureQueueEntry inserts a queue record if one doesn't exist yet for (line_id, event_id).
// Returns true if a new record was inserted.
func ensureQueueEntry(db *gorm.DB, p seedDef) bool {
	var existing domain.PatientQueue
	err := db.Where("line_id = ? AND event_id = ?", p.lineID(), 1).First(&existing).Error
	if err == nil {
		return false // already exists — idempotent
	}

	now := time.Now()
	q := &domain.PatientQueue{
		LineID:  p.lineID(),
		EventID: 1,
		Queue:   nextQueue(),
		Status:  p.QueueStatus,
		Station: p.QueueStation,
	}

	switch p.QueueStatus {
	case domain.QueueWaiting:
		start := now.Add(-20 * time.Minute)
		q.QStartTime = &start

	case domain.QueueAssigned:
		start := now.Add(-60 * time.Minute)
		q.QStartTime = &start
		if p.RoomID != "" {
			q.RoomID = ptr(p.RoomID)
		}

	case domain.QueueCompleted:
		var start, end time.Time
		// Passport completed patients joined later; Thai completed joined earlier.
		if p.PassportID != "" {
			start = now.Add(-90 * time.Minute)
			end = now.Add(-30 * time.Minute)
		} else {
			start = now.Add(-120 * time.Minute)
			end = now.Add(-60 * time.Minute)
		}
		q.QStartTime = &start
		q.QEndTime = &end
		if p.RoomID != "" {
			q.RoomID = ptr(p.RoomID)
		}

	case domain.QueueSkip:
		start := now.Add(-45 * time.Minute)
		end := now.Add(-10 * time.Minute)
		q.QStartTime = &start
		q.QEndTime = &end
	}

	if err := db.Create(q).Error; err != nil {
		log.Fatalf("create queue for patient #%d: %v", p.Idx, err)
	}
	return true
}

// updateRoomTimestamps sets last_assigned_at / last_completed_at based on seeded queue data.
func updateRoomTimestamps(db *gorm.DB) {
	now := time.Now()
	// Rooms with COMPLETED patients have last_completed_at set.
	completedRooms := map[string]time.Time{
		room3ID: now.Add(-60 * time.Minute),
		room4ID: now.Add(-60 * time.Minute),
		room5ID: now.Add(-60 * time.Minute),
	}
	for roomID, t := range completedRooms {
		ts := t
		db.Model(&domain.DoctorRoom{}).Where("room_id = ?", roomID).
			Update("last_completed_at", ts)
	}
	// All rooms that have ASSIGNED patients have last_assigned_at set,
	// ordered by round-robin assignment time (most recent first).
	assignedRooms := map[string]time.Time{
		room1ID: now.Add(-40 * time.Minute), // got 44 and 46 most recently
		room2ID: now.Add(-42 * time.Minute),
		room3ID: now.Add(-70 * time.Minute),
		room4ID: now.Add(-72 * time.Minute),
		room5ID: now.Add(-74 * time.Minute),
	}
	for roomID, t := range assignedRooms {
		ts := t
		db.Model(&domain.DoctorRoom{}).Where("room_id = ?", roomID).
			Update("last_assigned_at", ts)
	}
}

// ── Output ────────────────────────────────────────────────────────────────────

func printSummary(db *gorm.DB, patientCount, checkCount, queueCount int) {
	var eventCount, roomCount int64
	db.Model(&domain.EventInfo{}).Count(&eventCount)
	db.Model(&domain.DoctorRoom{}).Count(&roomCount)

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("KMCH Seed Complete")
	fmt.Println("========================================")
	fmt.Printf("Patients:      %d\n", patientCount)
	fmt.Printf("patient_check: %d\n", checkCount)
	fmt.Printf("patient_queue: %d\n", queueCount)
	fmt.Printf("event_info:    %d\n", eventCount)
	fmt.Printf("doctor_rooms:  %d\n", roomCount)
	fmt.Println("========================================")
}

func printTable() {
	fmt.Println()
	fmt.Printf("%-3s  %-36s  %-26s  %-8s  %-5s  %-14s  %s\n",
		"#", "line_id", "name", "id_type", "route", "queue_status", "test_case")
	fmt.Println(strings.Repeat("-", 130))
	for _, p := range patients {
		name := p.FirstName + " " + p.LastName
		if len(name) > 26 {
			name = name[:24] + ".."
		}
		fmt.Printf("%-3d  %-36s  %-26s  %-8s  %-5s  %-14s  %s\n",
			p.Idx,
			p.lineID(),
			name,
			p.idType(),
			p.routeLabel(),
			p.queueLabel(),
			p.TestCase,
		)
	}
	fmt.Println()
}
