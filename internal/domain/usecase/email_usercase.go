package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"

	"github.com/hibiken/asynq"
)

type EmailUseCase interface {
	SendOTP(ctx context.Context, payload entity.SendCodeEmail) (*asynq.TaskInfo, error)
	VerifyCodeEmail(ctx context.Context, payload entity.VerifyCodeEmail) response.StatusResponse
}
