package token

import (
	"context"
	"time"

	"github.com/DevKayoS/go-lambda/internal/errors"
	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/DevKayoS/go-lambda/internal/pgstore"
	"github.com/DevKayoS/go-lambda/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (pgstore.User, error)
	GetUserWithPermissions(ctx context.Context, email string) (pgstore.GetUserWithPermissionsRow, error)
	GetUserPermissions(ctx context.Context, email string) ([]string, error)
}

type TokenService struct {
	userRepository UserRepository
}

func NewTokenService(pool *pgxpool.Pool) *TokenService {
	return &TokenService{
		userRepository: pgstore.New(pool),
	}
}

func (ts *TokenService) Auth(ctx context.Context, gr models.GenerateTokenRequest) (string, error) {
	user, err := ts.userRepository.GetUserWithPermissions(ctx, gr.Email)
	if err != nil {
		return "", errors.Unathorized("not authorized")
	}

	if !utils.CheckPasswordHash(gr.Password, user.Password) {
		return "", errors.Unathorized("not authorized")
	}

	permissions, err := ts.userRepository.GetUserPermissions(ctx, gr.Email)
	if err != nil {
		return "", errors.Internal("something went wrong", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.ID,
		"email":       user.Email,
		"name":        user.Name,
		"role":        user.RoleName,
		"permissions": permissions,
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(models.SecretKey)
	if err != nil {
		return "", errors.Unathorized("not authorized")
	}

	return tokenString, nil
}
