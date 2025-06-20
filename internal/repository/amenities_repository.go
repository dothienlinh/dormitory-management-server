package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"

	"gorm.io/gorm"
)

type amenitiesRepository struct {
	db *gorm.DB
}

func NewAmenitiesRepository(db *gorm.DB) *amenitiesRepository {
	return &amenitiesRepository{db: db}
}

func (r *amenitiesRepository) Create(ctx context.Context, amenities *entity.Amenity) error {
	return r.db.WithContext(ctx).Table(amenities.TableName()).Create(amenities).Error
}

func (r *amenitiesRepository) List(ctx context.Context, amenities *[]entity.Amenity) error {
	return r.db.WithContext(ctx).Table(entity.Amenity{}.TableName()).Find(amenities).Error
}

func (r *amenitiesRepository) Update(ctx context.Context, amenities *entity.Amenity) error {
	return r.db.WithContext(ctx).Table(amenities.TableName()).Where("id = ?", amenities.ID).Updates(amenities).Error
}

func (r *amenitiesRepository) Detail(ctx context.Context, amenities *entity.Amenity) error {
	return r.db.WithContext(ctx).Table(amenities.TableName()).Preload("RoomAmenities").First(amenities).Error
}

func (r *amenitiesRepository) Delete(ctx context.Context, amenities *entity.Amenity) error {
	return r.db.WithContext(ctx).Table(amenities.TableName()).Delete(amenities).Error
}
