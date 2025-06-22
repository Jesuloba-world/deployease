package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Jesuloba-world/deployease/backend/internal/api/handler"

)

func RegisterHealthRoutes(humaAPI huma.API, healthHandler *handler.HealthHandler) {
	healthGroup := huma.NewGroup(humaAPI, "/health")

	huma.Register(healthGroup, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Health Check",
		Description: "Returns health status of the application",
		Tags:        []string{"Monitoring"},
	}, healthHandler.Health)

	huma.Register(healthGroup, huma.Operation{
		OperationID: "readiness-check",
		Method:      http.MethodGet,
		Path:        "/ready",
		Summary:     "Readiness Check",
		Description: "Returns readiness status of the application",
		Tags:        []string{"Monitoring"},
	}, healthHandler.Ready)
}
