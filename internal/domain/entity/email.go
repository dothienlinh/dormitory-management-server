package entity

type SendCodeEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyCodeEmail struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
