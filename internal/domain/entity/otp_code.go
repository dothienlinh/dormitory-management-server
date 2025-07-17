package entity

import (
	"time"

	"gorm.io/gorm"
)

type IdentifierTypeEnum int

const (
	IdentifierTypeEmail IdentifierTypeEnum = iota
	IdentifierTypePhone
)

func (i IdentifierTypeEnum) String() string {
	return []string{"email", "phone"}[i]
}

type OtpTypeEnum int

const (
	OtpTypeLogin OtpTypeEnum = iota
	OtpTypeRegister
	OtpTypeResetPassword
	OtpTypeVerifyEmail
	OtpTypeVerifyPhone
	OtpTypeTransaction
	OtpTypeVerifyAccount
	OtpTypeForgotPassword
)

func (i OtpTypeEnum) String() string {
	return []string{"login", "register", "reset_password", "verify_email", "verify_phone", "transaction", "verify_account", "forgot_password"}[i]
}

type OtpCode struct {
	ID             uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	UserID         *uint64        `json:"user_id" gorm:"index"`
	Identifier     string         `json:"identifier" gorm:"not null;index"`
	IdentifierType string         `json:"identifier_type" gorm:"not null;type:varchar(20)"`
	OtpCode        string         `json:"otp_code" gorm:"not null;index"`
	OtpType        string         `json:"otp_type" gorm:"not null;type:varchar(30)"`
	IsUsed         bool           `json:"is_used" gorm:"default:false"`
	ExpiresAt      string         `json:"expires_at" gorm:"not null"`
	VerifiedAt     *string        `json:"verified_at"`
}

func (OtpCode) TableName() string {
	return "otp_codes"
}

type CreateOtpCode struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifier_type"`
	OtpCode        string `json:"otp_code"`
	OtpType        string `json:"otp_type"`
	IsUsed         bool   `json:"is_used"`
	ExpiresAt      string `json:"expires_at"`
}

type UseOtpCode struct {
	IsUsed     bool
	VerifiedAt string
}
