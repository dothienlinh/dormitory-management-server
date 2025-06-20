package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"

	"gorm.io/gorm"
)

type maintenanceHistoryRepository struct {
	db *gorm.DB
}

func NewMaintenanceHistoryRepository(db *gorm.DB) *maintenanceHistoryRepository {
	return &maintenanceHistoryRepository{
		db: db,
	}
}

func (r *maintenanceHistoryRepository) Create(ctx context.Context, history *entity.MaintenanceHistory) error {
	return r.db.WithContext(ctx).Table(history.TableName()).Create(history).Error
}

func (r *maintenanceHistoryRepository) Update(ctx context.Context, history *entity.MaintenanceHistory) error {
	return r.db.WithContext(ctx).Table(history.TableName()).Where("id = ?", history.ID).Updates(history).Error
}

func (r *maintenanceHistoryRepository) Detail(ctx context.Context, history *entity.MaintenanceHistory) error {
	return r.db.WithContext(ctx).Table(history.TableName()).Where("id = ?", history.ID).First(history).Error
}

func (r *maintenanceHistoryRepository) Delete(ctx context.Context, history *entity.MaintenanceHistory) error {
	return r.db.WithContext(ctx).Table(history.TableName()).Where("id = ?", history.ID).Delete(history).Error
}
