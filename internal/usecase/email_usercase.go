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

type EmailUseCase struct {
	logger      logger.Logger
	config      *config.Config
	asynqClient *asynq.Client
	repos       repository.Repositories
}

func NewEmailUseCase(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) usecase.EmailUseCase {
	cfg := config.LoadConfig()

	return &EmailUseCase{logger: logger, config: cfg, asynqClient: asynqClient, repos: repos}
}

func (uc *EmailUseCase) SendOTP(ctx context.Context, payload entity.SendCodeEmail) (*asynq.TaskInfo, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return nil, err
	}

	task := asynq.NewTask(string(tasks.TypeSendCodeEmail), jsonPayload)

	return uc.asynqClient.EnqueueContext(ctx, task)
}

func (uc *EmailUseCase) VerifyCodeEmail(ctx context.Context, payload entity.VerifyCodeEmail) response.StatusResponse {
	otpCode, err := uc.repos.OtpCode().FindCodeByCode(ctx, payload.Code)
	if err != nil {
		uc.logger.Error("Failed to find OTP code", zap.Error(err))
		return response.InternalServerError("Failed to verify OTP code")
	}

	if otpCode == nil {
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
