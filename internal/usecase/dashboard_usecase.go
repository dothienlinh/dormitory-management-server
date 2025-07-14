package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"math"

	"go.uber.org/zap"
)

type dashboardUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewDashboardUseCase(repos repository.Repositories, logger logger.Logger) usecase.DashboardUseCase {
	return &dashboardUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *dashboardUseCase) GetStats(ctx context.Context) response.StatusResponse {
	stats, err := uc.repos.Dashboard().GetStats(ctx)
	if err != nil {
		uc.logger.Error("Failed to get stats", zap.Error(err))
		return response.BadRequest("Failed to get stats")
	}

	return response.Success(stats, 1)
}

func (uc *dashboardUseCase) StudentOverview(ctx context.Context, userID uint64) response.StatusResponse {
	overview, err := uc.repos.Dashboard().StudentOverview(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get overview", zap.Error(err))
		return response.InternalServerError(err.Error())
	}

	return response.Success(overview, 0)
}

// GetNotifications retrieves notifications for a user
func (uc *dashboardUseCase) GetNotifications(ctx context.Context, userID uint64, page, limit int, notificationType string, isRead *bool) response.StatusResponse {
	params := &repository.NotificationParams{
		Page:   page,
		Limit:  limit,
		Type:   notificationType,
		IsRead: isRead,
	}

	notifications, total, err := uc.repos.Dashboard().GetNotifications(ctx, userID, params)
	if err != nil {
		uc.logger.Error("Failed to get notifications", zap.Error(err))
		return response.InternalServerError("Failed to get notifications")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	result := &entity.NotificationResponse{
		Data: notifications,
		Pagination: &entity.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(result, 0)
}

// MarkNotificationAsRead marks a notification as read
func (uc *dashboardUseCase) MarkNotificationAsRead(ctx context.Context, notificationID uint64, userID uint64) response.StatusResponse {
	err := uc.repos.Dashboard().MarkNotificationAsRead(ctx, notificationID, userID)
	if err != nil {
		uc.logger.Error("Failed to mark notification as read", zap.Error(err))
		return response.InternalServerError("Failed to mark notification as read")
	}

	return response.Success(nil, 0)
}

// MarkAllNotificationsAsRead marks all notifications as read
func (uc *dashboardUseCase) MarkAllNotificationsAsRead(ctx context.Context, userID uint64, notificationType string) response.StatusResponse {
	err := uc.repos.Dashboard().MarkAllNotificationsAsRead(ctx, userID, notificationType)
	if err != nil {
		uc.logger.Error("Failed to mark all notifications as read", zap.Error(err))
		return response.InternalServerError("Failed to mark all notifications as read")
	}

	return response.Success(nil, 0)
}

// GetEvents retrieves events
func (uc *dashboardUseCase) GetEvents(ctx context.Context, page, limit int, eventType, fromDate, toDate string) response.StatusResponse {
	params := &repository.EventParams{
		Page:     page,
		Limit:    limit,
		Type:     eventType,
		FromDate: fromDate,
		ToDate:   toDate,
	}

	events, total, err := uc.repos.Dashboard().GetEvents(ctx, params)
	if err != nil {
		uc.logger.Error("Failed to get events", zap.Error(err))
		return response.InternalServerError("Failed to get events")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	result := &entity.EventResponse{
		Data: events,
		Pagination: &entity.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(result, 0)
}

// GetServiceRequests retrieves service requests for a user
func (uc *dashboardUseCase) GetServiceRequests(ctx context.Context, userID uint64, page, limit int, status, category string) response.StatusResponse {
	params := &repository.ServiceRequestParams{
		Page:     page,
		Limit:    limit,
		Status:   status,
		Category: category,
	}

	requests, total, err := uc.repos.Dashboard().GetServiceRequests(ctx, userID, params)
	if err != nil {
		uc.logger.Error("Failed to get service requests", zap.Error(err))
		return response.InternalServerError("Failed to get service requests")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	result := &entity.ServiceRequestResponse{
		Data: requests,
		Pagination: &entity.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(result, 0)
}

// GetQuickStats retrieves quick stats for dashboard
func (uc *dashboardUseCase) GetQuickStats(ctx context.Context, userID uint64) response.StatusResponse {
	// Get limited notifications, events, and service requests for quick view
	notifications, _, err := uc.repos.Dashboard().GetNotifications(ctx, userID, &repository.NotificationParams{
		Page:   1,
		Limit:  5,
		IsRead: &[]bool{false}[0], // Only unread notifications
	})
	if err != nil {
		uc.logger.Error("Failed to get notifications for quick stats", zap.Error(err))
	}

	events, _, err := uc.repos.Dashboard().GetEvents(ctx, &repository.EventParams{
		Page:  1,
		Limit: 5,
	})
	if err != nil {
		uc.logger.Error("Failed to get events for quick stats", zap.Error(err))
	}

	serviceRequests, _, err := uc.repos.Dashboard().GetServiceRequests(ctx, userID, &repository.ServiceRequestParams{
		Page:  1,
		Limit: 5,
	})
	if err != nil {
		uc.logger.Error("Failed to get service requests for quick stats", zap.Error(err))
	}

	quickStats := map[string]interface{}{
		"notifications":    notifications,
		"events":           events,
		"service_requests": serviceRequests,
	}

	return response.Success(quickStats, 0)
}
