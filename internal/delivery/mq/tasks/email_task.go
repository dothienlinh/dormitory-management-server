package tasks

import (
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
	"encoding/json"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type EmailTask struct {
	config *config.Config
	logger logger.Logger
	repos  repository.Repositories
}

func NewEmailTask(logger logger.Logger, config *config.Config, repos repository.Repositories) *EmailTask {
	return &EmailTask{
		config: config,
		logger: logger,
		repos:  repos,
	}
}

func (et *EmailTask) SendOTP(ctx context.Context, t *asynq.Task) error {
	var payload *entity.OtpCode
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal payload:", zap.Error(err))
		return err
	}

	to := []string{payload.Identifier}
	subject := "Your OTP Code"
	body := "Your OTP code is: " + payload.OtpCode
	msg := []byte("Subject: " + subject + "\r\n" + "\r\n" + body)

	if err := common.SendMail(to, msg, et.config.Email); err != nil {
		return err
	}

	return nil
}

func (et *EmailTask) SendVerifyAccount(ctx context.Context, t *asynq.Task) error {
	var payload entity.SendMailVerifyAccount
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal payload:", zap.Error(err))
		return err
	}

	return nil
}
