package models

type Room struct {
	BaseModel
	RoomNumber     string        `json:"room_number" gorm:"not null;unique"`
	Status         RoomStatus    `json:"status" gorm:"type:room_status;not null"`
	RoomCategoryID uint          `json:"-" gorm:"not null"`
	RoomCategory   *RoomCategory `json:"room_category"`
	RoomRents      []RoomRent    `json:"-" gorm:"foreignKey:RoomID"`
	Contracts      []Contract    `json:"-" gorm:"foreignKey:RoomID"`
}

type RoomSimple struct {
	BaseModel
	RoomNumber string `json:"room_number"`
	Status     string `json:"status"`
}

type RoomStatus string

const (
	Available   RoomStatus = "available"
	Occupied    RoomStatus = "occupied"
	Maintenance RoomStatus = "maintenance"
)

type FilterRoom struct {
	RoomNumber     string     `form:"room_number"`
	Status         RoomStatus `form:"status" validate:"oneof=available occupied maintenance"`
	RoomCategoryID uint       `form:"room_category_id"`
	Pagination
}

type CreateRoom struct {
	RoomNumber     string     `json:"room_number" validate:"required"`
	Status         RoomStatus `json:"status" validate:"required,oneof=available occupied maintenance"`
	RoomCategoryID uint       `json:"room_category_id" validate:"required"`
}

type UpdateRoom struct {
	Status         RoomStatus `json:"status" validate:"required,oneof=available occupied maintenance"`
	RoomCategoryID uint       `json:"room_category_id" validate:"required"`
}
