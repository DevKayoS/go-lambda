package controllers

import (
	"net/http"

	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/DevKayoS/go-lambda/internal/services/token"
	"github.com/gin-gonic/gin"
)

type TokenController struct{}

func NewTokenController() *TokenController {
	return &TokenController{}
}

func (tc *TokenController) Generate(ctx *gin.Context) {
	var body models.GenerateTokenRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"code":   http.StatusBadRequest,
			"msg":    "Invalid request body!",
		})

		return
	}

	token, err := token.Generate(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"code":   http.StatusBadRequest,
			"msg":    err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": true,
		"code":   http.StatusCreated,
		"msg":    "Token generated with successfully",
		"token":  token,
	})
}
