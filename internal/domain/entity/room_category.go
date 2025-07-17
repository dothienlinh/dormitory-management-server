package entity

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type RoomCategory struct {
	ID          uint64         `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Capacity    int            `json:"capacity"`
	Price       float64        `json:"price"`
	Acreage     int            `json:"acreage"`
	Rooms       []Room         `json:"rooms"`
}

func (RoomCategory) TableName() string {
	return "room_categories"
}

type CreateRoomCategory struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"omitempty"`
	Capacity    int     `json:"capacity" binding:"required,numeric"`
	Price       float64 `json:"price" binding:"required,numeric"`
}

type UpdateRoomCategory struct {
	Name        string  `json:"name" binding:"omitempty"`
	Description string  `json:"description" binding:"omitempty"`
	Capacity    int     `json:"capacity" binding:"omitempty,numeric"`
	Price       float64 `json:"price" binding:"omitempty,numeric"`
}

type RoomFilter struct {
	Status         RoomStatus `form:"status" binding:"omitempty,oneof=available unavailable maintenance"`
	Keyword        string     `form:"keyword"`
	RoomCategoryID uint       `form:"room_category_id"`
	Pagination
}

func (f RoomFilter) Build() (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}

	if f.Status != "" {
		conditions = append(conditions, "status = ?")
		values = append(values, f.Status)
	}

	if f.RoomCategoryID > 0 {
		conditions = append(conditions, "room_category_id = ?")
		values = append(values, f.RoomCategoryID)
	}

	if f.Keyword != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		keyword := "%" + f.Keyword + "%"
		values = append(values, keyword, keyword)
	}

	if len(conditions) == 0 {
		return "", values
	}

	whereClause := strings.Join(conditions, " AND ")

	return whereClause, values
}

type RoomCategoryFilter struct {
	Keyword string `form:"keyword"`
	Pagination
}

func (f RoomCategoryFilter) Build() (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}

	if f.Keyword != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		keyword := "%" + f.Keyword + "%"
		values = append(values, keyword, keyword)
	}

	if len(conditions) == 0 {
		return "", values
	}

	whereClause := strings.Join(conditions, " AND ")

	return whereClause, values
}

type RoomDTO struct {
	ID             uint         `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Status         RoomStatus   `json:"status"`
	RoomCategoryID uint         `json:"room_category_id"`
	RoomCategory   RoomCategory `json:"room_category"`
	Users          []User       `json:"users"`
}

type RoomCategoryDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
	Rooms       []Room  `json:"rooms"`
}
