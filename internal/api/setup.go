package api

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/DevKayoS/go-lambda/internal/services/transaction"
	"github.com/gin-gonic/gin"
)

func SetupAPI() *gin.Engine {
	r := gin.Default()

	healthController := controllers.NewHealthController()
	tokenController := controllers.NewTokenController()

	transactionService := transaction.NewTransactionService()
	transactionController := controllers.NewTransactionController(transactionService)

	apiHandler := NewAPI(healthController, tokenController, transactionController)
	apiHandler.BindRoutes(r)

	return r
}
