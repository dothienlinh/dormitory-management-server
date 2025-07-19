package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
)

type PaymentHistoryRepository interface {
	Create(ctx context.Context, paymentHistory *entity.PaymentHistory) error

	GetByID(ctx context.Context, id uint) (*entity.PaymentHistory, error)

	GetByContractID(ctx context.Context, contractID uint) ([]entity.PaymentHistory, error)

	List(ctx context.Context, filter *entity.PaymentHistoryFilter) ([]entity.PaymentHistory, int64, error)

	GetMyPaymentHistory(ctx context.Context, userID uint64, filter *entity.PaymentHistoryFilter) ([]entity.PaymentHistory, int64, error)

	Update(ctx context.Context, paymentHistory *entity.PaymentHistory) error

	Delete(ctx context.Context, id uint) error
}
