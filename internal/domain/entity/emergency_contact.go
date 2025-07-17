package entity

import (
	"time"

	"gorm.io/gorm"
)

type EmergencyContact struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	UserID       uint64         `json:"user_id" gorm:"not null;index"`
	Name         string         `json:"name" gorm:"not null"`
	Phone        string         `json:"phone" gorm:"not null"`
	Relationship string         `json:"relationship" gorm:"not null"`
}

func (EmergencyContact) TableName() string {
	return "emergency_contacts"
}
