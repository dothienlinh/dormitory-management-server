package entity

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

func (t TokenType) String() string {
	return []string{"access_token", "refresh_token"}[t]
}
