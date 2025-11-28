package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DevKayoS/go-lambda/internal/errors"
	"github.com/DevKayoS/go-lambda/internal/pgstore"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(ctx context.Context, body pgstore.InsertUserParams) error
	GetMe(ctx context.Context, userID int64) (pgstore.GetUserWithPermissionsByIdRow, error)
}

type UserController struct {
	service UserService
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewUserController(service UserService) *UserController {
	return &UserController{service}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errors.BadRequest("Invalid request body!"))
		return
	}

	userParams := pgstore.InsertUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := uc.service.CreateUser(ctx, userParams)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": true,
		"code":   http.StatusCreated,
		"msg":    "Usuario criado com sucesso!",
	})
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	id, exists := ctx.Get("user_id")
	if !exists {
		ctx.Error(errors.BadRequest("user not found"))
		ctx.Abort()
		return
	}

	userID, ok := id.(int64)
	if !ok {
		ctx.Error(errors.BadRequest("invalid user_id type"))
		ctx.Abort()
		return
	}

	user, err := uc.service.GetMe(ctx, userID)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"code":   http.StatusOK,
		"msg":    "Usuario encontrado com sucesso!",
		"data":   user,
	})
}
