package entity

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID                 uint64               `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt          time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt       `json:"-" gorm:"index"`
	RoomNumber         string               `json:"room_number" gorm:"not null;uniqueIndex"`
	Status             RoomStatus           `json:"status" gorm:"type:varchar(20);not null;default:'available'"`
	RoomCategoryID     uint64               `json:"room_category_id" gorm:"not null;index"`
	RoomCategory       RoomCategory         `json:"room_category" gorm:"foreignKey:RoomCategoryID"`
	Users              []User               `json:"users" gorm:"foreignKey:RoomID"`
	RoomAmenities      []RoomAmenities      `json:"room_amenities" gorm:"foreignKey:RoomID"`
	MaintenanceHistory []MaintenanceHistory `json:"maintenance_histories" gorm:"foreignKey:RoomID"`
}

type ListRooms struct {
	ID             uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	RoomNumber     string         `json:"room_number" gorm:"not null;uniqueIndex"`
	Status         RoomStatus     `json:"status" gorm:"type:varchar(20);not null;default:'available'"`
	RoomCategoryID uint64         `json:"room_category_id" gorm:"not null;index"`
	RoomCategory   RoomCategory   `json:"room_category" gorm:"foreignKey:RoomCategoryID"`
	UserCount      int64          `json:"user_count"`
}

type CreateRoom struct {
	RoomNumber     string     `json:"room_number" binding:"required"`
	Status         RoomStatus `json:"status" binding:"required,oneof=available unavailable maintenance"`
	RoomCategoryID uint64     `json:"room_category_id" binding:"required,numeric"`
	AmenityIDs     []uint64   `json:"amenity_ids" binding:"omitempty,dive,numeric"`
}

type UpdateRoom struct {
	RoomNumber     string     `json:"room_number" binding:"omitempty"`
	Status         RoomStatus `json:"status" binding:"omitempty,oneof=available unavailable maintenance"`
	RoomCategoryID uint64     `json:"room_category_id" binding:"omitempty,numeric"`
	AmenityIDs     []uint64   `json:"amenity_ids" binding:"omitempty,dive,numeric"`
}

func (Room) TableName() string {
	return "rooms"
}

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusUnavailable RoomStatus = "unavailable"
	RoomStatusMaintenance RoomStatus = "maintenance"
)
