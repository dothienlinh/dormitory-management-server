package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"

	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreateLinkPaymentVietQR(ctx context.Context, payload *entity.Payment) error {
	return r.db.WithContext(ctx).Table(payload.TableName()).Create(payload).Error
}
