package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (*entity.DashboardStats, error)
	StudentOverview(ctx context.Context, userID uint64) (*entity.ResponseOverview, error)

	// Notification methods
	GetNotifications(ctx context.Context, userID uint64, params *NotificationParams) ([]*entity.Notification, int, error)
	MarkNotificationAsRead(ctx context.Context, notificationID uint64, userID uint64) error
	MarkAllNotificationsAsRead(ctx context.Context, userID uint64, notificationType string) error

	// Event methods
	GetEvents(ctx context.Context, params *EventParams) ([]*entity.Event, int, error)

	// Service Request methods
	GetServiceRequests(ctx context.Context, userID uint64, params *ServiceRequestParams) ([]*entity.ServiceRequest, int, error)
}

type NotificationParams struct {
	Page   int
	Limit  int
	Type   string
	IsRead *bool
}

type EventParams struct {
	Page     int
	Limit    int
	Type     string
	FromDate string
	ToDate   string
}

type ServiceRequestParams struct {
	Page     int
	Limit    int
	Status   string
	Category string
}
