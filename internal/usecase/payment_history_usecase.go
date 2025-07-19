package usecase

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/repository"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type paymentHistoryUseCase struct {
	repositories repository.Repositories
	logger       logger.Logger
}

func NewPaymentHistoryUseCase(repositories repository.Repositories, logger logger.Logger) usecase.PaymentHistoryUseCase {
	return &paymentHistoryUseCase{
		repositories: repositories,
		logger:       logger,
	}
}

func (u *paymentHistoryUseCase) CreatePaymentHistory(ctx context.Context, createPaymentHistory *entity.CreatePaymentHistory) response.StatusResponse {
	// Validate contract exists
	_, err := u.repositories.Contract().GetByID(ctx, uint(createPaymentHistory.ContractID))
	if err != nil {
		u.logger.Error("failed to get contract", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	// Create payment history
	paymentHistory := &entity.PaymentHistory{
		ContractID:  createPaymentHistory.ContractID,
		Period:      createPaymentHistory.Period,
		Amount:      createPaymentHistory.Amount,
		Status:      entity.PaymentStatusPending,
		DueDate:     createPaymentHistory.DueDate,
		Description: &createPaymentHistory.Description,
	}

	if err := u.repositories.PaymentHistory().Create(ctx, paymentHistory); err != nil {
		u.logger.Error("failed to create payment history", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	return response.Created(paymentHistory)
}

func (u *paymentHistoryUseCase) GetPaymentHistoryByID(ctx context.Context, id uint) response.StatusResponse {
	paymentHistory, err := u.repositories.PaymentHistory().GetByID(ctx, id)
	if err != nil {
		u.logger.Error("failed to get payment history", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Payment history not found")
		}
		return response.InternalServerError("Internal server error")
	}

	return response.Success(paymentHistory, 1)
}

func (u *paymentHistoryUseCase) GetPaymentHistoryList(ctx context.Context, filter *entity.PaymentHistoryFilter) response.StatusResponse {
	paymentHistories, total, err := u.repositories.PaymentHistory().List(ctx, filter)
	if err != nil {
		u.logger.Error("failed to get payment history list", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	return response.Success(paymentHistories, total)
}

func (u *paymentHistoryUseCase) GetMyPaymentHistory(ctx context.Context, userID uint64, filter *entity.PaymentHistoryFilter) response.StatusResponse {
	paymentHistories, total, err := u.repositories.PaymentHistory().GetMyPaymentHistory(ctx, userID, filter)
	if err != nil {
		u.logger.Error("failed to get my payment history", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	return response.Success(paymentHistories, total)
}

func (u *paymentHistoryUseCase) UpdatePaymentHistory(ctx context.Context, id uint, updatePaymentHistory *entity.UpdatePaymentHistory) response.StatusResponse {
	// Get existing payment history
	paymentHistory, err := u.repositories.PaymentHistory().GetByID(ctx, id)
	if err != nil {
		u.logger.Error("failed to get payment history", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Payment history not found")
		}
		return response.InternalServerError("Internal server error")
	}

	// Update fields
	if updatePaymentHistory.Status != nil {
		paymentHistory.Status = *updatePaymentHistory.Status
	}
	if updatePaymentHistory.PaymentDate != nil {
		paymentHistory.PaymentDate = updatePaymentHistory.PaymentDate
	}
	if updatePaymentHistory.Method != nil {
		paymentHistory.Method = updatePaymentHistory.Method
	}
	if updatePaymentHistory.Description != nil {
		paymentHistory.Description = updatePaymentHistory.Description
	}
	if updatePaymentHistory.ReceiptURL != nil {
		paymentHistory.ReceiptURL = updatePaymentHistory.ReceiptURL
	}

	if err := u.repositories.PaymentHistory().Update(ctx, paymentHistory); err != nil {
		u.logger.Error("failed to update payment history", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	return response.Success(paymentHistory, 1)
}

func (u *paymentHistoryUseCase) DeletePaymentHistory(ctx context.Context, id uint) response.StatusResponse {
	// Check if payment history exists
	_, err := u.repositories.PaymentHistory().GetByID(ctx, id)
	if err != nil {
		u.logger.Error("failed to get payment history", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Payment history not found")
		}
		return response.InternalServerError("Internal server error")
	}

	if err := u.repositories.PaymentHistory().Delete(ctx, id); err != nil {
		u.logger.Error("failed to delete payment history", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	return response.Success(nil, 1)
}

func (u *paymentHistoryUseCase) MakePayment(ctx context.Context, paymentID uint, userID uint64, request *entity.MakePaymentRequest) response.StatusResponse {
	// Get payment history
	paymentHistory, err := u.repositories.PaymentHistory().GetByID(ctx, paymentID)
	if err != nil {
		u.logger.Error("failed to get payment history", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Payment not found")
		}
		return response.InternalServerError("Internal server error")
	}

	// Check if user owns this payment (through contract)
	if paymentHistory.Contract.UserID != userID {
		return response.Forbidden("You don't have permission to make this payment")
	}

	// Check if payment is already paid
	if paymentHistory.Status == entity.PaymentStatusPaid {
		return response.BadRequest("Payment already completed")
	}

	// Check if payment is cancelled
	if paymentHistory.Status == entity.PaymentStatusCancelled {
		return response.BadRequest("Payment has been cancelled")
	}

	// Update payment status
	now := time.Now()
	paymentHistory.Status = entity.PaymentStatusPaid
	paymentHistory.PaymentDate = &now
	paymentHistory.Method = &request.Method

	if err := u.repositories.PaymentHistory().Update(ctx, paymentHistory); err != nil {
		u.logger.Error("failed to update payment", zap.Error(err))
		return response.InternalServerError("Internal server error")
	}

	return response.Success(paymentHistory, 1)
}

func (u *paymentHistoryUseCase) DownloadReceipt(ctx context.Context, paymentID uint, userID uint64) response.StatusResponse {
	// Get payment history
	paymentHistory, err := u.repositories.PaymentHistory().GetByID(ctx, paymentID)
	if err != nil {
		u.logger.Error("failed to get payment history", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NotFound("Payment not found")
		}
		return response.InternalServerError("Internal server error")
	}

	// Check if user owns this payment (through contract)
	if paymentHistory.Contract.UserID != userID {
		return response.Forbidden("You don't have permission to download this receipt")
	}

	// Check if payment is paid
	if paymentHistory.Status != entity.PaymentStatusPaid {
		return response.BadRequest("Payment not completed yet")
	}

	// Check if receipt exists
	if paymentHistory.ReceiptURL == nil || *paymentHistory.ReceiptURL == "" {
		return response.NotFound("Receipt not found")
	}

	// Return receipt URL for download
	receiptData := map[string]interface{}{
		"receipt_url":  *paymentHistory.ReceiptURL,
		"payment_id":   paymentHistory.ID,
		"amount":       paymentHistory.Amount,
		"period":       paymentHistory.Period,
		"payment_date": paymentHistory.PaymentDate,
	}

	return response.Success(receiptData, 1)
}
