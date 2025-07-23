package tasks

import (
	"context"
	"dormitory_management/internal/config"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/pkg/logger"
	"encoding/json"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type PaymentTask struct {
	config *config.Config
	logger logger.Logger
	repos  repository.Repositories
}

func NewPaymentTask(logger logger.Logger, config *config.Config, repos repository.Repositories) *PaymentTask {
	return &PaymentTask{
		config: config,
		logger: logger,
		repos:  repos,
	}
}

func (pt *PaymentTask) CreateLinkPaymentVietQR(ctx context.Context, t *asynq.Task) error {
	var payload *entity.CreatePaymentLinkWorkerPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		pt.logger.Error("could not unmarshal payment payload", zap.Error(err))
		return err
	}

	payment := &entity.Payment{
		Amount:         payload.Payment.Amount,
		PaymentMethod:  entity.PaymentMethodVietQR,
		Currency:       payload.Payment.Currency,
		PaymentChannel: entity.PaymentChannelBankTransfer,
		Description:    payload.Payment.Description,
		OrderCode:      payload.Payment.OrderCode,
		PaymentLinkId:  payload.Payment.PaymentLinkId,
		Status:         payload.Payment.Status,
		ExpiredAt:      payload.Payment.ExpiredAt,
		Bin:            payload.Payment.Bin,
		AccountNumber:  payload.Payment.AccountNumber,
		AccountName:    payload.Payment.AccountName,
		UserId:         payload.Payment.UserId,
	}

	if err := pt.repos.Payment().CreateLinkPaymentVietQR(ctx, payment); err != nil {
		pt.logger.Error("could not create payment", zap.Error(err))
		return err
	}

	if err := pt.repos.Bill().UpdateBillsWithPayment(ctx, payload.BillIDs, payment.ID); err != nil {
		pt.logger.Error("could not update bills with payment", zap.Error(err))
		return err
	}
	return nil
}
