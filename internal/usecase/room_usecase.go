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

// roomUseCase implements the usecase.RoomUseCase interface
type roomUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

// NewRoomUseCase creates a new room use case
func NewRoomUseCase(repos repository.Repositories, logger logger.Logger) usecase.RoomUseCase {
	return &roomUseCase{
		repos:  repos,
		logger: logger,
	}
}

// CreateRoom creates a new room
func (uc *roomUseCase) CreateRoom(ctx context.Context, room *entity.CreateRoom) response.StatusResponse {
	// Validate room category
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

// GetRoomByID retrieves a room by ID
func (uc *roomUseCase) GetRoomByID(ctx context.Context, id uint) response.StatusResponse {
	room, err := uc.repos.Room().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room with ID %d not found", id))
	}

	return response.Success(room, 1)
}

// GetListRooms retrieves rooms based on filter
func (uc *roomUseCase) GetListRooms(ctx context.Context, filter *entity.RoomFilter) response.StatusResponse {
	rooms, total, err := uc.repos.Room().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of rooms", zap.Error(err))
		return response.InternalServerError("Failed to get list of rooms")
	}

	return response.Success(rooms, total)
}

// UpdateRoom updates a room
func (uc *roomUseCase) UpdateRoom(ctx context.Context, id uint, roomData *entity.UpdateRoom) response.StatusResponse {
	room, err := uc.repos.Room().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room with ID %d not found", id))
	}

	// Update room fields
	if roomData.Name != "" {
		room.Name = roomData.Name
	}
	if roomData.Description != "" {
		room.Description = roomData.Description
	}
	if roomData.Status != "" {
		room.Status = roomData.Status
	}
	if roomData.RoomCategoryID != 0 {
		// Validate room category
		if _, err := uc.repos.RoomCategory().GetByID(ctx, roomData.RoomCategoryID); err != nil {
			uc.logger.Error("Failed to validate room category", zap.Error(err))
			return response.BadRequest("Invalid room category")
		}
		room.RoomCategoryID = roomData.RoomCategoryID
	}

	if err := uc.repos.Room().Update(ctx, room); err != nil {
		uc.logger.Error("Failed to update room", zap.Error(err))
		return response.InternalServerError("Failed to update room")
	}

	return response.Success(room, 1)
}

// DeleteRoom deletes a room
func (uc *roomUseCase) DeleteRoom(ctx context.Context, id uint) response.StatusResponse {
	// Check if room has active students
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
		return response.InternalServerError("Failed to delete room")
	}

	return response.Success("Room deleted successfully", 0)
}

// roomCategoryUseCase implements the usecase.RoomCategoryUseCase interface
type roomCategoryUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}
