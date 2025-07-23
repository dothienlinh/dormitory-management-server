package repository

import (
	"context"
	"dormitory_management/internal/domain/entity"
	"errors"

	"github.com/payOSHQ/payos-lib-golang"
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

func (r *paymentRepository) ReceiveHookVietQR(ctx context.Context, webhookData *payos.WebhookDataType, bill *entity.Bill) error {
	payment := entity.Payment{PaymentLinkId: webhookData.PaymentLinkId}
	if err := r.db.WithContext(ctx).Table(payment.TableName()).Where(&payment).First(&payment).Error; err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		payment.Status = "PAID"
		if err := tx.WithContext(ctx).Table(payment.TableName()).Updates(&payment).Error; err != nil {
			return err
		}

		bill.UserID = payment.UserId
		bill.Status = entity.BillStatusPending
		bill.PaymentID = &payment.ID

		if err := tx.WithContext(ctx).Table(bill.TableName()).Where(&bill).First(&bill).Error; err != nil {
			return err
		}

		bill.Status = entity.BillStatusPaid
		if err := tx.WithContext(ctx).Table(bill.TableName()).Updates(&bill).Error; err != nil {
			return err
		}

		contract := entity.Contract{
			UserID: payment.UserId,
			Status: entity.ContractStatusActive,
		}
		if err := tx.WithContext(ctx).Table(contract.TableName()).Where(&contract).First(&contract).Error; err != nil {
			return err
		}

		paymentHistory := entity.PaymentHistory{
			ContractID:  contract.ID,
			Amount:      float64(payment.Amount),
			Status:      entity.PaymentStatusPaid,
			PaymentDate: &payment.CreatedAt,
		}
		if err := tx.WithContext(ctx).Table(paymentHistory.TableName()).Create(&paymentHistory).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *paymentRepository) CancelPaymentVietQR(ctx context.Context, paymentLinkId string) error {
	payment := entity.Payment{PaymentLinkId: paymentLinkId}
	if err := r.db.WithContext(ctx).Table(payment.TableName()).Where(&payment).First(&payment).Error; err != nil {
		return err
	}

	if payment.Status != "PENDING" {
		return errors.New("payment link is not in pending status")
	}

	payment.Status = "CANCELED"
	return r.db.WithContext(ctx).Table(payment.TableName()).Updates(&payment).Error
}
