package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type RoomCategoryRepository interface {
	// Create a new room category
	Create(ctx context.Context, category *entity.CreateRoomCategory) error

	// GetByID retrieves a room category by ID
	GetByID(ctx context.Context, id uint64) (*entity.RoomCategory, error)

	// List retrieves room categories based on filter
	List(ctx context.Context, filter *entity.RoomCategoryFilter) ([]entity.RoomCategory, int64, error)

	// Update updates an existing room category
	Update(ctx context.Context, category *entity.RoomCategory) error

	// Delete deletes a room category by ID
	Delete(ctx context.Context, id uint64) error
}
