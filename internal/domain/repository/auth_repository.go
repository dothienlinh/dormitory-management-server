package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type AuthRepository interface {
	CheckTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint64) (string, error)

	SetCacheTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint64, tokenVersion string, expiresIn int) error

	SetUserCache(ctx context.Context, user *entity.User, expiresIn int) error

	GetUserCache(ctx context.Context, user *entity.User) error

	DeleteUserCache(ctx context.Context, userID uint64) error

	InvalidateToken(ctx context.Context, tokenType entity.TokenType, userID uint64) error

	Register(ctx context.Context, user *entity.User, otpCode *entity.OtpCode) error

	Login(ctx context.Context, user *entity.User, loginType entity.LoginType) error

	VerifyAccount(ctx context.Context, otpCode *entity.OtpCode, user *entity.User) error

	ResetPassword(ctx context.Context, user *entity.User, otpCode *entity.OtpCode) error
}
