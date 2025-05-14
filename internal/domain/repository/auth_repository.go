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
	SetTokenVersion(ctx context.Context, tokenType entity.TokenType, userID uint, tokenVersion string, expiresIn int) error

	// InvalidateToken invalidates a token in the repository
	InvalidateToken(ctx context.Context, tokenType entity.TokenType, userID uint) error

	// Register registers a new user
	Register(ctx context.Context, user *entity.User) error

	// Login authenticates a user and returns user data
	Login(ctx context.Context, email, password string) (*entity.User, error)
}
