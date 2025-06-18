package entity

type Room struct {
	Base
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Status         RoomStatus   `json:"status"`
	RoomCategoryID uint64       `json:"room_category_id"`
	RoomCategory   RoomCategory `json:"room_category"`
	Users          []User       `json:"users,omitempty"`
}

type CreateRoom struct {
	Name           string     `json:"name" binding:"required"`
	Description    string     `json:"description" binding:"omitempty"`
	Status         RoomStatus `json:"status" binding:"required,oneof=available unavailable maintenance"`
	RoomCategoryID uint64     `json:"room_category_id" binding:"required,numeric"`
}

type UpdateRoom struct {
	Name           string     `json:"name" binding:"omitempty"`
	Description    string     `json:"description" binding:"omitempty"`
	Status         RoomStatus `json:"status" binding:"omitempty,oneof=available unavailable maintenance"`
	RoomCategoryID uint64     `json:"room_category_id" binding:"omitempty,numeric"`
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
