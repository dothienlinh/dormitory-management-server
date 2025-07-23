package entity

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentMethodVietQR PaymentMethod = "VIETQR"
	PaymentMethodMoMo   PaymentMethod = "MOMO"
	PaymentMethodVNPay  PaymentMethod = "VNPAY"
)

type PaymentChannel string

const (
	PaymentChannelCreditCard   PaymentChannel = "credit_card"
	PaymentChannelDebitCard    PaymentChannel = "debit_card"
	PaymentChannelBankTransfer PaymentChannel = "bank_transfer"
)

type Payment struct {
	ID             uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	UserId         uint64         `json:"user_id" gorm:"not null;index"`
	Amount         int            `json:"amount" gorm:"not null"`
	Currency       string         `json:"currency" gorm:"not null;default:'VND'"`
	PaymentMethod  PaymentMethod  `json:"payment_method" gorm:"type:varchar(20);not null"`
	PaymentChannel PaymentChannel `json:"payment_channel" gorm:"type:varchar(20);not null"`
	TransactionId  *string        `json:"transaction_id" gorm:"uniqueIndex"`
	Bin            string         `json:"bin" gorm:"not null"`
	AccountNumber  string         `json:"account_number" gorm:"not null"`
	AccountName    string         `json:"account_name" gorm:"not null"`
	Description    string         `json:"description" gorm:"type:text"`
	OrderCode      int64          `json:"order_code" gorm:"not null"`
	PaymentLinkId  string         `json:"payment_link_id" gorm:"not null"`
	Status         string         `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	ExpiredAt      *int           `json:"expired_at"`
}

func (Payment) TableName() string {
	return "payments"
}

type CreateLinkPaymentVietQR struct {
	Amount      int      `json:"amount" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	BillIDs     []uint64 `json:"bill_ids" binding:"required,dive,gt=0"`
}

type CreatePaymentLinkWorkerPayload struct {
	Payment Payment  `json:"payment"`
	BillIDs []uint64 `json:"bill_ids"`
}
