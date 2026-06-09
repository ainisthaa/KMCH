package excel

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
	"lineoa-miniapp/pkg/mentalhealthcache"
)

// RegistrationRepository checks whether a patient ID is in the pre-registration Excel.
type RegistrationRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
}

type ExcelRegistrationRepository struct {
	filePath string
}

func NewExcelRegistrationRepository(filePath string) *ExcelRegistrationRepository {
	return &ExcelRegistrationRepository{filePath: filePath}
}

func (r *ExcelRegistrationRepository) Exists(_ context.Context, id string) (bool, error) {
	f, err := excelize.OpenFile(r.filePath)
	if err != nil {
		log.Printf("excel_registration: cannot open file (%v) — skipping verification", err)
		return true, nil // fail-open: allow registration when file unavailable
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return false, fmt.Errorf("registration excel has no sheets")
	}
	rows, err := f.GetRows(sheets[0], excelize.Options{RawCellValue: true})
	if err != nil {
		return false, err
	}
	if len(rows) < 2 {
		return false, nil
	}

	idCol := -1
	for i, cell := range rows[0] {
		h := strings.TrimSpace(cell)
		if h == "id_passport" || strings.Contains(h, "เลขบัตร") || strings.Contains(h, "National ID") {
			idCol = i
			break
		}
	}
	if idCol == -1 {
		log.Printf("excel_registration: ID column not found — skipping verification")
		return true, nil // fail-open
	}

	norm := mentalhealthcache.NormalizeID(id)
	for _, row := range rows[1:] {
		if idCol >= len(row) {
			continue
		}
		if mentalhealthcache.NormalizeID(strings.TrimSpace(row[idCol])) == norm {
			return true, nil
		}
	}
	return false, nil
}
