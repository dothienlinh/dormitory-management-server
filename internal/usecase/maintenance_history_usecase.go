package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/pkg/logger"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type maintenanceHistoryUsecase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewMaintenanceHistoryUsecase(repos repository.Repositories, logger logger.Logger) *maintenanceHistoryUsecase {
	return &maintenanceHistoryUsecase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *maintenanceHistoryUsecase) Create(ctx context.Context, payload *entity.CreateMaintenanceHistory) response.StatusResponse {
	historyEntity := &entity.MaintenanceHistory{
		RoomID:          payload.RoomID,
		Description:     payload.Description,
		MaintenanceDate: payload.MaintenanceDate,
		Cost:            payload.Cost,
	}

	if err := uc.repos.MaintenanceHistory().Create(ctx, historyEntity); err != nil {
		uc.logger.Error("Failed to create maintenance history", zap.Error(err))
		return response.InternalServerError("Failed to create maintenance history")
	}

	return response.Created(historyEntity)
}

func (uc *maintenanceHistoryUsecase) Detail(ctx context.Context, id uint64) response.StatusResponse {
	history := &entity.MaintenanceHistory{Base: entity.Base{ID: id}}

	if err := uc.repos.MaintenanceHistory().Detail(ctx, history); err != nil {
		uc.logger.Error("Failed to get maintenance history detail", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Maintenance history not found")
		}
		return response.InternalServerError("Failed to get maintenance history detail")
	}

	return response.Success(history, 0)
}

func (uc *maintenanceHistoryUsecase) Update(ctx context.Context, id uint64) response.StatusResponse {
	history := &entity.MaintenanceHistory{Base: entity.Base{ID: id}}

	if err := uc.repos.MaintenanceHistory().Detail(ctx, history); err != nil {
		uc.logger.Error("Failed to get maintenance history detail", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Maintenance history not found")
		}
		return response.InternalServerError("Failed to get maintenance history detail")
	}

	if err := uc.repos.MaintenanceHistory().Update(ctx, history); err != nil {
		uc.logger.Error("Failed to update maintenance history", zap.Error(err))
		return response.InternalServerError("Failed to update maintenance history")
	}

	return response.Success(history, 0)
}

func (uc *maintenanceHistoryUsecase) Delete(ctx context.Context, id uint64) response.StatusResponse {
	history := &entity.MaintenanceHistory{Base: entity.Base{ID: id}}

	if err := uc.repos.MaintenanceHistory().Detail(ctx, history); err != nil {
		uc.logger.Error("Failed to get maintenance history detail", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Maintenance history not found")
		}
		return response.InternalServerError("Failed to get maintenance history detail")
	}

	if err := uc.repos.MaintenanceHistory().Delete(ctx, history); err != nil {
		uc.logger.Error("Failed to delete maintenance history", zap.Error(err))
		return response.InternalServerError("Failed to delete maintenance history")
	}

	return response.Success(nil, 0)
}
