package entity

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
	UserId         uint64         `json:"user_id"`
	Amount         int            `json:"amount"`
	Currency       string         `json:"currency"`
	PaymentMethod  PaymentMethod  `json:"payment_method"`
	PaymentChannel PaymentChannel `json:"payment_channel"`
	TransactionId  *string        `json:"transaction_id"`
	Bin            string         `json:"bin"`
	AccountNumber  string         `json:"accountNumber"`
	AccountName    string         `json:"accountName"`
	Description    string         `json:"description"`
	OrderCode      int64          `json:"orderCode"`
	PaymentLinkId  string         `json:"paymentLinkId"`
	Status         string         `json:"status"`
	ExpiredAt      *int           `json:"expiredAt"`
}

func (Payment) TableName() string {
	return "payments"
}

type CreateLinkPaymentVietQR struct {
	UserID      uint64 `json:"user_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Description string `json:"description" binding:"required"`
	Name        string `json:"name" binding:"required"`
}
