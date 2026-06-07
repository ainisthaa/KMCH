package domain

import "time"

// PatientQueue is the doctor-consultation queue entry.
// The Queue field is an internal FIFO timestamp — NEVER expose it in any API response.
type PatientQueue struct {
	QueueID    uint       `gorm:"primaryKey;autoIncrement;column:queue_id"`
	LineID     string     `gorm:"column:line_id;not null"`
	EventID    uint       `gorm:"column:event_id"`
	Queue      string     `gorm:"column:queue"` // internal FIFO order — never return to client
	RoomID     *string    `gorm:"column:room_id"`
	Status     string     `gorm:"column:status;type:varchar(20)"`
	Station    string     `gorm:"column:station;type:varchar(20)"`
	QStartTime *time.Time `gorm:"column:q_starttime"`
	QEndTime   *time.Time `gorm:"column:q_endtime"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
	Room       *DoctorRoom `gorm:"foreignKey:RoomID;references:RoomID"`
}

func (PatientQueue) TableName() string { return "patient_queue" }

// Queue status
const (
	QueueWaiting   = "waiting"
	QueueAssigned  = "assigned"
	QueueCompleted = "completed"
	QueueSkip      = "skip"
)

// Queue station
const (
	QStationRegistered = "registered"
	QStationPaid       = "paid"
	QStationQueue      = "queue"
	QStationCompleted  = "completed"
)
