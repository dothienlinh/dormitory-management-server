package models

import "strings"

type RoomCategory struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Capacity    int    `json:"capacity" gorm:"type:int;not null;check:capacity in (8, 6, 4)"`
	Price       int    `json:"price" gorm:"type:int;not null"`
	Acreage     int    `json:"acreage" gorm:"type:int;not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Rooms       []Room `json:"-"`
}

type CreateRoomCategory struct {
	Name        string `json:"name" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,oneof=8 6 4"`
	Price       int    `json:"price" binding:"required"`
	Acreage     int    `json:"acreage" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ListRoomCategory struct {
	BaseModel
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	Price       int    `json:"price"`
	Acreage     int    `json:"acreage"`
	Description string `json:"description"`
}

type UpdateRoomCategory struct {
	CreateRoomCategory
}

type RoomCategoryDetail struct {
	BaseModel
	Name        string       `json:"name"`
	Capacity    int          `json:"capacity"`
	Price       int          `json:"price"`
	Acreage     int          `json:"acreage"`
	Description string       `json:"description"`
	Rooms       []RoomSimple `json:"rooms"`
}

type FilterRoomCategory struct {
	Name string `form:"name"`
	Pagination
}

func (f *FilterRoomCategory) Build() (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}

	if f.Name != "" {
		conditions = append(conditions, "name ILIKE ?")
		values = append(values, "%"+f.Name+"%")
	}

	whereClause := strings.Join(conditions, " AND ")
	return whereClause, values
}
