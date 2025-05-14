package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

// userUseCase implements the usecase.UserUseCase interface
type userUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repos repository.Repositories, logger logger.Logger) usecase.UserUseCase {
	return &userUseCase{
		repos:  repos,
		logger: logger,
	}
}

// GetUserByID retrieves a user by ID
func (uc *userUseCase) GetUserByID(ctx context.Context, id uint) response.StatusResponse {
	user, err := uc.repos.User().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", id))
	}

	return response.Success(user, 1)
}

// GetListUsers retrieves users based on filter
func (uc *userUseCase) GetListUsers(ctx context.Context, filter *entity.UserFilter) response.StatusResponse {
	users, total, err := uc.repos.User().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of users", zap.Error(err))
		return response.InternalServerError("Failed to get list of users")
	}

	return response.Success(users, total)
}

// UpdateUser updates a user
func (uc *userUseCase) UpdateUser(ctx context.Context, id uint, userData *entity.User) response.StatusResponse {
	user, err := uc.repos.User().GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", id))
	}

	// Update user fields
	if userData.FullName != "" {
		user.FullName = userData.FullName
	}
	if userData.Gender != "" {
		user.Gender = userData.Gender
	}
	if userData.Status != "" {
		user.Status = userData.Status
	}
	if userData.Phone != nil {
		user.Phone = userData.Phone
	}
	if userData.Birthday != nil {
		user.Birthday = userData.Birthday
	}
	if userData.Avatar != nil {
		user.Avatar = userData.Avatar
	}

	if err := uc.repos.User().Update(ctx, user); err != nil {
		uc.logger.Error("Failed to update user", zap.Error(err))
		return response.InternalServerError("Failed to update user")
	}

	return response.Success(user, 1)
}

// DeleteUser deletes a user
func (uc *userUseCase) DeleteUser(ctx context.Context, id uint) response.StatusResponse {
	if _, err := uc.repos.User().GetByID(ctx, id); err != nil {
		uc.logger.Error("Failed to get user for deletion", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", id))
	}

	if err := uc.repos.User().Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete user", zap.Error(err))
		return response.InternalServerError("Failed to delete user")
	}

	return response.Success("User deleted successfully", 0)
}

// AddUserToRoom adds a user to a room
func (uc *userUseCase) AddUserToRoom(ctx context.Context, payload entity.CreateRoomRent) response.StatusResponse {
	if err := uc.repos.User().AddUserToRoom(ctx, payload); err != nil {
		uc.logger.Error("Failed to add user to room", zap.Error(err))
		return response.BadRequest(err.Error())
	}

	return response.Success("User added to room successfully", 0)
}

// RemoveUserFromRoom removes a user from a room
func (uc *userUseCase) RemoveUserFromRoom(ctx context.Context, payload entity.RemoveUserFromRoom) response.StatusResponse {
	if err := uc.repos.User().RemoveUserFromRoom(ctx, payload); err != nil {
		uc.logger.Error("Failed to remove user from room", zap.Error(err))
		return response.BadRequest(err.Error())
	}

	return response.Success("User removed from room successfully", 0)
}
