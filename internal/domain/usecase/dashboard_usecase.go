package usecase

import (
	"context"
	"dormitory_management/internal/domain/response"
)

type DashboardUseCase interface {
	GetStats(ctx context.Context) response.StatusResponse
}
