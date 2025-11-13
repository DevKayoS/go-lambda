package controllers

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"code":   http.StatusBadRequest,
			"msg":    "Invalid request body!",
		})

		return
	}

	transaction, err := tr.service.Create(body)
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
		"msg":    "Transaction registered with successfully",
		"data":   transaction,
	})
}
