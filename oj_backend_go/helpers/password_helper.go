package helpers

import (
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
			return ""
	}

	return string(hashedPassword)
}

func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
