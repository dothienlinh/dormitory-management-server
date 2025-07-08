package tasks

import (
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
	"encoding/json"
	"fmt"
	"net/url"

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
	var payload entity.OtpCode
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal OTP payload", zap.Error(err))
		return err
	}

	to := []string{payload.Identifier}
	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", payload.OtpCode)

	msg := []byte(
		"To: " + payload.Identifier + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	if err := common.SendMail(to, msg, et.config.Email); err != nil {
		et.logger.Error("failed to send OTP email", zap.Error(err))
		return err
	}

	et.logger.Info("OTP email sent", zap.String("email", payload.Identifier))
	return nil
}

func (et *EmailTask) SendVerifyAccount(ctx context.Context, t *asynq.Task) error {
	var payload entity.SendMailVerifyAccount
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal payload:", zap.Error(err))
		return err
	}

	escapeToken := url.QueryEscape(payload.Token)
	escapeEmail := url.QueryEscape(payload.Email)
	verifyLink := fmt.Sprintf("%s/auth/verify-account/%s/%s", et.config.Client.ClientDomain, escapeToken, escapeEmail)

	to := []string{payload.Email}
	subject := "Verify Account"
	body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello,</p>
				<p>Click the link below to verify your account:</p>
				<p><a href="%s">Verify Account</a></p>
			</body>
		</html>
	`, verifyLink)

	msg := []byte(
		"To: " + payload.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	if err := common.SendMail(to, msg, et.config.Email); err != nil {
		et.logger.Error("failed to send verification email", zap.Error(err))
		return err
	}

	return nil
}
