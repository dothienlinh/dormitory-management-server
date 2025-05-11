package models

import "time"

type RoomRent struct {
	BaseModel
	RoomID uint           `json:"room_id" gorm:"not null"`
	Room   *Room          `json:"room"`
	UserID uint           `json:"user_id" gorm:"not null"`
	User   *User          `json:"user"`
	Status RoomRentStatus `json:"status" gorm:"type:room_rent_status;not null"`
}

type RoomRentStatus string

const (
	RoomRentStatusActive    RoomRentStatus = "active"
	RoomRentStatusInactive  RoomRentStatus = "inactive"
	RoomRentStatusGraduated RoomRentStatus = "graduated"
)

type CreateRoomRent struct {
	RoomID uint           `json:"room_id" binding:"required"`
	UserID uint           `json:"user_id" binding:"required"`
	Status RoomRentStatus `json:"status" binding:"required,oneof=active inactive graduated"`
}

type UpdateRoomRent struct {
	Status  RoomRentStatus `json:"status" binding:"required,oneof=active inactive graduated"`
	EndDate *time.Time     `json:"end_date" binding:"required"`
}
