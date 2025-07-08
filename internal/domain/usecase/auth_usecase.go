package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *entity.UserRegister) response.StatusResponse

	Login(ctx context.Context, loginData *entity.UserLogin) response.StatusResponse

	RefreshToken(ctx context.Context, refreshToken string) response.StatusResponse

	Logout(ctx context.Context, userID uint64) response.StatusResponse

	GenerateTokens(ctx context.Context, userID uint64) (string, string, error)

	Me(ctx context.Context, userID uint64) response.StatusResponse
	VerifyAccount(ctx context.Context, payload entity.VerifyAccount) response.StatusResponse
	ResendVerifyAccount(ctx context.Context, payload entity.SendCodeEmail) response.StatusResponse
	ForgotPassword(ctx context.Context, payload entity.SendCodeEmail) response.StatusResponse
	ResetPassword(ctx context.Context, payload entity.ResetPassword) response.StatusResponse
	ChangePassword(ctx context.Context, userID uint64, payload entity.ChangePassword) response.StatusResponse
}
