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
func (r *roomRepository) Create(ctx context.Context, room *entity.Room) error {
	if err := r.db.WithContext(ctx).Table(room.TableName()).Create(room).Error; err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	return nil
}

// GetByID retrieves a room by ID
func (r *roomRepository) GetByID(ctx context.Context, id uint) (*entity.Room, error) {
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
func (r *roomRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Delete(&entity.Room{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

// roomCategoryRepository implements the repository.RoomCategoryRepository interface
type roomCategoryRepository struct {
	db *gorm.DB
}

// NewRoomCategoryRepository creates a new room category repository
func NewRoomCategoryRepository(db *gorm.DB) repository.RoomCategoryRepository {
	return &roomCategoryRepository{
		db: db,
	}
}

// Create creates a new room category
func (r *roomCategoryRepository) Create(ctx context.Context, category *entity.RoomCategory) error {
	if err := r.db.WithContext(ctx).Table(category.TableName()).Create(category).Error; err != nil {
		return fmt.Errorf("failed to create room category: %w", err)
	}
	return nil
}

// GetByID retrieves a room category by ID
func (r *roomCategoryRepository) GetByID(ctx context.Context, id uint) (*entity.RoomCategory, error) {
	var category entity.RoomCategory
	if err := r.db.WithContext(ctx).Table(category.TableName()).Preload("Rooms").First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room category with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get room category: %w", err)
	}
	return &category, nil
}

// List retrieves room categories based on filter
func (r *roomCategoryRepository) List(ctx context.Context, filter *entity.RoomCategoryFilter) ([]entity.RoomCategory, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var categories []entity.RoomCategory
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count room categories: %w", err)
	}

	// Get records with pagination
	if err := r.db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Where(whereClause, values...).
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get room categories: %w", err)
	}

	return categories, total, nil
}

// Update updates a room category
func (r *roomCategoryRepository) Update(ctx context.Context, category *entity.RoomCategory) error {
	if err := r.db.WithContext(ctx).Table(category.TableName()).Save(category).Error; err != nil {
		return fmt.Errorf("failed to update room category: %w", err)
	}
	return nil
}

// Delete deletes a room category
func (r *roomCategoryRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Delete(&entity.RoomCategory{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room category: %w", err)
	}
	return nil
}
