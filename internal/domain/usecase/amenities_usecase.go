package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type AmenitiesUseCase interface {
	Create(ctx context.Context, payload *entity.CreateAmenity) response.StatusResponse

	List(ctx context.Context) response.StatusResponse

	// Detail(ctx context.Context, id uint) response.StatusResponse

	Update(ctx context.Context, payload *entity.UpdateAmenity, id uint64) response.StatusResponse

	Delete(ctx context.Context, id uint64) response.StatusResponse
}
