package usecase

import (
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/mq/tasks"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type emailUseCase struct {
	logger      logger.Logger
	config      *config.Config
	asynqClient *asynq.Client
	repos       repository.Repositories
}

func NewEmailUseCase(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) usecase.EmailUseCase {
	cfg := config.LoadConfig()

	return &emailUseCase{logger: logger, config: cfg, asynqClient: asynqClient, repos: repos}
}

func (uc *emailUseCase) SendOTP(ctx context.Context, payload entity.SendCodeEmail) (*asynq.TaskInfo, error) {
	otpCode := &entity.OtpCode{
		OtpCode:        common.GenerateCode(6),
		Identifier:     payload.Email,
		IdentifierType: entity.IdentifierTypeEmail.String(),
		IsUsed:         false,
		OtpType:        entity.OtpTypeVerifyEmail.String(),
		ExpiresAt:      common.GetExpireTime(15),
	}

	if err := uc.repos.OtpCode().CreateOtpCode(ctx, otpCode); err != nil {
		uc.logger.Error("Failed to create OTP code", zap.Error(err))
		return nil, err
	}

	jsonPayload, err := json.Marshal(otpCode)
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return nil, err
	}

	task := asynq.NewTask(string(tasks.TaskSendCodeEmail), jsonPayload)

	return uc.asynqClient.EnqueueContext(ctx, task)
}

func (uc *emailUseCase) VerifyCodeEmail(ctx context.Context, payload entity.VerifyCodeEmail) response.StatusResponse {
	otpCode := &entity.OtpCode{
		Identifier: payload.Email,
		OtpCode:    payload.Code,
	}
	if err := uc.repos.OtpCode().FindCode(ctx, otpCode); err != nil {
		uc.logger.Error("Failed to find OTP code", zap.Error(err))
		return response.InternalServerError("Failed to verify OTP code")
	}

	if otpCode.ID == 0 {
		uc.logger.Error("OTP code not found")
		return response.NotFound("OTP code not found")
	}

	if otpCode.Identifier != payload.Email && otpCode.IdentifierType != entity.IdentifierTypeEmail.String() && otpCode.OtpType != entity.OtpTypeVerifyEmail.String() {
		uc.logger.Error("Invalid OTP code")
		return response.BadRequest("Invalid OTP code")
	}

	if otpCode.IsUsed {
		uc.logger.Error("OTP code already used")
		return response.BadRequest("OTP code already used")
	}

	expiresAt, err := common.ParsedTime(otpCode.ExpiresAt)
	if err != nil {
		uc.logger.Error("Failed to parse expiresAt", zap.Error(err))
		return response.InternalServerError("Failed to verify OTP code")
	}

	if time.Now().After(expiresAt) {
		uc.logger.Error("OTP code expired")
		return response.BadRequest("OTP code expired")
	}

	if err := uc.repos.OtpCode().UseOtpCode(ctx, otpCode.ID); err != nil {
		uc.logger.Error("Failed to use OTP code", zap.Error(err))
		return response.InternalServerError("Failed to verify OTP code")
	}

	return response.Success(nil, 0)
}
