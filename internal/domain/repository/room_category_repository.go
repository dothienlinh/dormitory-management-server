package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type RoomCategoryRepository interface {
	Create(ctx context.Context, category *entity.CreateRoomCategory) error

	GetByID(ctx context.Context, id uint64) (*entity.RoomCategory, error)

	List(ctx context.Context, filter *entity.RoomCategoryFilter) ([]entity.RoomCategory, int64, error)

	Update(ctx context.Context, category *entity.RoomCategory) error

	Delete(ctx context.Context, id uint64) error
}
