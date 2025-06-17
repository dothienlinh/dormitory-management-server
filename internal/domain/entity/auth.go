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
