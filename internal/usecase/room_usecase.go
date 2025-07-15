package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

type roomUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewRoomUseCase(repos repository.Repositories, logger logger.Logger) usecase.RoomUseCase {
	return &roomUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *roomUseCase) CreateRoom(ctx context.Context, room *entity.CreateRoom) response.StatusResponse {
	if _, err := uc.repos.RoomCategory().GetByID(ctx, room.RoomCategoryID); err != nil {
		uc.logger.Error("Failed to validate room category", zap.Error(err))
		return response.BadRequest("Invalid room category")
	}

	if err := uc.repos.Room().Create(ctx, room); err != nil {
		uc.logger.Error("Failed to create room", zap.Error(err))
		return response.InternalServerError("Failed to create room")
	}

	return response.Created(room)
}

func (uc *roomUseCase) GetRoomByID(ctx context.Context, id uint64) response.StatusResponse {
	room, err := uc.repos.Room().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room with ID %d not found", id))
	}

	return response.Success(room, 1)
}

func (uc *roomUseCase) GetListRooms(ctx context.Context, filter *entity.RoomFilter) response.StatusResponse {
	rooms, total, err := uc.repos.Room().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of rooms", zap.Error(err))
		return response.InternalServerError("Failed to get list of rooms")
	}

	return response.Success(rooms, total)
}

func (uc *roomUseCase) UpdateRoom(ctx context.Context, id uint64, roomData *entity.UpdateRoom) response.StatusResponse {
	room, err := uc.repos.Room().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room with ID %d not found", id))
	}

	if roomData.RoomNumber != "" {
		room.RoomNumber = roomData.RoomNumber
	}
	if roomData.Status != "" {
		room.Status = roomData.Status
	}
	if roomData.RoomCategoryID != 0 {
		if _, err := uc.repos.RoomCategory().GetByID(ctx, roomData.RoomCategoryID); err != nil {
			uc.logger.Error("Failed to validate room category", zap.Error(err))
			return response.BadRequest("Invalid room category")
		}
		room.RoomCategoryID = roomData.RoomCategoryID
	}

	if err := uc.repos.Room().Update(ctx, room, roomData.AmenityIDs); err != nil {
		uc.logger.Error("Failed to update room", zap.Error(err))
		return response.InternalServerError("Failed to update room")
	}

	return response.Success(room, 1)
}

func (uc *roomUseCase) DeleteRoom(ctx context.Context, id uint64) response.StatusResponse {
	room, err := uc.repos.Room().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room for deletion", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room with ID %d not found", id))
	}

	if len(room.Users) > 0 {
		return response.BadRequest("Cannot delete room with active students")
	}

	if err := uc.repos.Room().Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete room", zap.Error(err))
		return response.InternalServerError(err.Error())
	}

	return response.Success("Room deleted successfully", 0)
}

// GetStudentRoomDetails retrieves detailed room information for a student
func (uc *roomUseCase) GetStudentRoomDetails(ctx context.Context, userID uint64) response.StatusResponse {
	if userID == 0 {
		uc.logger.Error("Missing user ID in GetStudentRoomDetails")
		return response.BadRequest("User ID is required")
	}
	details, err := uc.repos.Room().GetStudentRoomDetails(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get student room details", zap.Error(err))
		return response.InternalServerError("Failed to get room details")
	}

	return response.Success(details, 0)
}

// GetRoomStats retrieves room statistics for a student
func (uc *roomUseCase) GetRoomStats(ctx context.Context, userID uint64) response.StatusResponse {
	stats, err := uc.repos.Room().GetRoomStats(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get room stats", zap.Error(err))
		return response.InternalServerError("Failed to get room stats")
	}

	return response.Success(stats, 0)
}

