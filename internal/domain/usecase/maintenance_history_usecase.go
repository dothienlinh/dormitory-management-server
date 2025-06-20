package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type MaintenanceHistoryUsecase interface {
	Create(ctx context.Context, history *entity.CreateMaintenanceHistory) response.StatusResponse
	// List(ctx context.Context) ([]entity.MaintenanceHistory, response.StatusResponse)
	Detail(ctx context.Context, id uint64) response.StatusResponse
	Update(ctx context.Context, id uint64) response.StatusResponse
	Delete(ctx context.Context, id uint64) response.StatusResponse
}
