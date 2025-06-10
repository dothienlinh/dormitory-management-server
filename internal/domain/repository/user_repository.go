package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// List retrieves users based on filter
	List(ctx context.Context, filter *entity.UserFilter) ([]entity.User, int64, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uint) error

	// AddUserToRoom adds a user to a room
	AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) error

	// UserLeavesRoom removes a user from a room
	UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) error
}
