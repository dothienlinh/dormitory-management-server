package common

import (
	"math/rand"
	"strconv"
	"time"
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

func GenerateNumber() int {
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	millisStr := strconv.FormatInt(millis, 10)
	number, _ := strconv.Atoi(millisStr[len(millisStr)-6:])
	return number
}
