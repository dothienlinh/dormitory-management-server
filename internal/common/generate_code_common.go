package common

import (
	"math/rand"
)

func GenerateCode(length int) string {
	digitNumber := "0123456789"
	code := make([]byte, length)
	for i := range length {
		code[i] = digitNumber[i]
	}

	for i := len(code) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		code[i], code[j] = code[j], code[i]
	}

	return string(code)
}