// GetRoomIssues retrieves room issues with pagination
func (uc *roomUseCase) GetRoomIssues(ctx context.Context, userID uint64, page, limit int, status, category string) response.StatusResponse {
	if userID == 0 {
		uc.logger.Error("Missing user ID in GetRoomIssues")
		return response.BadRequest("User ID is required")
	}
	params := &repository.RoomIssueParams{
		Page:     page,
		Limit:    limit,
		Status:   status,
		Category: category,
	}

	issues, total, err := uc.repos.Room().GetRoomIssues(ctx, userID, params)
	if err != nil {
		uc.logger.Error("Failed to get room issues", zap.Error(err))
		return response.InternalServerError("Failed to get room issues")
	}

	result := map[string]interface{}{
		"data": issues,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	return response.Success(result, 0)
}

// CreateRoomIssue creates a new room issue
func (uc *roomUseCase) CreateRoomIssue(ctx context.Context, userID uint64, title, description, category, priority string) response.StatusResponse {
	// Get user's room first
	user := &entity.User{}
	user.ID = userID
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user", zap.Error(err))
		return response.BadRequest("User not found")
	}

	if user.RoomID == nil {
		return response.BadRequest("User has no room assigned")
	}

	issue := &entity.RoomIssue{
		Title:       title,
		Description: description,
		Category:    category,
		Priority:    priority,
		RoomID:      uint64(*user.RoomID),
		ReportedBy:  userID,
		Status:      "pending",
	}

	if err := uc.repos.Room().CreateRoomIssue(ctx, issue); err != nil {
		uc.logger.Error("Failed to create room issue", zap.Error(err))
		return response.InternalServerError("Failed to create room issue")
	}

	return response.Created(issue)
}

// GetRoomIssueDetails retrieves details of a specific room issue
func (uc *roomUseCase) GetRoomIssueDetails(ctx context.Context, issueID uint64, userID uint64) response.StatusResponse {
	issue, err := uc.repos.Room().GetRoomIssueByID(ctx, issueID, userID)
	if err != nil {
		uc.logger.Error("Failed to get room issue details", zap.Error(err))
		return response.NotFound("Issue not found")
	}

	return response.Success(issue, 0)
}

// GetRoomBills retrieves room bills with pagination
func (uc *roomUseCase) GetRoomBills(ctx context.Context, userID uint64, page, limit, year, month int, status string) response.StatusResponse {
	params := &repository.RoomBillParams{
		Page:   page,
		Limit:  limit,
		Year:   year,
		Month:  month,
		Status: status,
	}

	bills, total, err := uc.repos.Room().GetRoomBills(ctx, userID, params)
	if err != nil {
		uc.logger.Error("Failed to get room bills", zap.Error(err))
		return response.InternalServerError("Failed to get room bills")
	}

	result := map[string]interface{}{
		"data": bills,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	return response.Success(result, 0)
}

// GetRoomRules retrieves all room rules
func (uc *roomUseCase) GetRoomRules(ctx context.Context) response.StatusResponse {
	rules, err := uc.repos.Room().GetRoomRules(ctx)
	if err != nil {
		uc.logger.Error("Failed to get room rules", zap.Error(err))
		return response.InternalServerError("Failed to get room rules")
	}

	return response.Success(rules, 0)
}

// GetCleaningSchedule retrieves cleaning schedule for user's room
func (uc *roomUseCase) GetCleaningSchedule(ctx context.Context, userID uint64) response.StatusResponse {
	if userID == 0 {
		uc.logger.Error("Missing user ID in GetCleaningSchedule")
		return response.BadRequest("User ID is required")
	}
	// Get user's room first
	user := &entity.User{}
	user.ID = userID
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user", zap.Error(err))
		return response.BadRequest("User not found")
	}

	if user.RoomID == nil {
		return response.BadRequest("User has no room assigned")
	}

	schedule, err := uc.repos.Room().GetCleaningSchedule(ctx, uint64(*user.RoomID))
	if err != nil {
		uc.logger.Error("Failed to get cleaning schedule", zap.Error(err))
		return response.InternalServerError("Failed to get cleaning schedule")
	}

	return response.Success(schedule, 0)
}

// GetRoommates retrieves roommates for a user
func (uc *roomUseCase) GetRoommates(ctx context.Context, userID uint64) response.StatusResponse {
	roommates, err := uc.repos.Room().GetRoommates(ctx, userID)
	if err != nil {
		uc.logger.Error("Failed to get roommates", zap.Error(err))
		return response.InternalServerError("Failed to get roommates")
	}

	return response.Success(roommates, 0)
}
