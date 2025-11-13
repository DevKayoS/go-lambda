package routes

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(rg *gin.RouterGroup, tr *controllers.TransactionController) {
	transaction := rg.Group("/transaction")
	{
		transaction.POST("", tr.Create)
	}
}
