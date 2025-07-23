package entity

import (
	"time"

	"gorm.io/gorm"
)

type BillStatus string

const (
	BillStatusPending BillStatus = "PENDING"
	BillStatusPaid    BillStatus = "PAID"
	BillStatusOverdue BillStatus = "OVERDUE"
)

type Bill struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	UserID      uint64         `json:"user_id" gorm:"not null"`
	PaymentID   *uint64        `json:"payment_id" gorm:"not null"`
	Amount      float64        `json:"amount" gorm:"not null"`
	Status      BillStatus     `json:"status" gorm:"not null"`
	Description *string        `json:"description" gorm:"type:text"`
	DueDate     time.Time      `json:"due_date" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Payment     Payment        `json:"payment" gorm:"foreignKey:PaymentID"`
}

func (Bill) TableName() string {
	return "bills"
}

type QueryBill struct {
	Description *string     `form:"description" binding:"omitempty"`
	Status      *BillStatus `form:"status" binding:"omitempty,oneof=PENDING PAID OVERDUE"`
	Pagination
}
