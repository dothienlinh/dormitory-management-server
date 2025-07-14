package entity

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypePayment      NotificationType = "payment"
	NotificationTypeMaintenance  NotificationType = "maintenance"
	NotificationTypeEvent        NotificationType = "event"
	NotificationTypeService      NotificationType = "service"
	NotificationTypeAnnouncement NotificationType = "announcement"
)

type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityMedium NotificationPriority = "medium"
	NotificationPriorityHigh   NotificationPriority = "high"
)

type Notification struct {
	ID        uint64               `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string               `json:"title" gorm:"not null"`
	Content   string               `json:"content" gorm:"type:text"`
	Type      NotificationType     `json:"type" gorm:"not null"`
	Priority  NotificationPriority `json:"priority" gorm:"default:'medium'"`
	IsRead    bool                 `json:"is_read" gorm:"default:false"`
	UserID    uint64               `json:"user_id" gorm:"not null"`
	User      *User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ReadAt    *time.Time           `json:"read_at,omitempty"`
	CreatedAt time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt       `json:"-" gorm:"index"`
}

func (Notification) TableName() string {
	return "notifications"
}
