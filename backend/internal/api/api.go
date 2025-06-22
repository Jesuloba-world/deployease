package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humabunrouter"
	"github.com/uptrace/bunrouter"

	"github.com/Jesuloba-world/deployease/backend/internal/api/handler"
	"github.com/Jesuloba-world/deployease/backend/internal/api/routes"
	"github.com/Jesuloba-world/deployease/backend/internal/config"
)

type API struct {
	config  *config.Config
	router  *bunrouter.Router
	humaAPI huma.API
}

func NewAPI(cfg config.Config, router *bunrouter.Router) *API {
	config := huma.DefaultConfig("DeployEase API", "1.0.0")
	config.Info.Description = "DeployEase REST API - A modern deployment automation platform that streamlines application deployment workflows."

	api := humabunrouter.New(router, config)

	apiInstance := &API{
		config:  &cfg,
		router:  router,
		humaAPI: api,
	}

	return apiInstance
}

func (a *API) GetHumaAPI() huma.API {
	return a.humaAPI
}

func (a *API) InitializeAndRegisterRoutes() {
	healthHandler := handler.NewHealthHandler("1.0.0")
	routes.RegisterHealthRoutes(a.humaAPI, healthHandler)
}
