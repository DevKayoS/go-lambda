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
}

func NewAPI(hc *controllers.HealthController, tc *controllers.TokenController, tr *controllers.TransactionController) *API {
	return &API{
		HealthController:      hc,
		TokenController:       tc,
		TransactionController: tr,
	}
}

func (a *API) BindRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		routes.SetupHealthRoutes(api, a.HealthController)
		routes.SetupTokenRoutes(api, a.TokenController)

		// protected routes
		api.Use(middleware.AuthMiddleware())
		routes.SetupTransactionRoutes(api, a.TransactionController)

	}
}
