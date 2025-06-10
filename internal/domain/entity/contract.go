package entity

import (
	"time"
)

// ContractStatus represents the status of a contract
type ContractStatus string

const (
	ContractStatusActive    ContractStatus = "active"
	ContractStatusInactive  ContractStatus = "inactive"
	ContractStatusCancelled ContractStatus = "cancelled"
)

// Contract entity
type Contract struct {
	Base
	UserID      uint           `json:"user_id"`
	User        User           `json:"user"`
	RoomID      uint           `json:"room_id"`
	Room        Room           `json:"room"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Price       float64        `json:"price"`
	Status      ContractStatus `json:"status"`
	Description string         `json:"description"`
	Code        string         `json:"code"`
}

func (Contract) TableName() string {
	return "contracts"
}

// CreateContract is the data transfer object for creating a contract
type CreateContract struct {
	UserID      uint           `json:"user_id" binding:"required,numeric"`
	RoomID      uint           `json:"room_id" binding:"required,numeric"`
	StartDate   time.Time      `json:"start_date" binding:"required,datetime"`
	EndDate     time.Time      `json:"end_date" binding:"required,datetime"`
	Price       float64        `json:"price" binding:"required,numeric"`
	Status      ContractStatus `json:"status" binding:"omitempty,oneof=active inactive cancelled"`
	Description string         `json:"description" binding:"omitempty"`
}

// UpdateContract is the data transfer object for updating a contract
type UpdateContract struct {
	Status      ContractStatus `json:"status" binding:"required,oneof=active inactive cancelled"`
	StartDate   *time.Time     `json:"start_date" binding:"omitempty,datetime"`
	EndDate     *time.Time     `json:"end_date" binding:"omitempty,datetime"`
	Price       *float64       `json:"price" binding:"omitempty,numeric"`
	Description *string        `json:"description" binding:"omitempty"`
}

// ContractDTO is a data transfer object for Contract entity
type ContractDTO struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user,omitempty"`
	RoomID      uint           `json:"room_id"`
	Room        Room           `json:"room,omitempty"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Price       float64        `json:"price"`
	Status      ContractStatus `json:"status"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ContractFilter for filtering contracts
type ContractFilter struct {
	Status  ContractStatus `form:"status" binding:"omitempty,oneof=active inactive cancelled"`
	UserID  uint           `form:"user_id"`
	RoomID  uint           `form:"room_id"`
	Keyword string         `form:"keyword"`
	Pagination
}
