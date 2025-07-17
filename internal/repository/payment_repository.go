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
		bill.PaymentID = payment.ID
		bill.Amount = float64(payment.Amount)
		bill.Status = "PAID"
		bill.Description = payment.Description

		findBill := entity.Bill{PaymentID: payment.ID, UserID: payment.UserId}
		if err := tx.WithContext(ctx).Table(findBill.TableName()).Where(&findBill).First(&findBill).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.WithContext(ctx).Table(bill.TableName()).Create(bill).Error; err != nil {
					return err
				}
				return nil
			}
			return err
		}

		return nil
	})
}
