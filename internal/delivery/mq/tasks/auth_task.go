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

type AuthTask struct {
	config *config.Config
	logger logger.Logger
	repos  repository.Repositories
}

func NewAuthTask(logger logger.Logger, config *config.Config, repos repository.Repositories) *AuthTask {
	return &AuthTask{
		config: config,
		logger: logger,
		repos:  repos,
	}
}

func (at *AuthTask) SendEmailForgotPassword(ctx context.Context, t *asynq.Task) error {
	var payload entity.SendMailForgotPassword
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		at.logger.Error("could not unmarshal forgot password payload", zap.Error(err))
		return err
	}

	to := []string{payload.Email}
	subject := "Forgot Password Verification Code"
	body := "Your verification code is: " + payload.Code

	msg := []byte(
		"To: " + payload.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	if err := common.SendMail(to, msg, at.config.Email); err != nil {
		at.logger.Error("failed to send forgot password email", zap.Error(err))
		return err
	}

	at.logger.Info("Forgot password email sent", zap.String("email", payload.Email))

	return nil
}
