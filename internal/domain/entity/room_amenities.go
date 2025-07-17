package entity

import (
	"time"

	"gorm.io/gorm"
)

type RoomAmenities struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	RoomID    uint64         `json:"room_id" gorm:"not null"`
	AmenityID uint64         `json:"amenity_id" gorm:"not null"`
	Room      *Room          `json:"room" gorm:"foreignKey:RoomID"`
	Amenity   *Amenity       `json:"amenity" gorm:"foreignKey:AmenityID"`
}

func (RoomAmenities) TableName() string {
	return "room_amenities"
}

type CreateRoomAmenities struct {
	RoomID    uint64 `json:"room_id" binding:"required,numeric"`
	AmenityID uint64 `json:"amenity_id" binding:"required,numeric"`
}
