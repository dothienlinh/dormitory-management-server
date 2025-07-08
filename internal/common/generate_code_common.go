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
		code[i] = digitNumber[rand.Intn(len(digitNumber))]
	}
	return string(code)
}

func GenerateNumber() int {
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	millisStr := strconv.FormatInt(millis, 10)
	number, _ := strconv.Atoi(millisStr[len(millisStr)-6:])
	return number
}

func GenerateStudentCode() string {
	code := time.Now().Format("060102")
	code += GenerateCode(4)
	return code
}
