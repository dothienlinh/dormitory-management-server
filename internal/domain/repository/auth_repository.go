package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

// AuthRepository defines the interface for authentication operations
type AuthRepository interface {
	// StoreToken stores a token in the repository (could be a refresh token)
	StoreToken(ctx context.Context, userID uint, token string, expiresIn int) error

	// GetToken retrieves a token from the repository
	GetToken(ctx context.Context, token string) (uint, error)

	// InvalidateToken invalidates a token in the repository
	InvalidateToken(ctx context.Context, userID uint, token string) error

	// Register registers a new user
	Register(ctx context.Context, user *entity.User) error

	// Login authenticates a user and returns user data
	Login(ctx context.Context, email, password string) (*entity.User, error)
}
