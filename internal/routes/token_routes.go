package routes

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupTokenRoutes(rg *gin.RouterGroup, tc *controllers.TokenController) {
	token := rg.Group("/auth")
	{
		token.POST("", tc.Generate)
	}
}
