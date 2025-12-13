package user

import (
	"context"
	"errors"
	"testing"

	"github.com/DevKayoS/go-lambda/internal/pgstore"
	"github.com/DevKayoS/go-lambda/internal/services/user/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	body := pgstore.InsertUserParams{
		Email:    "test@example.com",
		Password: "senha123",
	}

	mockRepo.On("GetUserByEmail", ctx, body.Email).Return(pgstore.User{}, pgx.ErrNoRows)
	mockRepo.On("InsertUser", ctx, mock.AnythingOfType("pgstore.InsertUserParams")).Return(int64(1), nil)

	err := service.CreateUser(ctx, body)

	assert.NoError(t, err)
}

func TestCreateUser_EmailEmpty(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	body := pgstore.InsertUserParams{
		Email:    "",
		Password: "senha123",
	}

	err := service.CreateUser(ctx, body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "e-mail é obrigatório")
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	body := pgstore.InsertUserParams{
		Email:    "existing@example.com",
		Password: "senha123",
	}

	existingUser := pgstore.User{
		ID:    1,
		Email: "existing@example.com",
	}

	mockRepo.On("GetUserByEmail", ctx, body.Email).Return(existingUser, nil)

	err := service.CreateUser(ctx, body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email já está sendo utilizado")
}

func TestCreateUser_PasswordEmpty(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	body := pgstore.InsertUserParams{
		Email:    "test@example.com",
		Password: "",
	}

	mockRepo.On("GetUserByEmail", ctx, body.Email).Return(pgstore.User{}, pgx.ErrNoRows)

	err := service.CreateUser(ctx, body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "senha não pode estar vazia")
}

func TestCreateUser_PasswordTooShort(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	body := pgstore.InsertUserParams{
		Email:    "test@example.com",
		Password: "123",
	}

	mockRepo.On("GetUserByEmail", ctx, body.Email).Return(pgstore.User{}, pgx.ErrNoRows)

	err := service.CreateUser(ctx, body)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tamanho valido")
}

func TestCreateUser_DatabaseError(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	body := pgstore.InsertUserParams{
		Email:    "test@example.com",
		Password: "senha123",
	}

	mockRepo.On("GetUserByEmail", ctx, body.Email).Return(pgstore.User{}, pgx.ErrNoRows)
	mockRepo.On("InsertUser", ctx, mock.AnythingOfType("pgstore.InsertUserParams")).
		Return(int64(0), errors.New("database error"))

	err := service.CreateUser(ctx, body)

	assert.Error(t, err)
}

// =============================================================================
// TESTES - GetMe
// =============================================================================

func TestGetMe_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	expectedUser := pgstore.GetUserWithPermissionsByIdRow{
		ID:    1,
		Email: "test@example.com",
	}

	mockRepo.On("GetUserWithPermissionsById", ctx, int64(1)).Return(expectedUser, nil)

	user, err := service.GetMe(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetMe_UserNotFound(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("GetUserWithPermissionsById", ctx, int64(999)).
		Return(pgstore.GetUserWithPermissionsByIdRow{}, pgx.ErrNoRows)

	user, err := service.GetMe(ctx, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Empty(t, user)
}

// =============================================================================
// TESTES - List
// =============================================================================

func TestList_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	expectedUsers := []pgstore.ListUserRow{
		{Name: "user1", Email: "user1@example.com"},
		{Name: "user2", Email: "user2@example.com"},
	}

	mockRepo.On("ListUser", ctx).Return(expectedUsers, nil)

	users, err := service.List(ctx)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers, users)
}

func TestList_EmptyList(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("ListUser", ctx).Return([]pgstore.ListUserRow{}, nil)

	users, err := service.List(ctx)

	assert.Error(t, err)
	assert.Empty(t, users)
}

func TestList_DatabaseError(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := &UserService{userRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("ListUser", ctx).Return([]pgstore.ListUserRow{}, errors.New("db error"))

	users, err := service.List(ctx)

	assert.Error(t, err)
	assert.Empty(t, users)
}
