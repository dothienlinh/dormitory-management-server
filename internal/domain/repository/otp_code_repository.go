package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type OtpCodeRepository interface {
	FindCode(ctx context.Context, otpCode *entity.OtpCode) error
	CreateOtpCode(ctx context.Context, otpCode *entity.OtpCode) error
	UseOtpCode(ctx context.Context, otpCodeId uint) error
}
