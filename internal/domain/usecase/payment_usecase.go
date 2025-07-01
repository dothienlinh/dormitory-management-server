package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
)

type PaymentUseCase interface {
	CreateLinkPaymentVietQR(ctx context.Context, payload *entity.CreateLinkPaymentVietQR) response.StatusResponse
}
