package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type FacilitiesRepository interface {
	Create(ctx context.Context, facilities *entity.Facilities) error

	List(ctx context.Context, facilities *[]entity.Facilities) error

	Detail(ctx context.Context, facilities *entity.Facilities) error

	Update(ctx context.Context, facilities *entity.Facilities) error

	Delete(ctx context.Context, facilities *entity.Facilities) error
}
