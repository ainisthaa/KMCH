package domain

import "time"

type DoctorRoom struct {
	RoomID          string     `gorm:"primaryKey;column:room_id;type:varchar(50)"`
	RoomName        string     `gorm:"column:room_name"`
	LastAssignedAt  *time.Time `gorm:"column:last_assigned_at"`  // round-robin ordering
	LastCompletedAt *time.Time `gorm:"column:last_completed_at"` // least-loaded tie-break
}

func (DoctorRoom) TableName() string { return "doctor_room" }
