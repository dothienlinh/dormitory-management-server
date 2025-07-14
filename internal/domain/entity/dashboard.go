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

type Stats struct {
	TotalRequests       int `json:"total_requests"`
	PendingRequests     int `json:"pending_requests"`
	CompletedRequests   int `json:"completed_requests"`
	UnreadNotifications int `json:"unread_notifications"`
}

type ResponseOverview struct {
	Student         *User             `json:"student"`
	Room            *Room             `json:"room"`
	Contract        *Contract         `json:"contract"`
	Payment         *Payment          `json:"payment"`
	RoomMates       []*User           `json:"room_mates"`
	Stats           *Stats            `json:"stats"`
	Notifications   []*Notification   `json:"notifications,omitempty"`
	Events          []*Event          `json:"events,omitempty"`
	ServiceRequests []*ServiceRequest `json:"service_requests,omitempty"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type NotificationResponse struct {
	Data       []*Notification     `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

type EventResponse struct {
	Data       []*Event            `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

type ServiceRequestResponse struct {
	Data       []*ServiceRequest   `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}
