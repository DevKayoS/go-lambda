package controllers

import (
	"context"
	"net/http"

	"github.com/DevKayoS/go-lambda/internal/errors"
	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/gin-gonic/gin"
)

type tokenService interface {
	Auth(ctx context.Context, gr models.GenerateTokenRequest) (string, error)
}

type TokenController struct {
	tokenService tokenService
}

func NewTokenController(tokenService tokenService) *TokenController {
	return &TokenController{tokenService}
}

func (tc *TokenController) Generate(ctx *gin.Context) {
	var body models.GenerateTokenRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(errors.Unathorized("Invalid request body"))
		return
	}

	token, err := tc.tokenService.Auth(ctx, body)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"code":   http.StatusOK,
		"msg":    "Token generated with successfully",
		"token":  token,
	})
}
