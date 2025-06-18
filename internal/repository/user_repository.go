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
func (r *userRepository) GetByID(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Preload("Room.RoomCategory").First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", user.ID)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Table(user.TableName()).Where(user).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with email %s not found", user.Email)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
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
		Preload("Room.RoomCategory").
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
func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// AddUserToRoom adds a user to a room
func (r *userRepository) AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) error {
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
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("room_id = ?", payload.RoomID).Count(&countStudentsInRoom).Error; err != nil {
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
	if user.RoomID != nil {
		return errors.New("student already has a room")
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ?", payload.UserID).Update("room_id", room.ID).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UserLeavesRoom removes a user from a room
func (r *userRepository) UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) error {
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
	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ? AND role = ?", payload.UserID, entity.UserRoleStudent).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("student with ID %d not found", payload.UserID)
		}
		return fmt.Errorf("failed to get student: %w", err)
	}

	// Check if user has a room
	if user.RoomID == nil {
		return fmt.Errorf("student not in room")
	}

	if err := r.db.WithContext(ctx).Table(entity.User{}.TableName()).Where("id = ?", payload.UserID).Update("room_id", nil).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
