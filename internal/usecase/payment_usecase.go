package usecase

import (
	"context"
	"dormitory_management/internal/common"
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/mq/tasks"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/pkg/logger"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/payOSHQ/payos-lib-golang"
	"go.uber.org/zap"
)

type paymentUseCase struct {
	repos       repository.Repositories
	logger      logger.Logger
	config      *config.Config
	asynqClient *asynq.Client
}

func NewPaymentUseCase(repos repository.Repositories, logger logger.Logger, asynqClient *asynq.Client) *paymentUseCase {
	cfg := config.LoadConfig()
	return &paymentUseCase{
		repos:       repos,
		logger:      logger,
		config:      cfg,
		asynqClient: asynqClient,
	}
}

func (uc *paymentUseCase) CreateLinkPaymentVietQR(ctx context.Context, payload *entity.CreateLinkPaymentVietQR) response.StatusResponse {

	user := &entity.User{ID: payload.UserID}

	if err := uc.repos.User().GetByID(ctx, user); err != nil {
		uc.logger.Error("Failed to get user by ID", zap.Error(err))
		return response.NotFound(err.Error())
	}

	expiredAt := int(time.Now().Add(15 * time.Minute).Unix())

	body := payos.CheckoutRequestType{
		OrderCode: int64(common.GenerateNumber()),
		Amount:    payload.Amount,
		Items: []payos.Item{
			{
				Name:     payload.Name,
				Price:    payload.Amount,
				Quantity: 1,
			},
		},
		Description: payload.Description,
		CancelUrl:   uc.config.Client.ClientDomain + "/payment/cancel",
		ReturnUrl:   uc.config.Client.ClientDomain + "/payment/success",
		ExpiredAt:   &expiredAt,
	}

	data, err := payos.CreatePaymentLink(body)
	if err != nil {
		uc.logger.Error("Failed to create payment link", zap.Error(err))
		return response.InternalServerError("Failed to create payment link")
	}

	jsonPayload, err := json.Marshal(&entity.Payment{
		PaymentMethod:  entity.PaymentMethodVietQR,
		PaymentChannel: entity.PaymentChannelBankTransfer,
		UserId:         payload.UserID,
		Amount:         payload.Amount,
		Currency:       data.Currency,
		Bin:            data.Bin,
		AccountNumber:  data.AccountNumber,
		AccountName:    data.AccountName,
		Description:    payload.Description,
		OrderCode:      data.OrderCode,
		PaymentLinkId:  data.PaymentLinkId,
		Status:         data.Status,
		ExpiredAt:      data.ExpiredAt,
	})
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return response.InternalServerError("Failed send mail verify account")
	}

	task := asynq.NewTask(string(tasks.TaskCreateLinkPaymentVietQR), jsonPayload)
	uc.asynqClient.EnqueueContext(ctx, task)

	return response.Success(data.CheckoutUrl, 0)
}

func (uc *paymentUseCase) ReceiveHookVietQR(ctx context.Context, webhookData *payos.WebhookDataType) error {
	if webhookData.OrderCode == 123 {
		return nil
	}
	bill := entity.Bill{}
	if err := uc.repos.Payment().ReceiveHookVietQR(ctx, webhookData, &bill); err != nil {
		uc.logger.Error("Failed to receive hook vietqr", zap.Error(err))
		return err
	}

	jsonPayload, err := json.Marshal(&bill)
	if err != nil {
		uc.logger.Error("Failed to marshal payload", zap.Error(err))
		return err
	}

	if bill.ID != 0 {
		task := asynq.NewTask(string(tasks.TaskSendEmailBill), jsonPayload)
		uc.asynqClient.EnqueueContext(ctx, task)
	}

	return nil
}
