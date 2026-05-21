package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/httpx"
)

type HealthData struct {
	Service   string `json:"service" example:"to-osn-backend"`
	Version   string `json:"version" example:"dev"`
	Timestamp string `json:"timestamp" example:"2026-04-17T11:43:45Z"`
}

type HealthResponse struct {
	Success   bool       `json:"success" example:"true"`
	Message   string     `json:"message" example:"service is healthy"`
	RequestID string     `json:"request_id" example:"12cd4404-a385-47fb-bc1a-eea247de4516"`
	Data      HealthData `json:"data"`
}

// healthHandler godoc
// @Summary Health check service
// @Description Mengembalikan status service, versi aplikasi, dan timestamp server UTC.
// @Tags System
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return httpx.Success(c, fiber.StatusOK, "service is healthy", fiber.Map{
			"service":   cfg.AppName,
			"version":   cfg.AppVersion,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}
