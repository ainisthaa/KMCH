package domain

import "time"

// PatientCheck tracks one patient's registration for one event.
// UNIQUE INDEX on (line_id, event_id) — one record per patient per event.
type PatientCheck struct {
	CheckID           uint       `gorm:"primaryKey;autoIncrement;column:check_id"`
	LineID            string     `gorm:"column:line_id;not null;uniqueIndex:idx_line_event"`
	EventID           uint       `gorm:"column:event_id;uniqueIndex:idx_line_event"`
	PsyevalForm       bool       `gorm:"column:psyeval_form;default:false"`       // completed mental-health form
	IsSV              bool       `gorm:"column:is_sv;default:false"`              // has mental-health issue/risk
	IsPaid            bool       `gorm:"column:is_paid;default:false"`            // payment done
	NeedsTransfer     bool       `gorm:"column:needs_transfer;default:false"`     // needs rights-transfer (set at payment scan)
	TransferCompleted bool       `gorm:"column:transfer_completed;default:false"` // rights-transfer actually done
	NeedsPsychologist bool       `gorm:"column:needs_psychologist;default:false"` // computed at payment scan
	PsychologistDone  bool       `gorm:"column:psychologist_done;default:false"`  // psychologist visit done
	RouteType         string     `gorm:"column:route_type;type:varchar(1)"`       // A/B/C/D
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
}

func (PatientCheck) TableName() string { return "patient_check" }

// Route type constants
const (
	RouteA = "A" // no psychologist, no transfer
	RouteB = "B" // no psychologist, transfer
	RouteC = "C" // psychologist, no transfer
	RouteD = "D" // psychologist + transfer
)

// Station codes
const (
	StationRegistration      = "registration"
	StationPayment           = "payment"
	StationPsychologist      = "psychologist"
	StationRightsTransfer    = "rights_transfer"
	StationDoctorConsultation = "doctor_consultation"
	StationXray              = "xray"
)
