package api

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupAPI() *gin.Engine {
	r := gin.Default()

	healthController := controllers.NewHealthController()
	tokenController := controllers.NewTokenController()

	apiHandler := NewAPI(healthController, tokenController)
	apiHandler.BindRoutes(r)

	return r
}
