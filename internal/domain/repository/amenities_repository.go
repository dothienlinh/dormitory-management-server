package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type AmenitiesRepository interface {
	Create(ctx context.Context, amenities *entity.Amenities) error

	List(ctx context.Context, amenities *[]entity.Amenities) error

	Detail(ctx context.Context, amenities *entity.Amenities) error

	Update(ctx context.Context, amenities *entity.Amenities) error

	Delete(ctx context.Context, amenities *entity.Amenities) error
}
