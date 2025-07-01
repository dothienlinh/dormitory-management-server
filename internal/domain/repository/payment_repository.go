package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type PaymentRepository interface {
	CreateLinkPaymentVietQR(ctx context.Context, payload *entity.Payment) error
}
