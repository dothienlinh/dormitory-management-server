package entity

type Room struct {
	Base
	RoomNumber         string               `json:"room_number"`
	Status             RoomStatus           `json:"status"`
	RoomCategoryID     uint64               `json:"room_category_id"`
	RoomCategory       RoomCategory         `json:"room_category"`
	Users              []User               `json:"users"`
	RoomAmenities      []RoomAmenities      `json:"room_amenities"`
	MaintenanceHistory []MaintenanceHistory `json:"maintenance_histories"`
}

type ListRooms struct {
	Base
	RoomNumber     string       `json:"room_number"`
	Status         RoomStatus   `json:"status"`
	RoomCategoryID uint64       `json:"room_category_id"`
	RoomCategory   RoomCategory `json:"room_category"`
	UserCount      int64        `json:"user_count"`
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
