package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id uint64) response.StatusResponse

	// GetListUsers retrieves users based on filter
	GetListUsers(ctx context.Context, filter *entity.UserFilter) response.StatusResponse

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, id uint64, user *entity.User) response.StatusResponse

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uint64) response.StatusResponse

	// AddUserToRoom adds a user to a room
	AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) response.StatusResponse

	// UserLeavesRoom removes a user from a room
	UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) response.StatusResponse
}
