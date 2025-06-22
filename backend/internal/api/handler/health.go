package handler

import (
	"context"
	"time"
)

type HealthHandler struct {
	version string
}

func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		version: version,
	}
}

type HealthInput struct{}

type HealthResponseBody struct {
	Status    string    `json:"status" doc:"Current health status of the service" enum:"healthy,unhealthy" example:"healthy"`
	TimeStamp time.Time `json:"timestamp" doc:"Timestamp when the health check was performed" format:"date-time"`
	Version   string    `json:"version" doc:"Version of the service" minLength:"1" maxLength:"50" example:"1.0.0"`
}

type HealthResponse struct {
	Body HealthResponseBody `json:"body,inline"`
}

func (h *HealthHandler) Health(ctx context.Context, input *HealthInput) (*HealthResponse, error) {
	response := &HealthResponse{
		Body: HealthResponseBody{
			Status:    "healthy",
			TimeStamp: time.Now(),
			Version:   h.version,
		},
	}

	return response, nil
}

type ReadyResponseBody struct {
	Status    string    `json:"status" doc:"Current readiness status of the service" enum:"ready,not_ready" example:"ready"`
	Timestamp time.Time `json:"timestamp" doc:"Timestamp when the readiness check was performed" format:"date-time"`
}

type ReadyResponse struct {
	Body ReadyResponseBody `json:"body,inline"`
}

type ReadyInput struct{}

func (h *HealthHandler) Ready(ctx context.Context, input *ReadyInput) (*ReadyResponse, error) {
	response := &ReadyResponse{
		Body: ReadyResponseBody{
			Status:    "ready",
			Timestamp: time.Now(),
		},
	}

	return response, nil
}
