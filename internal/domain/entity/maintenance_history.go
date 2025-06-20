package entity

type MaintenanceHistory struct {
	Base
	Description     string  `json:"description"`
	RoomID          uint64  `json:"room_id"`
	MaintenanceDate string  `json:"maintenance_date"`
	Cost            float64 `json:"cost"`
}

func (MaintenanceHistory) TableName() string {
	return "maintenance_histories"
}

type CreateMaintenanceHistory struct {
	Description     string  `json:"description" binding:"required"`
	RoomID          uint64  `json:"room_id" binding:"required,numeric"`
	MaintenanceDate string  `json:"maintenance_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Cost            float64 `json:"cost" binding:"required,numeric"`
}
