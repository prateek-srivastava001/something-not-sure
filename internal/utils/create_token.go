package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	Exp          time.Duration
	Email        string
	Role         string
	TokenVersion int
}

func CreateToken(payload TokenPayload) (string, error) {
	secret := os.Getenv("ACCESS_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("secret key not set for token")
	}

	claims := jwt.MapClaims{
		"exp":  time.Now().Add(payload.Exp).Unix(),
		"sub":  payload.Email,
		"role": payload.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}
