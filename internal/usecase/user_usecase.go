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

type userUseCase struct {
	repos  repository.Repositories
	logger logger.Logger
}

func NewUserUseCase(repos repository.Repositories, logger logger.Logger) usecase.UserUseCase {
	return &userUseCase{
		repos:  repos,
		logger: logger,
	}
}

func (uc *userUseCase) GetUserByID(ctx context.Context, id uint64) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: id},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", id))
	}

	return response.Success(user, 1)
}

func (uc *userUseCase) GetListUsers(ctx context.Context, filter *entity.UserFilter) response.StatusResponse {
	users, total, err := uc.repos.User().List(ctx, filter)
	if err != nil {
		uc.logger.Error("Failed to get list of users", zap.Error(err))
		return response.InternalServerError("Failed to get list of users")
	}

	return response.Success(users, total)
}

func (uc *userUseCase) UpdateUser(ctx context.Context, id uint64, userData *entity.User) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: id},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", id))
	}

	if userData.FullName != "" {
		user.FullName = userData.FullName
	}
	if userData.Gender != "" {
		user.Gender = userData.Gender
	}
	if userData.Status != "" {
		user.Status = userData.Status
	}
	if userData.Phone != "" {
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

func (uc *userUseCase) DeleteUser(ctx context.Context, id uint64) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: id},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user for deletion", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", id))
	}

	if err := uc.repos.User().Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete user", zap.Error(err))
		return response.InternalServerError("Failed to delete user")
	}

	return response.Success("User deleted successfully", 0)
}

func (uc *userUseCase) AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) response.StatusResponse {
	if err := uc.repos.User().AddUserToRoom(ctx, payload); err != nil {
		uc.logger.Error("Failed to add user to room", zap.Error(err))
		return response.BadRequest(err.Error())
	}

	return response.Success("User added to room successfully", 0)
}

func (uc *userUseCase) UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) response.StatusResponse {
	if err := uc.repos.User().UserLeavesRoom(ctx, payload); err != nil {
		uc.logger.Error("Failed to remove user from room", zap.Error(err))
		return response.BadRequest(err.Error())
	}

	return response.Success("User removed from room successfully", 0)
}

func (uc *userUseCase) UpdateUserStatusAccount(ctx context.Context, userID uint64, statusAccount entity.StatusAccount) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: userID},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user by ID", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", userID))
	}

	user.StatusAccount = statusAccount
	if err := uc.repos.User().UpdateUserStatusAccount(ctx, user); err != nil {
		uc.logger.Error("Failed to update user status account", zap.Error(err))
		return response.InternalServerError("Failed to update user status account")
	}

	return response.Success("User status account updated successfully", 0)
}

func (uc *userUseCase) UpdateMe(ctx context.Context, userID uint64, payload entity.UserUpdateMe) response.StatusResponse {
	user := &entity.User{
		Base: entity.Base{ID: userID},
	}
	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user for update", zap.Error(err))
		return response.NotFound(fmt.Sprintf("User with ID %d not found", userID))
	}

	if err := uc.repos.User().UpdateMe(ctx, user, payload); err != nil {
		uc.logger.Error("Failed to update user", zap.Error(err))
		return response.InternalServerError("Failed to update user")
	}

	return response.Success(user, 0)

}
