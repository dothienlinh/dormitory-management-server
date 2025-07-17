package entity

type SendCodeEmail struct {
	Email string   `json:"email" binding:"required,email"`
	Type  UserRole `json:"type" binding:"required,oneof=student manager"`
}

type VerifyCodeEmail struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type SendMailVerifyAccount struct {
	UserID uint64 `json:"user_id" binding:"required,numeric"`
	Email  string `json:"email" binding:"required,email"`
	Token  string `json:"token" binding:"required"`
}

type SendMailForgotPassword struct {
	UserID uint64 `json:"user_id" binding:"required,numeric"`
	Email  string `json:"email" binding:"required,email"`
	Code   string `json:"code" binding:"required,len=6"`
}

type CompanyInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type BillEmailData struct {
	Bill    Bill        `json:"bill"`
	Company CompanyInfo `json:"company"`
}
