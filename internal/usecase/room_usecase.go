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

type roomCategoryUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}
