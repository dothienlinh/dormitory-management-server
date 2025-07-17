package entity

import (
	"time"

	"gorm.io/gorm"
)

type Amenity struct {
	ID            uint64          `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `json:"-" gorm:"index"`
	Name          string          `json:"name" gorm:"not null"`
	RoomAmenities []RoomAmenities `json:"room_amenities" gorm:"foreignKey:AmenityID"`
}

func (Amenity) TableName() string {
	return "amenities"
}

type CreateAmenity struct {
	Name string `json:"name" binding:"required"`
}

type UpdateAmenity struct {
	Name string `json:"name" binding:"required"`
}
