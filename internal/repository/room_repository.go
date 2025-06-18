package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// roomRepository implements the repository.RoomRepository interface
type roomRepository struct {
	db *gorm.DB
}

// NewRoomRepository creates a new room repository
func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &roomRepository{
		db: db,
	}
}

// Create creates a new room
func (r *roomRepository) Create(ctx context.Context, room *entity.CreateRoom) error {
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Create(room).Error; err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

// GetByID retrieves a room by ID
func (r *roomRepository) GetByID(ctx context.Context, id uint64) (*entity.Room, error) {
	var room entity.Room
	if err := r.db.WithContext(ctx).Table(room.TableName()).Preload("RoomCategory").Preload("Users").First(&room, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}
	return &room, nil
}

// List retrieves rooms based on filter
func (r *roomRepository) List(ctx context.Context, filter *entity.RoomFilter) ([]entity.Room, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var rooms []entity.Room
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Get records with pagination
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Where(whereClause, values...).
		Preload("RoomCategory").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&rooms).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get rooms: %w", err)
	}

	return rooms, total, nil
}

// Update updates a room
func (r *roomRepository) Update(ctx context.Context, room *entity.Room) error {
	if err := r.db.WithContext(ctx).Table(room.TableName()).Save(room).Error; err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}
	return nil
}

// Delete deletes a room
func (r *roomRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Delete(&entity.Room{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}
