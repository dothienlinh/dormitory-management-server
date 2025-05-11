package models

import "time"

type Contract struct {
	BaseModel
	RoomID    uint           `json:"room_id" gorm:"not null"`
	Room      *Room          `json:"room"`
	UserId    uint           `json:"user_id" gorm:"not null"`
	User      *User          `json:"user"`
	StartDate *time.Time     `json:"start_date" gorm:"not null"`
	EndDate   *time.Time     `json:"end_date" gorm:"not null"`
	Status    ContractStatus `json:"status" gorm:"type:contract_status;not null"`
}

type ContractStatus string

const (
	ContractStatusActive   ContractStatus = "active"
	ContractStatusInactive ContractStatus = "inactive"
)

type CreateContract struct {
	RoomID    uint           `json:"room_id" binding:"required"`
	UserId    uint           `json:"user_id" binding:"required"`
	StartDate *time.Time     `json:"start_date" binding:"required"`
	EndDate   *time.Time     `json:"end_date" binding:"required"`
	Status    ContractStatus `json:"status" binding:"required,oneof=active inactive"`
}

type UpdateContract struct {
	EndDate *time.Time     `json:"end_date" binding:"required"`
	Status  ContractStatus `json:"status" binding:"required,oneof=active inactive"`
}
