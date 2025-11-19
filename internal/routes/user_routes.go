package routes

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(rg *gin.RouterGroup, uc *controllers.UserController) {
	user := rg.Group("/user")
	{
		user.POST("", uc.CreateUser)
	}
}
