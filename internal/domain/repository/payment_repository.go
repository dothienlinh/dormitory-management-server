package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"

	"github.com/payOSHQ/payos-lib-golang"
)

type PaymentRepository interface {
	CreateLinkPaymentVietQR(ctx context.Context, payload *entity.Payment) error
	ReceiveHookVietQR(ctx context.Context, webhookData *payos.WebhookDataType, bill entity.Bill) error
}
