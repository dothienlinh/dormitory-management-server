package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type OtpCodeRepository interface {
	FindCodeByCode(ctx context.Context, otpCode string) (*entity.OtpCode, error)
	CreateOtpCode(ctx context.Context, otpCode *entity.CreateOtpCode) error
	UseOtpCode(ctx context.Context, otpCodeId uint) error
}
