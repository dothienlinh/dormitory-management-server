package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type roomCategoryRepository struct {
	db *gorm.DB
}

func NewRoomCategoryRepository(db *gorm.DB) repository.RoomCategoryRepository {
	return &roomCategoryRepository{
		db: db,
	}
}

func (r *roomCategoryRepository) Create(ctx context.Context, createRoomCategory *entity.CreateRoomCategory) error {

	category := &entity.RoomCategory{
		Name:        createRoomCategory.Name,
		Description: createRoomCategory.Description,
		Capacity:    createRoomCategory.Capacity,
		Price:       createRoomCategory.Price,
	}

	if err := r.db.WithContext(ctx).Table(category.TableName()).Create(category).Error; err != nil {
		return fmt.Errorf("failed to create room category: %w", err)
	}
	return nil
}

func (r *roomCategoryRepository) GetByID(ctx context.Context, id uint64) (*entity.RoomCategory, error) {
	var category entity.RoomCategory
	if err := r.db.WithContext(ctx).Table(category.TableName()).Preload("Rooms").First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room category with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get room category: %w", err)
	}
	return &category, nil
}

func (r *roomCategoryRepository) List(ctx context.Context, filter *entity.RoomCategoryFilter) ([]entity.RoomCategory, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var categories []entity.RoomCategory
	var total int64

	if err := r.db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count room categories: %w", err)
	}

	if err := r.db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Where(whereClause, values...).
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get room categories: %w", err)
	}

	return categories, total, nil
}

func (r *roomCategoryRepository) Update(ctx context.Context, category *entity.RoomCategory) error {
	if err := r.db.WithContext(ctx).Table(category.TableName()).Save(category).Error; err != nil {
		return fmt.Errorf("failed to update room category: %w", err)
	}
	return nil
}

func (r *roomCategoryRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Table(entity.RoomCategory{}.TableName()).Delete(&entity.RoomCategory{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room category: %w", err)
	}
	return nil
}
