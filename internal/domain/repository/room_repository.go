package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

// RoomRepository defines the interface for room data operations
type RoomRepository interface {
	// Create a new room
	Create(ctx context.Context, room *entity.Room) error

	// GetByID retrieves a room by ID
	GetByID(ctx context.Context, id uint) (*entity.Room, error)

	// List retrieves rooms based on filter
	List(ctx context.Context, filter *entity.RoomFilter) ([]entity.Room, int64, error)

	// Update updates an existing room
	Update(ctx context.Context, room *entity.Room) error

	// Delete deletes a room by ID
	Delete(ctx context.Context, id uint) error
}

// RoomCategoryRepository defines the interface for room category data operations
type RoomCategoryRepository interface {
	// Create a new room category
	Create(ctx context.Context, category *entity.RoomCategory) error

	// GetByID retrieves a room category by ID
	GetByID(ctx context.Context, id uint) (*entity.RoomCategory, error)

	// List retrieves room categories based on filter
	List(ctx context.Context, filter *entity.RoomCategoryFilter) ([]entity.RoomCategory, int64, error)

	// Update updates an existing room category
	Update(ctx context.Context, category *entity.RoomCategory) error

	// Delete deletes a room category by ID
	Delete(ctx context.Context, id uint) error
}
