package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

// AuthRepository defines the interface for authentication operations
type AuthRepository interface {
	// CheckTokenVersion checks the token version in the repository
	CheckTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint) (string, error)

	// SetTokenVersion sets the token version in the repository
	SetCacheTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint, tokenVersion string, expiresIn int) error

	// SetUserCache sets user data in the cache
	SetUserCache(ctx context.Context, user *entity.User, expiresIn int) error

	// GetUserCache retrieves user data from the cache
	GetUserCache(ctx context.Context, userID uint) (*entity.User, error)

	// DeleteUserCache deletes user data from the cache
	DeleteUserCache(ctx context.Context, userID uint) error

	// InvalidateToken invalidates a token in the repository
	InvalidateToken(ctx context.Context, tokenType entity.TokenType, userID uint) error

	// Register registers a new user
	Register(ctx context.Context, user *entity.User, otpCode *entity.OtpCode) error

	// Login authenticates a user and returns user data
	Login(ctx context.Context, email, password string) (*entity.User, error)
}
