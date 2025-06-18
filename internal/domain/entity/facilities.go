package entity

type Facilities struct {
	Base
	Name string `json:"name"`
}

func (Facilities) TableName() string {
	return "facilities"
}

type CreateFacility struct {
	Name string `json:"name" binding:"required"`
}

type UpdateFacility struct {
	Name string `json:"name" binding:"required"`
}
