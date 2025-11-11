package api

import (
	"github.com/DevKayoS/go-lambda/internal/controllers"
	"github.com/DevKayoS/go-lambda/internal/routes"
	"github.com/gin-gonic/gin"
)

type API struct {
	HealthController *controllers.HealthController
}

func NewAPI(hc *controllers.HealthController) *API {
	return &API{
		HealthController: hc,
	}
}

func (a *API) BindRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		routes.SetupHealthRoutes(api, a.HealthController)
	}
}
