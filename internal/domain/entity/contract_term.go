package entity

import (
	"time"

	"gorm.io/gorm"
)

type ContractTerm struct {
	ID         uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	ContractID uint64         `json:"contract_id" gorm:"not null;index"`
	Contract   Contract       `json:"contract" gorm:"foreignKey:ContractID"`
	Title      string         `json:"title" gorm:"not null"`
	Content    string         `json:"content" gorm:"type:text;not null"`
	Order      int            `json:"order" gorm:"not null;default:0"`
}

func (ContractTerm) TableName() string {
	return "contract_terms"
}

type CreateContractTerm struct {
	ContractID uint64 `json:"contract_id" binding:"required,numeric"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Order      int    `json:"order" binding:"omitempty,numeric"`
}

type UpdateContractTerm struct {
	Title   *string `json:"title" binding:"omitempty"`
	Content *string `json:"content" binding:"omitempty"`
	Order   *int    `json:"order" binding:"omitempty,numeric"`
}

type ContractTermDTO struct {
	ID         uint64    `json:"id"`
	ContractID uint64    `json:"contract_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Order      int       `json:"order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
