package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"

	"gorm.io/gorm"
)

type facilitiesRepository struct {
	db *gorm.DB
}

func NewFacilitiesRepository(db *gorm.DB) *facilitiesRepository {
	return &facilitiesRepository{db: db}
}

func (r *facilitiesRepository) Create(ctx context.Context, facilities *entity.Facilities) error {
	return r.db.WithContext(ctx).Table(facilities.TableName()).Create(facilities).Error
}

func (r *facilitiesRepository) List(ctx context.Context, facilities *[]entity.Facilities) error {
	return r.db.WithContext(ctx).Table(entity.Facilities{}.TableName()).Find(facilities).Error
}

func (r *facilitiesRepository) Update(ctx context.Context, facilities *entity.Facilities) error {
	return r.db.WithContext(ctx).Table(facilities.TableName()).Where("id = ?", facilities.ID).Updates(facilities).Error
}

func (r *facilitiesRepository) Detail(ctx context.Context, facilities *entity.Facilities) error {
	return r.db.WithContext(ctx).Table(facilities.TableName()).First(facilities).Error
}

func (r *facilitiesRepository) Delete(ctx context.Context, facilities *entity.Facilities) error {
	return r.db.WithContext(ctx).Table(facilities.TableName()).Delete(facilities).Error
}
