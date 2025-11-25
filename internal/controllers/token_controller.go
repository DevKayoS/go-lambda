package controllers

import (
	"context"
	"net/http"

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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"code":   http.StatusUnauthorized,
			"msg":    "Invalid request body!",
		})

		return
	}

	token, err := tc.tokenService.Auth(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"code":   http.StatusUnauthorized,
			"msg":    err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"code":   http.StatusOK,
		"msg":    "Token generated with successfully",
		"token":  token,
	})
}
