package helpers

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
)

func DecodeSessionToken(jwtToken string) (jwt.MapClaims, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err 
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid %s", "token or claims!")
	}
}

func CreateSessionToken(user auth.User) string {
	claims := jwt.MapClaims{
		"performed_by_id": user.ID,
		"token": user.Token,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	sessionToken, _ := token.SignedString(secretKey)

	return sessionToken
}
