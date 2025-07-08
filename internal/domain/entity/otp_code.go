package entity

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
	Base
	UserID         *uint64 `json:"user_id"`
	Identifier     string  `json:"identifier"`
	IdentifierType string  `json:"identifier_type"`
	OtpCode        string  `json:"otp_code"`
	OtpType        string  `json:"otp_type"`
	IsUsed         bool    `json:"is_used"`
	ExpiresAt      string  `json:"expires_at"`
	VerifiedAt     *string `json:"verified_at"`
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
