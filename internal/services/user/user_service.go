package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	pgstore "github.com/DevKayoS/go-lambda/internal/pgstore"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	InsertUser(ctx context.Context, arg pgstore.InsertUserParams) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (pgstore.User, error)
}

type UserService struct {
	userRepository userRepository
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{
		userRepository: pgstore.New(pool),
	}
}

func (u *UserService) CreateUser(ctx context.Context, body pgstore.InsertUserParams) error {
	if body.Email == "" {
		return fmt.Errorf("O campo e-mail é obrigatório.")
	}

	_, err := u.userRepository.GetUserByEmail(ctx, body.Email)
	if err == nil {
		return fmt.Errorf("email já está sendo utilizado")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		// Erro inesperado
		slog.Error("Erro ao buscar usuário por email: ", err)
		return fmt.Errorf("erro ao validar email")
	}

	if body.Password == "" {
		return fmt.Errorf("A senha não pode estar vazia")
	}

	if len(body.Password) < 6 {
		return fmt.Errorf("A senha não possui um tamanho valido")
	}

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		return fmt.Errorf("Erro inesperado!")
	}

	body.Password = hashedPassword

	_, err = u.userRepository.InsertUser(ctx, body)
	if err != nil {
		return fmt.Errorf("Algo deu errado ao criar o usuario: ", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// TODO: passar para ca a logica de login
