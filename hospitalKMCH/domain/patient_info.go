package domain

import "time"

type PatientInfo struct {
	LineID       string     `gorm:"primaryKey;column:line_id;type:varchar(100)"`
	FirstName    string     `gorm:"column:first_name"`
	LastName     string     `gorm:"column:last_name"`
	TelNo        string     `gorm:"column:tel_no"`
	NationalID   string     `gorm:"column:national_id"`
	PassportID   string     `gorm:"column:passport_id"`
	RegisterDate *time.Time `gorm:"column:register_date"`
	StudentID    string     `gorm:"column:student_id"`
	EmployeeID   string     `gorm:"column:employee_id"`
}

func (PatientInfo) TableName() string { return "patient_info" }
