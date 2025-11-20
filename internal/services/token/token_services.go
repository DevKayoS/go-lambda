package token

import (
	"errors"
	"time"

	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/golang-jwt/jwt"
)

// TODO: ajustar esse metodo para ficar igual aos outros para ele ter a dependencia do userRepository
func Generate(gr models.GenerateTokenRequest) (string, error) {
	if len(gr.Password) < 3 {
		return "", errors.New("password is too short")
	}

	if len(gr.Password) > 72 {
		return "", errors.New("password is too long")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": gr.Email,
		"pass": gr.Password,
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(models.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
