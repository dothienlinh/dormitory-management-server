package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (*entity.DashboardStats, error)
}
