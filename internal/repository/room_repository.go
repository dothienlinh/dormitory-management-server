package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/helper"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) Create(ctx context.Context, payload *entity.CreateRoom) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		roomEntity := &entity.Room{
			RoomNumber:     payload.RoomNumber,
			Status:         payload.Status,
			RoomCategoryID: payload.RoomCategoryID,
		}

		if err := tx.WithContext(ctx).Create(roomEntity).Error; err != nil {
			return fmt.Errorf("Create to create room: %w", err)
		}

		if len(payload.AmenityIDs) > 0 {
			uniqueAmenityIDs := helper.RemoveDuplicateUint64(payload.AmenityIDs)
			var existingCount int64
			if err := tx.WithContext(ctx).
				Table(entity.RoomAmenities{}.TableName()).
				Where("room_id = ? AND amenity_id IN ?", roomEntity.ID, uniqueAmenityIDs).
				Count(&existingCount).Error; err != nil {
				return fmt.Errorf("failed to check existing room amenities: %w", err)
			}

			if existingCount > 0 {
				return fmt.Errorf("some amenities are already associated with this room")
			}

			if err := r.createRoomAmenities(ctx, tx, roomEntity.ID, uniqueAmenityIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roomRepository) createRoomAmenities(ctx context.Context, tx *gorm.DB, roomID uint64, amenityIDs []uint64) error {
	var amenityCount int64
	if err := tx.WithContext(ctx).
		Table(entity.Amenity{}.TableName()).
		Where("id IN ?", amenityIDs).
		Count(&amenityCount).Error; err != nil {
		return fmt.Errorf("failed to validate amenities: %w", err)
	}

	if int64(len(amenityIDs)) != amenityCount {
		return fmt.Errorf("some amenity IDs do not exist")
	}

	roomAmenities := make([]entity.RoomAmenities, 0, len(amenityIDs))
	for _, amenityID := range amenityIDs {
		roomAmenities = append(roomAmenities, entity.RoomAmenities{
			RoomID:    roomID,
			AmenityID: amenityID,
		})
	}

	if err := tx.WithContext(ctx).Create(&roomAmenities).Error; err != nil {
		return fmt.Errorf("failed to create room amenities: %w", err)
	}

	return nil
}

func (r *roomRepository) GetByID(ctx context.Context, id uint64) (*entity.Room, error) {
	room := entity.Room{}
	if err := r.db.WithContext(ctx).Table("rooms r").
		Preload("Users").
		Preload("RoomCategory").
		Preload("RoomAmenities.Amenity").
		Preload("MaintenanceHistory").
		Select("r.*").
		Where("r.id = ?", id).
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}
	return &room, nil
}

func (r *roomRepository) List(ctx context.Context, filter *entity.RoomFilter) ([]entity.ListRooms, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var results []entity.ListRooms
	var total int64

	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	if err := r.db.WithContext(ctx).Table("rooms r").Where(whereClause, values...).
		Select("r.*, COUNT(u.id) as user_count").
		Joins("LEFT JOIN users u ON r.id = u.room_id").
		Preload("RoomCategory").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Group("r.id").
		Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get rooms: %w", err)
	}

	return results, total, nil
}

func (r *roomRepository) Update(ctx context.Context, room *entity.Room, amenityIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Table(room.TableName()).Save(room).Error; err != nil {
			return fmt.Errorf("failed to update room: %w", err)
		}

		roomAmenities := entity.RoomAmenities{
			RoomID: room.ID,
		}

		if err := tx.WithContext(ctx).Table(roomAmenities.TableName()).Unscoped().Where("room_id = ?", room.ID).Delete(&roomAmenities).Error; err != nil {
			return fmt.Errorf("failed to delete existing room amenities: %w", err)
		}

		if len(amenityIDs) > 0 {
			uniqueAmenityIDs := helper.RemoveDuplicateUint64(amenityIDs)
			if err := r.createRoomAmenities(ctx, tx, room.ID, uniqueAmenityIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roomRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Delete(&entity.Room{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}
	return nil
}

func (r *roomRepository) AmountStudentsInRoom(ctx context.Context, roomID uint64, amount *int64) error {
	return r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("room_id = ?", roomID).Count(amount).Error
}
