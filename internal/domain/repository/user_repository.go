package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error

	GetByID(ctx context.Context, user *entity.User) error

	GetByEmail(ctx context.Context, user *entity.User) error

	List(ctx context.Context, filter *entity.UserFilter) ([]entity.User, int64, error)

	Update(ctx context.Context, user *entity.User) error

	Delete(ctx context.Context, id uint64) error

	AddUserToRoom(ctx context.Context, payload entity.AddUserToRoom) error

	UserLeavesRoom(ctx context.Context, payload entity.UserLeavesRoom) error

	CheckStudentCodeExists(ctx context.Context, studentCode string) error

	UpdateUserStatusAccount(ctx context.Context, user *entity.User) error

	GetUserByRoles(ctx context.Context, user *entity.User, roles []entity.UserRole) error
}
