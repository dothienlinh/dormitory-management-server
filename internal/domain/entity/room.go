package entity

import (
	"strings"
)

// Room entity
type Room struct {
	Base
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Status         RoomStatus   `json:"status"`
	RoomCategoryID uint         `json:"room_category_id"`
	RoomCategory   RoomCategory `json:"room_category"`
	Users          []User       `json:"users,omitempty"`
}

func (Room) TableName() string {
	return "rooms"
}

// RoomStatus represents the status of a room
type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusUnavailable RoomStatus = "unavailable"
	RoomStatusMaintenance RoomStatus = "maintenance"
)

// RoomCategory entity
type RoomCategory struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
	Rooms       []Room  `json:"rooms,omitempty"`
}

func (RoomCategory) TableName() string {
	return "room_categories"
}

// RoomFilter for filtering rooms
type RoomFilter struct {
	Status         RoomStatus `form:"status" binding:"omitempty,oneof=available unavailable maintenance"`
	Keyword        string     `form:"keyword"`
	RoomCategoryID uint       `form:"room_category_id"`
	Pagination
}

// Build creates the SQL WHERE clause and parameters for the filter
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

// RoomCategoryFilter for filtering room categories
type RoomCategoryFilter struct {
	Keyword string `form:"keyword"`
	Pagination
}

// Build creates the SQL WHERE clause and parameters for the filter
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

// RoomDTO is a data transfer object for Room entity
type RoomDTO struct {
	ID             uint         `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Status         RoomStatus   `json:"status"`
	RoomCategoryID uint         `json:"room_category_id"`
	RoomCategory   RoomCategory `json:"room_category,omitempty"`
	Users          []User       `json:"users,omitempty"`
}

// RoomCategoryDTO is a data transfer object for RoomCategory entity
type RoomCategoryDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
	Rooms       []Room  `json:"rooms,omitempty"`
}
