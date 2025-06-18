package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Encrypt(plaintext string, secretKey string) (string, error) {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)

	return encodedCipherText, nil
}

func Decrypt(encodedCipherText string, secretKey string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encodedCipherText)
	if err != nil {
		return "", err
	}

	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(cipherText), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
