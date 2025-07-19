package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type PaymentHistoryUseCase interface {
	CreatePaymentHistory(ctx context.Context, paymentHistory *entity.CreatePaymentHistory) response.StatusResponse

	GetPaymentHistoryByID(ctx context.Context, id uint) response.StatusResponse

	GetPaymentHistoryList(ctx context.Context, filter *entity.PaymentHistoryFilter) response.StatusResponse

	GetMyPaymentHistory(ctx context.Context, userID uint64, filter *entity.PaymentHistoryFilter) response.StatusResponse

	UpdatePaymentHistory(ctx context.Context, id uint, paymentHistory *entity.UpdatePaymentHistory) response.StatusResponse

	DeletePaymentHistory(ctx context.Context, id uint) response.StatusResponse

	MakePayment(ctx context.Context, paymentID uint, userID uint64, request *entity.MakePaymentRequest) response.StatusResponse

	DownloadReceipt(ctx context.Context, paymentID uint, userID uint64) response.StatusResponse
}
