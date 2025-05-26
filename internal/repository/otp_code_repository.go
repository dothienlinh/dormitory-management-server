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

func (r *OtpCodeRepository) FindCodeByCode(ctx context.Context, otpCode string) (*entity.OtpCode, error) {
	var code entity.OtpCode

	if err := r.db.WithContext(ctx).Table(code.TableName()).Where(&entity.OtpCode{OtpCode: otpCode}).First(&code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &code, nil
}

func (r *OtpCodeRepository) CreateOtpCode(ctx context.Context, otpCode *entity.CreateOtpCode) error {
	return r.db.WithContext(ctx).Table(entity.OtpCode{}.TableName()).Create(otpCode).Error
}

func (r *OtpCodeRepository) UseOtpCode(ctx context.Context, otpCodeId uint) error {
	useOtpCode := entity.UseOtpCode{IsUsed: true, VerifiedAt: time.Now().Format(time.RFC3339)}
	return r.db.WithContext(ctx).Table(entity.OtpCode{}.TableName()).Where("id = ?", otpCodeId).Updates(useOtpCode).Error
}
