package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"

	"github.com/payOSHQ/payos-lib-golang"
)

type PaymentUseCase interface {
	CreateLinkPaymentVietQR(ctx context.Context, userID uint64, payload *entity.CreateLinkPaymentVietQR) response.StatusResponse
	ReceiveHookVietQR(ctx context.Context, webhookData *payos.WebhookDataType) error
	CancelPaymentVietQR(ctx context.Context, paymentLinkId string) response.StatusResponse
}
