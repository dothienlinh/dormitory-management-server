package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

// AuthUseCase defines the interface for authentication business logic
type AuthUseCase interface {
	// Register registers a new user
	Register(ctx context.Context, user *entity.UserRegister) response.StatusResponse

	// Login authenticates a user and returns tokens
	Login(ctx context.Context, loginData *entity.UserLogin) response.StatusResponse

	// RefreshToken refreshes an access token
	RefreshToken(ctx context.Context, refreshToken string) response.StatusResponse

	// Logout invalidates tokens
	Logout(ctx context.Context, userID uint) response.StatusResponse

	// GenerateTokens generates access and refresh tokens
	GenerateTokens(ctx context.Context, userID uint) (string, string, error)

	// Me returns the current user
	Me(ctx context.Context, userID uint) response.StatusResponse
}
