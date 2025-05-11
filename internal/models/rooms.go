package models

import "strings"

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
	Status         RoomStatus `form:"status" binding:"oneof=available occupied maintenance"`
	RoomCategoryID uint       `form:"room_category_id"`
	Pagination
}

func (f *FilterRoom) Build() (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}

	if f.RoomNumber != "" {
		conditions = append(conditions, "room_number ILIKE ?")
		values = append(values, "%"+f.RoomNumber+"%")
	}

	if f.Status != "" {
		conditions = append(conditions, "status = ?")
		values = append(values, f.Status)
	}

	if f.RoomCategoryID != 0 {
		conditions = append(conditions, "room_category_id = ?")
		values = append(values, f.RoomCategoryID)
	}

	whereClause := strings.Join(conditions, " AND ")

	return whereClause, values
}

type CreateRoom struct {
	RoomNumber     string     `json:"room_number" binding:"required"`
	Status         RoomStatus `json:"status" binding:"required,oneof=available occupied maintenance"`
	RoomCategoryID uint       `json:"room_category_id" binding:"required"`
}

type UpdateRoom struct {
	Status         RoomStatus `json:"status" binding:"required,oneof=available occupied maintenance"`
	RoomCategoryID uint       `json:"room_category_id" binding:"required"`
}
