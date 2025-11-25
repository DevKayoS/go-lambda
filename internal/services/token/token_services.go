package token

import (
	"context"
	"fmt"
	"time"

	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/DevKayoS/go-lambda/internal/pgstore"
	"github.com/DevKayoS/go-lambda/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository interface {
	GetUserByEmail(ctx context.Context, email string) (pgstore.User, error)
}

type TokenService struct {
	userRepository userRepository
}

func NewTokenService(pool *pgxpool.Pool) *TokenService {
	return &TokenService{
		userRepository: pgstore.New(pool),
	}
}

func (ts *TokenService) Auth(ctx context.Context, gr models.GenerateTokenRequest) (string, error) {
	user, err := ts.userRepository.GetUserByEmail(ctx, gr.Email)
	if err != nil {
		return "", fmt.Errorf("not authorized")
	}

	if !utils.CheckPasswordHash(gr.Password, user.Password) {
		return "", fmt.Errorf("not authorized")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.Name,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(models.SecretKey)
	if err != nil {
		return "", fmt.Errorf("not authorized")
	}

	return tokenString, nil
}
