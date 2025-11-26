package controllers

import (
	"net/http"

	"github.com/DevKayoS/go-lambda/internal/errors"
	"github.com/DevKayoS/go-lambda/internal/models"
	"github.com/gin-gonic/gin"
)

type TransactionService interface {
	Create(body models.TransactionRequest) (any, error)
}

type TransactionController struct {
	service TransactionService
}

func NewTransactionController(service TransactionService) *TransactionController {
	return &TransactionController{service}
}

func (tr *TransactionController) Create(ctx *gin.Context) {
	var body models.TransactionRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Error(errors.BadRequest("Invalid request body!"))
		return
	}

	transaction, err := tr.service.Create(body)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": true,
		"code":   http.StatusCreated,
		"msg":    "Transaction registered with successfully",
		"data":   transaction,
	})
}
