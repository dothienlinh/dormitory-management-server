package entity

type RoomAmenities struct {
	Base
	RoomID    uint64   `json:"room_id"`
	AmenityID uint64   `json:"amenity_id"`
	Room      *Room    `json:"room"`
	Amenity   *Amenity `json:"amenity"`
}

func (RoomAmenities) TableName() string {
	return "room_amenities"
}

type CreateRoomAmenities struct {
	RoomID    uint64 `json:"room_id" binding:"required,numeric"`
	AmenityID uint64 `json:"amenity_id" binding:"required,numeric"`
}
