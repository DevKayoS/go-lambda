package api

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/DevKayoS/go-lambda/internal/middleware"
	"github.com/DevKayoS/go-lambda/internal/routes"
	"github.com/gin-gonic/gin"
)

type API struct {
	HealthController      *controllers.HealthController
	TokenController       *controllers.TokenController
	TransactionController *controllers.TransactionController
	UserController        *controllers.UserController
}

func NewAPI(
	hc *controllers.HealthController,
	tc *controllers.TokenController,
	tr *controllers.TransactionController,
	uc *controllers.UserController,
) *API {
	return &API{
		HealthController:      hc,
		TokenController:       tc,
		TransactionController: tr,
		UserController:        uc,
	}
}

func (a *API) BindRoutes(r *gin.Engine) {
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api/v1")

	apiPublic := api.Group("/")
	{
		routes.SetupHealthRoutes(apiPublic, a.HealthController)
		routes.SetupTokenRoutes(apiPublic, a.TokenController)
	}

	apiProtected := api.Group("/")
	apiProtected.Use(middleware.AuthMiddleware())
	{
		routes.SetupUserRoutes(apiProtected, a.UserController)
		routes.SetupTransactionRoutes(apiProtected, a.TransactionController)
	}
}
