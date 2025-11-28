package user

import (
	"context"
	"errors"
	"log/slog"

	errorsHandler "github.com/DevKayoS/go-lambda/internal/errors"
	pgstore "github.com/DevKayoS/go-lambda/internal/pgstore"
	"github.com/DevKayoS/go-lambda/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository interface {
	InsertUser(ctx context.Context, arg pgstore.InsertUserParams) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (pgstore.User, error)
	GetUserWithPermissionsById(ctx context.Context, id int64) (pgstore.GetUserWithPermissionsByIdRow, error)
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
		return errorsHandler.BadRequest("O campo e-mail é obrigatório.")
	}

	_, err := u.userRepository.GetUserByEmail(ctx, body.Email)
	if err == nil {
		return errorsHandler.BadRequest("email já está sendo utilizado")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Erro ao buscar usuário por email: ", err)
		return errorsHandler.BadRequest("erro ao validar email")
	}

	if body.Password == "" {
		return errorsHandler.BadRequest("A senha não pode estar vazia")
	}

	if len(body.Password) < 6 {
		return errorsHandler.BadRequest("A senha não possui um tamanho valido")
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		return errorsHandler.Internal("Erro inesperado!", err)
	}

	body.Password = hashedPassword

	_, err = u.userRepository.InsertUser(ctx, body)
	if err != nil {
		return errorsHandler.Internal("Algo deu errado ao criar o usuario", err)
	}

	return nil
}

func (u *UserService) GetMe(ctx context.Context, userID int64) (pgstore.GetUserWithPermissionsByIdRow, error) {
	slog.Info("que usuario ta vindo", userID)
	user, err := u.userRepository.GetUserWithPermissionsById(ctx, userID)
	if err != nil {
		return pgstore.GetUserWithPermissionsByIdRow{}, errorsHandler.BadRequest("user not found")
	}

	return user, nil
}
