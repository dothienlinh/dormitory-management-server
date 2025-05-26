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
	"time"

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
	var payload entity.SendCodeEmail
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		et.logger.Error("could not unmarshal payload:", zap.Error(err))
		return err
	}
	code := common.GenerateCode(6)

	otpCode, err := et.repos.OtpCode().FindCodeByCode(ctx, code)
	if err != nil {
		et.logger.Error("Failed to find OTP code", zap.Error(err))
		return err
	}

	var isExpired bool
	if otpCode != nil {

		expiresAt, err := common.ParsedTime(otpCode.ExpiresAt)
		if err != nil {
			et.logger.Error("Failed to parse expiration time", zap.Error(err))
			return err
		}

		isExpired = time.Now().After(expiresAt) || !otpCode.IsUsed || otpCode.Identifier != payload.Email || otpCode.IdentifierType != entity.IdentifierTypeEmail.String() || otpCode.OtpType != entity.OtpTypeVerifyEmail.String()
	}

	if otpCode == nil || isExpired {
		createOtpCode := &entity.CreateOtpCode{
			OtpCode:        code,
			IsUsed:         false,
			Identifier:     payload.Email,
			IdentifierType: entity.IdentifierTypeEmail.String(),
			OtpType:        entity.OtpTypeVerifyEmail.String(),
			ExpiresAt:      common.GetExpireTime(15), // 15 minutes expiration
		}

		if err := et.repos.OtpCode().CreateOtpCode(ctx, createOtpCode); err != nil {
			et.logger.Error("Failed to create OTP code", zap.Error(err))
			return err
		}

		to := []string{payload.Email}
		subject := "Your OTP Code"
		body := "Your OTP code is: " + code
		msg := []byte("Subject: " + subject + "\r\n" + "\r\n" + body)

		err = common.SendMail(to, msg, et.config.Email)
		if err != nil {
			et.logger.Error("Failed to send email", zap.Error(err))
			return err
		}

		return nil
	}

	return fmt.Errorf("OTP code is already used or expired for email: %s", payload.Email)
}
