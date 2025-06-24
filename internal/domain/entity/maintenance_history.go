package entity

import "time"

type MaintenanceHistory struct {
	Base
	Description     string    `json:"description"`
	RoomID          uint64    `json:"room_id"`
	MaintenanceDate time.Time `json:"maintenance_date"`
	Cost            float64   `json:"cost"`
}

func (MaintenanceHistory) TableName() string {
	return "maintenance_histories"
}

type CreateMaintenanceHistory struct {
	Description     string    `json:"description" binding:"required"`
	RoomID          uint64    `json:"room_id" binding:"required,numeric"`
	MaintenanceDate time.Time `json:"maintenance_date" binding:"required,validdate"`
	Cost            float64   `json:"cost" binding:"required,numeric"`
}
