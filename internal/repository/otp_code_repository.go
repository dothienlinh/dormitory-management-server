package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OtpCodeRepository struct {
	db *gorm.DB
}

func NewOtpCodeRepository(db *gorm.DB) *OtpCodeRepository {
	return &OtpCodeRepository{db: db}
}

func (r *OtpCodeRepository) FindCode(ctx context.Context, otpCode *entity.OtpCode) error {

	if err := r.db.WithContext(ctx).Table(otpCode.TableName()).Where(otpCode).First(&otpCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return nil
}

func (r *OtpCodeRepository) CreateOtpCode(ctx context.Context, payload *entity.OtpCode) error {
	return r.db.WithContext(ctx).Table(payload.TableName()).Create(payload).Error
}

func (r *OtpCodeRepository) UseOtpCode(ctx context.Context, otpCodeId uint) error {
	useOtpCode := entity.UseOtpCode{IsUsed: true, VerifiedAt: time.Now().Format(time.RFC3339)}
	return r.db.WithContext(ctx).Table(entity.OtpCode{}.TableName()).Where("id = ?", otpCodeId).Updates(useOtpCode).Error
}
