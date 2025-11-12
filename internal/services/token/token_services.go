package token

import (
	"errors"
	"os"
	"time"

	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/golang-jwt/jwt"
)

var secretKey = func() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		return []byte("secretKey")
	}
	return []byte(key)
}()

func Generate(gr models.GenerateRequest) (string, error) {
	if len(gr.Password) < 3 {
		return "", errors.New("password is too short")
	}

	if len(gr.Password) > 72 {
		return "", errors.New("password is too long")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": gr.User,
		"pass": gr.Password,
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
