package healthcheck

import (
	"context"
	"micro/transport/grpc/dependency"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Handler is a struct represent itself.
type Handler struct {
	Dependency *dependency.Dependency
}

// Check is a method returns health response.
func (h *Handler) Check(context.Context, *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	h.Dependency.Logger.Log.Info("Serving the Check request for health check")
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is a method returns health response.
func (h *Handler) Watch(r *health.HealthCheckRequest, server health.Health_WatchServer) error {
	h.Dependency.Logger.Log.Info("Serving the Watch request for health check")
	return server.Send(&health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	})
}
