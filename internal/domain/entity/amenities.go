package entity

type Amenity struct {
	Base
	Name          string          `json:"name"`
	RoomAmenities []RoomAmenities `json:"room_amenities"`
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
