package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type MaintenanceHistoryRepository interface {
	Create(ctx context.Context, history *entity.MaintenanceHistory) error
	// List(ctx context.Context, histories *[]entity.MaintenanceHistory) error
	Detail(ctx context.Context, history *entity.MaintenanceHistory) error
	Update(ctx context.Context, history *entity.MaintenanceHistory) error
	Delete(ctx context.Context, history *entity.MaintenanceHistory) error
}
