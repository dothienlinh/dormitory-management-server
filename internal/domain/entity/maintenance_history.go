package entity

import (
	"time"

	"gorm.io/gorm"
)

type MaintenanceHistory struct {
	ID              uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	Description     string         `json:"description" gorm:"type:text;not null"`
	RoomID          uint64         `json:"room_id" gorm:"not null;index"`
	MaintenanceDate time.Time      `json:"maintenance_date" gorm:"not null"`
	Cost            float64        `json:"cost" gorm:"not null"`
}

func (MaintenanceHistory) TableName() string {
	return "maintenance_histories"
}

type CreateMaintenanceHistory struct {
	Description     string    `json:"description" binding:"required"`
	RoomID          uint64    `json:"room_id" binding:"required,numeric"`
	MaintenanceDate time.Time `json:"maintenance_date" binding:"required,validdate"`
	Cost            float64   `json:"cost" binding:"required,numeric"`
}
