package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type UserUseCase interface {
	GetUserByID(ctx context.Context, id uint64) response.StatusResponse

	GetListUsers(ctx context.Context, filter *entity.UserFilter) response.StatusResponse

	UpdateUser(ctx context.Context, id uint64, user *entity.User) response.StatusResponse

	DeleteUser(ctx context.Context, id uint64) response.StatusResponse

	AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) response.StatusResponse

	UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) response.StatusResponse

	UpdateUserStatusAccount(ctx context.Context, userID uint64, statusAccount entity.StatusAccount) response.StatusResponse
}
