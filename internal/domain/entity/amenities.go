package entity

type Amenities struct {
	Base
	Name string `json:"name"`
}

func (Amenities) TableName() string {
	return "amenities"
}

type CreateAmenity struct {
	Name string `json:"name" binding:"required"`
}

type UpdateAmenity struct {
	Name string `json:"name" binding:"required"`
}
