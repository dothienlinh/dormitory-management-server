package entity

type DashboardStats struct {
	TotalStudents  int64        `json:"total_students"`
	TotalRooms     int64        `json:"total_rooms"`
	TotalContracts int64        `json:"total_contracts"`
	MonthRevenue   MonthRevenue `json:"month_revenue"`
}

type MonthRevenue struct {
	Month  string `json:"month"`
	Amount int64  `json:"amount"`
}
