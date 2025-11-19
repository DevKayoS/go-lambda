package api

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/DevKayoS/go-lambda/internal/pgstore/database"
	"github.com/DevKayoS/go-lambda/internal/services/transaction"
	"github.com/DevKayoS/go-lambda/internal/services/user"
	"github.com/gin-gonic/gin"
)

func SetupAPI() *gin.Engine {
	r := gin.Default()

	healthController := controllers.NewHealthController()
	tokenController := controllers.NewTokenController()

	transactionService := transaction.NewTransactionService()
	transactionController := controllers.NewTransactionController(transactionService)

	userService := user.NewUserService(database.Pool)
	userController := controllers.NewUserController(userService)

	apiHandler := NewAPI(healthController, tokenController, transactionController, userController)
	apiHandler.BindRoutes(r)

	return r
}
