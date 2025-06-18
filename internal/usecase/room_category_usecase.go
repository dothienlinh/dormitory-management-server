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

func NewRoomCategoryUseCase(repos repository.Repositories, logger logger.Logger) usecase.RoomCategoryUseCase {
	return &roomCategoryUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *roomCategoryUseCase) CreateRoomCategory(ctx context.Context, category *entity.CreateRoomCategory) response.StatusResponse {
	if err := uc.repos.RoomCategory().Create(ctx, category); err != nil {
		uc.logger.Error("Failed to create room category", zap.Error(err))
		return response.InternalServerError("Failed to create room category")
	}

	return response.Created(category)
}

func (uc *roomCategoryUseCase) GetRoomCategoryByID(ctx context.Context, id uint64) response.StatusResponse {
	category, err := uc.repos.RoomCategory().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room category by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room category with ID %d not found", id))
	}

	return response.Success(category, 1)
}

func (uc *roomCategoryUseCase) GetListRoomCategories(ctx context.Context, filter *entity.RoomCategoryFilter) response.StatusResponse {
	categories, total, err := uc.repos.RoomCategory().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of room categories", zap.Error(err))
		return response.InternalServerError("Failed to get list of room categories")
	}

	return response.Success(categories, total)
}

func (uc *roomCategoryUseCase) UpdateRoomCategory(ctx context.Context, id uint64, categoryData *entity.UpdateRoomCategory) response.StatusResponse {
	category, err := uc.repos.RoomCategory().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get room category for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("Room category with ID %d not found", id))
	}

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

func (uc *roomCategoryUseCase) DeleteRoomCategory(ctx context.Context, id uint64) response.StatusResponse {
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
