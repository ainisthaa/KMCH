package mentalhealthcache

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

// Cache loads mental health patient IDs from Excel ONCE at startup.
// All lookups are O(1) — the file is never re-read after startup.
type Cache struct {
	mu  sync.RWMutex
	ids map[string]struct{}
}

// Mock IDs used when Excel file is unavailable (development / testing).
var mockIDs = map[string]struct{}{
	"1459901255427": {},
	"1600100708962": {},
	"ME4166375915":  {},
}

func NewCache(filePath string) *Cache {
	c := &Cache{ids: make(map[string]struct{})}
	if err := c.loadFromExcel(filePath); err != nil {
		log.Printf("mental_health_cache: cannot load Excel (%v) — using mock data", err)
		c.ids = mockIDs
	} else {
		log.Printf("mental_health_cache: loaded %d records from %s", len(c.ids), filePath)
	}
	return c
}

// HasIssue returns true if the patient is listed in the mental health issues Excel.
func (c *Cache) HasIssue(citizenID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.ids[NormalizeID(citizenID)]
	return ok
}

// HasCompletedScreening returns true if the patient has completed the mental health screening form.
// Currently backed by the same Excel as HasIssue.
// Future: use a separate "completed, no issue" list to support the false case.
func (c *Cache) HasCompletedScreening(citizenID string) bool {
	return c.HasIssue(citizenID)
}

// NeedsPsychologist returns true when the patient must visit the psychologist station.
//
//   - Not completed screening → true  (未筛查 → 需要)
//   - Completed + mental health issue → true
//   - Completed + no issue → false
func (c *Cache) NeedsPsychologist(citizenID string) bool {
	hasCompleted := c.HasCompletedScreening(citizenID)
	hasIssue := c.HasIssue(citizenID)
	return !hasCompleted || hasIssue
}

func (c *Cache) loadFromExcel(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("no sheets found")
	}
	// RawCellValue avoids excelize formatting large numeric IDs into truncated
	// scientific notation, which would lose precision before NormalizeID runs.
	rows, err := f.GetRows(sheets[0], excelize.Options{RawCellValue: true})
	if err != nil {
		return err
	}
	if len(rows) < 2 {
		return nil
	}

	idColIdx := -1
	for i, cell := range rows[0] {
		h := strings.TrimSpace(cell)
		if strings.Contains(h, "เลขบัตรประชาชน") || strings.Contains(h, "National ID") {
			idColIdx = i
			break
		}
	}
	if idColIdx == -1 {
		return fmt.Errorf("citizen ID column not found")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, row := range rows[1:] {
		if idColIdx >= len(row) {
			continue
		}
		if norm := NormalizeID(strings.TrimSpace(row[idColIdx])); norm != "" {
			c.ids[norm] = struct{}{}
		}
	}
	return nil
}

// NormalizeID handles decimal (1600100708962.00), scientific (1.6E+12), and passport IDs.
func NormalizeID(raw string) string {
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
