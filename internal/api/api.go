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
	api := r.Group("/api/v1")
	{
		api.Use(middleware.ErrorHandler())

		routes.SetupHealthRoutes(api, a.HealthController)
		routes.SetupTokenRoutes(api, a.TokenController)
		routes.SetupUserRoutes(api, a.UserController)
		// protected routes
		api.Use(middleware.AuthMiddleware())
		routes.SetupTransactionRoutes(api, a.TransactionController)

	}
}
