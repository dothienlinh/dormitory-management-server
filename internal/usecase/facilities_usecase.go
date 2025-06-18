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

type facilitiesUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewFacilitiesUseCase(repos repository.Repositories, logger logger.Logger) *facilitiesUseCase {
	return &facilitiesUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *facilitiesUseCase) Create(ctx context.Context, payload *entity.CreateFacility) response.StatusResponse {
	uc.logger.Info("Create facility")

	facility := &entity.Facilities{
		Name: payload.Name,
	}

	if err := uc.repos.Facilities().Create(ctx, facility); err != nil {
		uc.logger.Error("Failed to create facility", zap.Error(err))
		response.InternalServerError("Failed to create facility")
	}

	return response.Created(facility)
}

func (uc *facilitiesUseCase) List(ctx context.Context) response.StatusResponse {
	uc.logger.Info("Get list facilities")
	facilities := []entity.Facilities{}

	if err := uc.repos.Facilities().List(ctx, &facilities); err != nil {
		uc.logger.Error("Failed to get list facilities", zap.Error(err))
		response.InternalServerError("Failed to get list facilities")
	}

	return response.Success(facilities, int64(len(facilities)))
}

func (uc *facilitiesUseCase) Update(ctx context.Context, payload *entity.UpdateFacility, id uint64) response.StatusResponse {
	facilities := &entity.Facilities{
		Base: entity.Base{ID: id},
	}

	if err := uc.repos.Facilities().Detail(ctx, facilities); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Facility not found")
		}
		uc.logger.Error("Failed to get detail facility", zap.Error(err))
		response.InternalServerError("Failed to get detail facilitys")
	}

	facilities.Name = payload.Name

	if err := uc.repos.Facilities().Update(ctx, facilities); err != nil {
		uc.logger.Error("Failed to update facility", zap.Error(err))
		response.InternalServerError("Failed to update facility")
	}

	return response.Success(facilities, 0)
}

func (uc *facilitiesUseCase) Delete(ctx context.Context, id uint64) response.StatusResponse {
	facilities := &entity.Facilities{
		Base: entity.Base{ID: id},
	}

	if err := uc.repos.Facilities().Detail(ctx, facilities); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Facility not found")
		}
		uc.logger.Error("Failed to get detail facility", zap.Error(err))
		response.InternalServerError("Failed to get detail facilitys")
	}

	if err := uc.repos.Facilities().Delete(ctx, facilities); err != nil {
		uc.logger.Error("Failed to delete facility", zap.Error(err))
		response.InternalServerError("Failed to delete facility")
	}

	return response.Success(nil, 0)
}
