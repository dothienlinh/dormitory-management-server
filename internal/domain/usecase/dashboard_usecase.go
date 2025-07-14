package usecase

import (
	"context"
	"dormitory_management/internal/domain/response"
)

type DashboardUseCase interface {
	GetStats(ctx context.Context) response.StatusResponse
	StudentOverview(ctx context.Context, userID uint64) response.StatusResponse

	// Notification methods
	GetNotifications(ctx context.Context, userID uint64, page, limit int, notificationType string, isRead *bool) response.StatusResponse
	MarkNotificationAsRead(ctx context.Context, notificationID uint64, userID uint64) response.StatusResponse
	MarkAllNotificationsAsRead(ctx context.Context, userID uint64, notificationType string) response.StatusResponse

	// Event methods
	GetEvents(ctx context.Context, page, limit int, eventType, fromDate, toDate string) response.StatusResponse

	// Service Request methods
	GetServiceRequests(ctx context.Context, userID uint64, page, limit int, status, category string) response.StatusResponse

	// Quick stats
	GetQuickStats(ctx context.Context, userID uint64) response.StatusResponse
}
