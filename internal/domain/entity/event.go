package entity

import (
	"time"

	"gorm.io/gorm"
)

type EventType string

const (
	EventTypeDeadline    EventType = "deadline"
	EventTypeMaintenance EventType = "maintenance"
	EventTypeEvent       EventType = "event"
	EventTypeMeeting     EventType = "meeting"
)

type Event struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	EventDate   time.Time      `json:"event_date" gorm:"not null"`
	StartTime   string         `json:"start_time" gorm:"not null"`
	EndTime     string         `json:"end_time" gorm:"not null"`
	Location    string         `json:"location" gorm:"not null"`
	Type        EventType      `json:"type" gorm:"not null"`
	IsMandatory bool           `json:"is_mandatory" gorm:"default:false"`
	Organizer   string         `json:"organizer" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Event) TableName() string {
	return "events"
}
