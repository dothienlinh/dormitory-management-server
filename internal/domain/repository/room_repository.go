package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

// RoomRepository defines the interface for room data operations
type RoomRepository interface {
	// Create a new room
	Create(ctx context.Context, room *entity.CreateRoom) error

	// GetByID retrieves a room by ID
	GetByID(ctx context.Context, id uint64) (*entity.Room, error)

	// List retrieves rooms based on filter
	List(ctx context.Context, filter *entity.RoomFilter) ([]entity.Room, int64, error)

	// Update updates an existing room
	Update(ctx context.Context, room *entity.Room) error

	// Delete deletes a room by ID
	Delete(ctx context.Context, id uint64) error
}

// RoomCategoryRepository defines the interface for room category data operations
