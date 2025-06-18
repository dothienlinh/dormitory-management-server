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

type amenitiesUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewAmenitiesUseCase(repos repository.Repositories, logger logger.Logger) *amenitiesUseCase {
	return &amenitiesUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *amenitiesUseCase) Create(ctx context.Context, payload *entity.CreateAmenity) response.StatusResponse {
	uc.logger.Info("Create amenity")

	amenity := &entity.Amenities{
		Name: payload.Name,
	}

	if err := uc.repos.Amenities().Create(ctx, amenity); err != nil {
		uc.logger.Error("Failed to create amenity", zap.Error(err))
		response.InternalServerError("Failed to create amenity")
	}

	return response.Created(amenity)
}

func (uc *amenitiesUseCase) List(ctx context.Context) response.StatusResponse {
	uc.logger.Info("Get list amenities")
	amenities := []entity.Amenities{}

	if err := uc.repos.Amenities().List(ctx, &amenities); err != nil {
		uc.logger.Error("Failed to get list amenities", zap.Error(err))
		response.InternalServerError("Failed to get list amenities")
	}

	return response.Success(amenities, int64(len(amenities)))
}

func (uc *amenitiesUseCase) Update(ctx context.Context, payload *entity.UpdateAmenity, id uint64) response.StatusResponse {
	amenities := &entity.Amenities{
		Base: entity.Base{ID: id},
	}

	if err := uc.repos.Amenities().Detail(ctx, amenities); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Amenity not found")
		}
		uc.logger.Error("Failed to get detail amenity", zap.Error(err))
		response.InternalServerError("Failed to get detail amenity")
	}

	amenities.Name = payload.Name

	if err := uc.repos.Amenities().Update(ctx, amenities); err != nil {
		uc.logger.Error("Failed to update amenity", zap.Error(err))
		response.InternalServerError("Failed to update amenity")
	}

	return response.Success(amenities, 0)
}

func (uc *amenitiesUseCase) Delete(ctx context.Context, id uint64) response.StatusResponse {
	amenities := &entity.Amenities{
		Base: entity.Base{ID: id},
	}

	if err := uc.repos.Amenities().Detail(ctx, amenities); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Amenity not found")
		}
		uc.logger.Error("Failed to get detail amenity", zap.Error(err))
		response.InternalServerError("Failed to get detail amenity")
	}

	if err := uc.repos.Amenities().Delete(ctx, amenities); err != nil {
		uc.logger.Error("Failed to delete amenity", zap.Error(err))
		response.InternalServerError("Failed to delete amenity")
	}

	return response.Success(nil, 0)
}
