package entity

import (
	"time"

	"gorm.io/gorm"
)

type ContractStatus string

const (
	ContractStatusActive    ContractStatus = "active"
	ContractStatusInactive  ContractStatus = "inactive"
	ContractStatusCancelled ContractStatus = "cancelled"
)

type PaymentCycle string

const (
	PaymentCycleMonthly   PaymentCycle = "monthly"
	PaymentCycleQuarterly PaymentCycle = "quarterly"
	PaymentCycleYearly    PaymentCycle = "yearly"
)

type Contract struct {
	ID            uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	UserID        uint64         `json:"user_id" gorm:"not null;index"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	RoomID        uint64         `json:"room_id" gorm:"not null;index"`
	Room          Room           `json:"room" gorm:"foreignKey:RoomID"`
	StartDate     time.Time      `json:"start_date" gorm:"not null"`
	EndDate       time.Time      `json:"end_date" gorm:"not null"`
	SignedDate    *time.Time     `json:"signed_date"`
	Price         float64        `json:"price" gorm:"not null"`
	DepositAmount *float64       `json:"deposit_amount"`
	PaymentCycle  *PaymentCycle  `json:"payment_cycle" gorm:"type:varchar(20);default:'monthly'"`
	Status        ContractStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	Description   string         `json:"description" gorm:"type:text"`
	Code          string         `json:"code" gorm:"uniqueIndex;not null"`

	// Relations
	ContractTerms  []ContractTerm   `json:"contract_terms" gorm:"foreignKey:ContractID"`
	PaymentHistory []PaymentHistory `json:"payment_history" gorm:"foreignKey:ContractID"`
}

func (Contract) TableName() string {
	return "contracts"
}

type CreateContract struct {
	UserID        uint64       `json:"user_id" binding:"required,numeric" gorm:"not null;index"`
	RoomID        uint64       `json:"room_id" binding:"required,numeric" gorm:"not null;index"`
	StartDate     string       `json:"start_date" binding:"required,validdate" gorm:"not null"`
	EndDate       string       `json:"end_date" binding:"required,validdate,gtefield=StartDate" gorm:"not null"`
	SignedDate    *string      `json:"signed_date" binding:"omitempty,validdate"`
	Price         float64      `json:"price" binding:"required,numeric" gorm:"not null"`
	DepositAmount *float64     `json:"deposit_amount" binding:"omitempty,numeric"`
	PaymentCycle  PaymentCycle `json:"payment_cycle" binding:"omitempty,oneof=monthly quarterly yearly"`
	Description   string       `json:"description" binding:"omitempty" gorm:"type:text"`
}

type UpdateContract struct {
	Status        ContractStatus `json:"status" binding:"required,oneof=active inactive cancelled" gorm:"type:varchar(20);not null"`
	StartDate     *time.Time     `json:"start_date" binding:"omitempty,validdate" gorm:"not null"`
	EndDate       *time.Time     `json:"end_date" binding:"omitempty,validdate,gtefield=StartDate" gorm:"not null"`
	SignedDate    *time.Time     `json:"signed_date" binding:"omitempty,validdate"`
	Price         *float64       `json:"price" binding:"omitempty,numeric" gorm:"not null"`
	DepositAmount *float64       `json:"deposit_amount" binding:"omitempty,numeric"`
	PaymentCycle  *PaymentCycle  `json:"payment_cycle" binding:"omitempty,oneof=monthly quarterly yearly"`
	Description   *string        `json:"description" binding:"omitempty" gorm:"type:text"`
}

type ContractDTO struct {
	ID             uint64           `json:"id"`
	UserID         uint64           `json:"user_id"`
	User           User             `json:"user"`
	RoomID         uint64           `json:"room_id"`
	Room           Room             `json:"room"`
	StartDate      time.Time        `json:"start_date"`
	EndDate        time.Time        `json:"end_date"`
	SignedDate     *time.Time       `json:"signed_date"`
	Price          float64          `json:"price"`
	DepositAmount  *float64         `json:"deposit_amount"`
	PaymentCycle   *PaymentCycle    `json:"payment_cycle"`
	Status         ContractStatus   `json:"status"`
	Description    string           `json:"description"`
	Code           string           `json:"code"`
	ContractTerms  []ContractTerm   `json:"contract_terms"`
	PaymentHistory []PaymentHistory `json:"payment_history"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

type ContractFilter struct {
	Status  ContractStatus `form:"status" binding:"omitempty,oneof=active inactive cancelled"`
	Keyword string         `form:"keyword"`
	Pagination
}
