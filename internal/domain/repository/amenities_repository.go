package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type AmenitiesRepository interface {
	Create(ctx context.Context, amenities *entity.Amenity) error

	List(ctx context.Context, amenities *[]entity.Amenity) error

	Detail(ctx context.Context, amenities *entity.Amenity) error

	Update(ctx context.Context, amenities *entity.Amenity) error

	Delete(ctx context.Context, amenities *entity.Amenity) error
}
