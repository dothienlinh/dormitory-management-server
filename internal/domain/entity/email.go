package entity

type SendCodeEmail struct {
	Email string `json:"email" binding:"required,email"`
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
