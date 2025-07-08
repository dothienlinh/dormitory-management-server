package entity

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

func (t TokenType) String() string {
	return []string{"access_token", "refresh_token"}[t]
}

type VerifyAccount struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

type ResetPassword struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
	Code        string `json:"code" binding:"required,len=6"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
