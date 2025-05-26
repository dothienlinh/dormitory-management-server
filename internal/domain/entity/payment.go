package entity

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
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
	Base
	UserId         uint           `json:"user_id"`
	Amount         int            `json:"amount"`
	Currency       string         `json:"currency"`
	PaymentStatus  PaymentStatus  `json:"payment_status"`
	PaymentMethod  PaymentMethod  `json:"payment_method"`
	PaymentChannel PaymentChannel `json:"payment_channel"`
	TransactionId  string         `json:"transaction_id"`
}

func (Payment) TableName() string {
	return "payments"
}
