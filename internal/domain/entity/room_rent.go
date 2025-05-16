package entity

// RoomRentStatus represents the status of a room rental
type RoomRentStatus string

const (
	RoomRentStatusActive   RoomRentStatus = "active"
	RoomRentStatusInactive RoomRentStatus = "inactive"
)

// RoomRent entity represents the relationship between users and rooms
type RoomRent struct {
	Base
	RoomID uint           `json:"room_id"`
	Room   Room           `json:"room"`
	UserID uint           `json:"user_id"`
	User   User           `json:"user"`
	Status RoomRentStatus `json:"status"`
}

func (RoomRent) TableName() string {
	return "room_rents"
}

// CreateRoomRent is the data transfer object for creating a room rental
type CreateRoomRent struct {
	RoomID uint           `json:"room_id" binding:"required"`
	UserID uint           `json:"user_id" binding:"required"`
	Status RoomRentStatus `json:"status" binding:"omitempty,oneof=active inactive"`
}

// RoomRentDTO is a data transfer object for RoomRent entity
type RoomRentDTO struct {
	ID     uint           `json:"id"`
	RoomID uint           `json:"room_id"`
	Room   Room           `json:"room,omitempty"`
	UserID uint           `json:"user_id"`
	User   User           `json:"user,omitempty"`
	Status RoomRentStatus `json:"status"`
}

type RemoveUserFromRoom struct {
	UserID uint `json:"user_id" binding:"required"`
	RoomID uint `json:"room_id" binding:"required"`
}
