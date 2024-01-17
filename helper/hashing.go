package helper

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHashing(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func HashEmail(email string) string {
	hasher := sha256.New()
	hasher.Write([]byte(email))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}