package models

type Amenities struct {
	BaseModel
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
}
