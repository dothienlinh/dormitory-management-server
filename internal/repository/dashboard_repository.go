package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) repository.DashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (*entity.DashboardStats, error) {
	stats := &entity.DashboardStats{
		MonthRevenue: entity.MonthRevenue{
			Month: time.Now().Month().String(),
		},
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("role = ?", entity.UserRoleStudent).Count(&stats.TotalStudents).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Count(&stats.TotalRooms).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table(entity.Contract{}.TableName()).Count(&stats.TotalContracts).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Table(entity.Payment{}.TableName()).Select("MONTH(created_at) as month, SUM(amount) as amount").Group("MONTH(created_at)").Scan(&stats.MonthRevenue).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *dashboardRepository) StudentOverview(ctx context.Context, userID uint64) (*entity.ResponseOverview, error) {
	overview := &entity.ResponseOverview{}

	user := &entity.User{}
	if err := r.db.WithContext(ctx).Table(user.TableName()).Where("id = ?", userID).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	overview.Student = user

	if user.RoomID != nil {
		room := &entity.Room{}
		if err := r.db.WithContext(ctx).Table(room.TableName()).Where("id = ?", user.RoomID).First(room).Error; err != nil {
			return nil, err
		}
		overview.Room = room

		roomMates := []*entity.User{}
		if err := r.db.WithContext(ctx).Table(user.TableName()).Where("room_id = ? AND id != ?", user.RoomID, userID).Find(&roomMates).Error; err != nil {
			return nil, err
		}
		overview.RoomMates = roomMates

		contract := &entity.Contract{}
		if err := r.db.WithContext(ctx).Table(contract.TableName()).Where("room_id = ? AND user_id = ?", user.RoomID, userID).First(contract).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			overview.Contract = contract
		}

		payment := &entity.Payment{}
		if err := r.db.WithContext(ctx).Table(payment.TableName()).Where("user_id = ?", userID).Order("created_at DESC").First(payment).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		} else {
			overview.Payment = payment
		}
	}
	// Get actual stats
	stats := &entity.Stats{}

	// Count service requests
	var totalRequests, pendingRequests, completedRequests int64
	r.db.WithContext(ctx).Model(&entity.ServiceRequest{}).Where("user_id = ?", userID).Count(&totalRequests)
	r.db.WithContext(ctx).Model(&entity.ServiceRequest{}).Where("user_id = ? AND status = ?", userID, entity.ServiceRequestStatusPending).Count(&pendingRequests)
	r.db.WithContext(ctx).Model(&entity.ServiceRequest{}).Where("user_id = ? AND status = ?", userID, entity.ServiceRequestStatusCompleted).Count(&completedRequests)

	stats.TotalRequests = int(totalRequests)
	stats.PendingRequests = int(pendingRequests)
	stats.CompletedRequests = int(completedRequests)

	// Count unread notifications
	var unreadCount int64
	r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)
	stats.UnreadNotifications = int(unreadCount)

	overview.Stats = stats

	return overview, nil
}

// GetNotifications retrieves notifications for a user with pagination
func (r *dashboardRepository) GetNotifications(ctx context.Context, userID uint64, params *repository.NotificationParams) ([]*entity.Notification, int, error) {
	if params == nil {
		params = &repository.NotificationParams{Page: 1, Limit: 10}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var notifications []*entity.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ?", userID)

	// Apply filters
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, int(total), nil
}

// MarkNotificationAsRead marks a single notification as read
func (r *dashboardRepository) MarkNotificationAsRead(ctx context.Context, notificationID uint64, userID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

// MarkAllNotificationsAsRead marks all notifications as read for a user
func (r *dashboardRepository) MarkAllNotificationsAsRead(ctx context.Context, userID uint64, notificationType string) error {
	now := time.Now()
	query := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false)

	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	return query.Updates(map[string]interface{}{
		"is_read": true,
		"read_at": &now,
	}).Error
}

// GetEvents retrieves events with pagination
func (r *dashboardRepository) GetEvents(ctx context.Context, params *repository.EventParams) ([]*entity.Event, int, error) {
	if params == nil {
		params = &repository.EventParams{Page: 1, Limit: 10}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var events []*entity.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Event{})

	// Apply filters
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.FromDate != "" {
		query = query.Where("event_date >= ?", params.FromDate)
	}
	if params.ToDate != "" {
		query = query.Where("event_date <= ?", params.ToDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("event_date ASC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, int(total), nil
}

// GetServiceRequests retrieves service requests for a user with pagination
func (r *dashboardRepository) GetServiceRequests(ctx context.Context, userID uint64, params *repository.ServiceRequestParams) ([]*entity.ServiceRequest, int, error) {
	if params == nil {
		params = &repository.ServiceRequestParams{Page: 1, Limit: 10}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	var requests []*entity.ServiceRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ServiceRequest{}).Where("user_id = ?", userID)

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}

	return requests, int(total), nil
}
