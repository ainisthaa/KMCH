package domain

import "time"

type EventInfo struct {
	EventID       uint      `gorm:"primaryKey;autoIncrement;column:event_id"`
	EventName     string    `gorm:"column:event_name"`
	EventDateFrom time.Time `gorm:"column:event_date_from"`
	EventDateTo   time.Time `gorm:"column:event_date_to"`
}

func (EventInfo) TableName() string { return "event_info" }
