package entity

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusOverdue   PaymentStatus = "overdue"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentHistoryMethod string

const (
	PaymentHistoryMethodCash          PaymentHistoryMethod = "cash"
	PaymentHistoryMethodBankTransfer  PaymentHistoryMethod = "bank_transfer"
	PaymentHistoryMethodCard          PaymentHistoryMethod = "card"
	PaymentHistoryMethodMobilePayment PaymentHistoryMethod = "mobile_payment"
)

type PaymentHistory struct {
	ID          uint64                `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt        `json:"-" gorm:"index"`
	ContractID  uint64                `json:"contract_id" gorm:"not null;index"`
	Contract    Contract              `json:"contract" gorm:"foreignKey:ContractID"`
	Period      string                `json:"period" gorm:"not null"`
	Amount      float64               `json:"amount" gorm:"not null"`
	Status      PaymentStatus         `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentDate *time.Time            `json:"payment_date"`
	DueDate     time.Time             `json:"due_date" gorm:"not null"`
	Method      *PaymentHistoryMethod `json:"method" gorm:"type:varchar(20)"`
	Description *string               `json:"description" gorm:"type:text"`
	ReceiptURL  *string               `json:"receipt_url"`
}

func (PaymentHistory) TableName() string {
	return "payment_histories"
}

type CreatePaymentHistory struct {
	ContractID  uint64    `json:"contract_id" binding:"required,numeric"`
	Period      string    `json:"period" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,numeric"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Description string    `json:"description" binding:"omitempty"`
}

type UpdatePaymentHistory struct {
	Status      *PaymentStatus        `json:"status" binding:"omitempty,oneof=pending paid overdue cancelled"`
	PaymentDate *time.Time            `json:"payment_date"`
	Method      *PaymentHistoryMethod `json:"method" binding:"omitempty,oneof=cash bank_transfer card mobile_payment"`
	Description *string               `json:"description"`
	ReceiptURL  *string               `json:"receipt_url"`
}

type PaymentHistoryFilter struct {
	ContractID  uint64               `form:"contract_id"`
	Status      PaymentStatus        `form:"status" binding:"omitempty,oneof=pending paid overdue cancelled"`
	PaymentType PaymentHistoryMethod `form:"payment_type" binding:"omitempty,oneof=cash bank_transfer card mobile_payment"`
	Keyword     string               `form:"keyword"`
	Pagination
}

type PaymentHistoryDTO struct {
	ID          uint64                `json:"id"`
	ContractID  uint64                `json:"contract_id"`
	Period      string                `json:"period"`
	Amount      float64               `json:"amount"`
	Status      PaymentStatus         `json:"status"`
	PaymentDate *time.Time            `json:"payment_date"`
	DueDate     time.Time             `json:"due_date"`
	Method      *PaymentHistoryMethod `json:"method"`
	Description *string               `json:"description"`
	ReceiptURL  *string               `json:"receipt_url"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type MakePaymentRequest struct {
	Method PaymentHistoryMethod `json:"method" binding:"required,oneof=cash bank_transfer card mobile_payment"`
}
