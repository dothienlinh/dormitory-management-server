package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// userRepository implements the repository.UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Table(user.TableName()).Preload("RoomRent.Room.RoomCategory").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Table(user.TableName()).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// List retrieves users based on filter
func (r *userRepository) List(ctx context.Context, filter *entity.UserFilter) ([]entity.User, int64, error) {
	filter.Parse()
	whereClause, values := filter.Build()

	var users []entity.User
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where(whereClause, values...).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get records with pagination
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where(whereClause, values...).
		Preload("RoomRent.Room.RoomCategory").
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.GetOffset()).
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// AddUserToRoom adds a user to a room
func (r *userRepository) AddUserToRoom(ctx context.Context, payload entity.CreateRoomRent) error {
	// Check if room exists
	var room entity.Room
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Preload("RoomCategory").Where("id = ?", payload.RoomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("room with ID %d not found", payload.RoomID)
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Check if room is full
	var countStudentsInRoom int64
	if err := r.db.WithContext(ctx).Table(entity.RoomRent{}.TableName()).Where("room_id = ?", payload.RoomID).Count(&countStudentsInRoom).Error; err != nil {
		return fmt.Errorf("failed to count students in room: %w", err)
	}

	if countStudentsInRoom >= int64(room.RoomCategory.Capacity) {
		return errors.New("room is full")
	}

	// Check if user exists and is a student
	var user entity.User
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ? AND role = ?", payload.UserID, entity.UserRoleStudent).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student with ID %d not found", payload.UserID)
		}
		return fmt.Errorf("failed to get student: %w", err)
	}

	// Check if user already has a room
	if user.RoomRentID != nil {
		return errors.New("student already has a room")
	}

	// Create room rent
	roomRent := entity.RoomRent{
		RoomID: payload.RoomID,
		UserID: payload.UserID,
		Status: payload.Status,
	}

	// Transaction to create room rent and update user
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(entity.RoomRent{}.TableName()).Create(&roomRent).Error; err != nil {
			return fmt.Errorf("failed to create room rent: %w", err)
		}

		if err := tx.Table(entity.User{}.TableName()).Where("id = ?", payload.UserID).Update("room_rent_id", roomRent.ID).Error; err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		return nil
	})
}

// RemoveUserFromRoom removes a user from a room
func (r *userRepository) RemoveUserFromRoom(ctx context.Context, payload entity.RemoveUserFromRoom) error {
	// Check if room exists
	var room entity.Room
	if err := r.db.WithContext(ctx).Table(entity.Room{}.TableName()).Where("id = ?", payload.RoomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("room with ID %d not found", payload.RoomID)
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Check if user exists and is a student
	var user entity.User
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Preload("RoomRent").Where("id = ? AND role = ?", payload.UserID, entity.UserRoleStudent).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student with ID %d not found", payload.UserID)
		}
		return fmt.Errorf("failed to get student: %w", err)
	}

	// Check if user has a room
	if user.RoomRentID == nil {
		return errors.New("student not in room")
	}

	// Transaction to remove room rent and update user
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(entity.User{}.TableName()).Where("id = ?", payload.UserID).Update("room_rent_id", nil).Error; err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		if err := tx.Table(entity.RoomRent{}.TableName()).Where("id = ?", user.RoomRentID).Delete(&entity.RoomRent{}).Error; err != nil {
			return fmt.Errorf("failed to delete room rent: %w", err)
		}

		return nil
	})
}
