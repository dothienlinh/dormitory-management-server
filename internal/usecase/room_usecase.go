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
func (uc *roomUseCase) CreateRoom(ctx context.Context, room *entity.Room) response.StatusResponse {
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
func (uc *roomUseCase) UpdateRoom(ctx context.Context, id uint, roomData *entity.Room) response.StatusResponse {
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

	if len(room.RoomRents) > 0 {
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

// NewRoomCategoryUseCase creates a new room category use case
func NewRoomCategoryUseCase(repos repository.Repositories, logger logger.Logger) usecase.RoomCategoryUseCase {
	return &roomCategoryUseCase{
		repos:  repos,
		logger: logger,
	}
}

// CreateRoomCategory creates a new room category
func (uc *roomCategoryUseCase) CreateRoomCategory(ctx context.Context, category *entity.RoomCategory) response.StatusResponse {
	if err := uc.repos.RoomCategory().Create(ctx, category); err != nil {
		uc.logger.Error("Failed to create room category", zap.Error(err))
		return response.InternalServerError("Failed to create room category")
	}

	return response.Created(category)
}

// GetRoomCategoryByID retrieves a room category by ID
func (uc *roomCategoryUseCase) GetRoomCategoryByID(ctx context.Context, id uint) response.StatusResponse {
	category, err := uc.repos.RoomCategory().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room category by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room category with ID %d not found", id))
	}

	return response.Success(category, 1)
}

// GetListRoomCategories retrieves room categories based on filter
func (uc *roomCategoryUseCase) GetListRoomCategories(ctx context.Context, filter *entity.RoomCategoryFilter) response.StatusResponse {
	categories, total, err := uc.repos.RoomCategory().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of room categories", zap.Error(err))
		return response.InternalServerError("Failed to get list of room categories")
	}

	return response.Success(categories, total)
}

// UpdateRoomCategory updates a room category
func (uc *roomCategoryUseCase) UpdateRoomCategory(ctx context.Context, id uint, categoryData *entity.RoomCategory) response.StatusResponse {
	category, err := uc.repos.RoomCategory().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room category for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room category with ID %d not found", id))
	}

	// Update category fields
	if categoryData.Name != "" {
		category.Name = categoryData.Name
	}
	if categoryData.Description != "" {
		category.Description = categoryData.Description
	}
	if categoryData.Capacity > 0 {
		category.Capacity = categoryData.Capacity
	}
	if categoryData.Price > 0 {
		category.Price = categoryData.Price
	}

	if err := uc.repos.RoomCategory().Update(ctx, category); err != nil {
		uc.logger.Error("Failed to update room category", zap.Error(err))
		return response.InternalServerError("Failed to update room category")
	}

	return response.Success(category, 1)
}

// DeleteRoomCategory deletes a room category
func (uc *roomCategoryUseCase) DeleteRoomCategory(ctx context.Context, id uint) response.StatusResponse {
	// Check if category has active rooms
	category, err := uc.repos.RoomCategory().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room category for deletion", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room category with ID %d not found", id))
	}

	if len(category.Rooms) > 0 {
		return response.BadRequest("Cannot delete room category with active rooms")
	}

	if err := uc.repos.RoomCategory().Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete room category", zap.Error(err))
		return response.InternalServerError("Failed to delete room category")
	}

	return response.Success("Room category deleted successfully", 0)
}
