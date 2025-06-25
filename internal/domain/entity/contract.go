package entity

import (
	"time"
)

type ContractStatus string

const (
	ContractStatusActive    ContractStatus = "active"
	ContractStatusInactive  ContractStatus = "inactive"
	ContractStatusCancelled ContractStatus = "cancelled"
)

type Contract struct {
	Base
	UserID      uint64         `json:"user_id"`
	User        User           `json:"user"`
	RoomID      uint64         `json:"room_id"`
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

type CreateContract struct {
	UserID      uint64  `json:"user_id" binding:"required,numeric"`
	RoomID      uint64  `json:"room_id" binding:"required,numeric"`
	StartDate   string  `json:"start_date" binding:"required,validdate"`
	EndDate     string  `json:"end_date" binding:"required,validdate,gtefield=StartDate"`
	Price       float64 `json:"price" binding:"required,numeric"`
	Description string  `json:"description" binding:"omitempty"`
}

type UpdateContract struct {
	Status      ContractStatus `json:"status" binding:"required,oneof=active inactive cancelled"`
	StartDate   *time.Time     `json:"start_date" binding:"omitempty,validdate"`
	EndDate     *time.Time     `json:"end_date" binding:"omitempty,validdate,gtefield=StartDate"`
	Price       *float64       `json:"price" binding:"omitempty,numeric"`
	Description *string        `json:"description" binding:"omitempty"`
}

type ContractDTO struct {
	ID          uint64         `json:"id"`
	UserID      uint64         `json:"user_id"`
	User        User           `json:"user"`
	RoomID      uint64         `json:"room_id"`
	Room        Room           `json:"room"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Price       float64        `json:"price"`
	Status      ContractStatus `json:"status"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ContractFilter struct {
	Status  ContractStatus `form:"status" binding:"omitempty,oneof=active inactive cancelled"`
	Keyword string         `form:"keyword"`
	Pagination
}
