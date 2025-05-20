package usecase

import (
	"context"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"

	"go.uber.org/zap"
)

type dashboardUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewDashboardUseCase(repos repository.Repositories, logger logger.Logger) usecase.DashboardUseCase {
	return &dashboardUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *dashboardUseCase) GetStats(ctx context.Context) response.StatusResponse {
	stats, err := uc.repos.Dashboard().GetStats(ctx)
	if err != nil {
		uc.logger.Error("Failed to get stats", zap.Error(err))
		return response.BadRequest("Failed to get stats")
	}

	return response.Success(stats, 1)
}
