package handlers

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zercle/gofiber-skeleton/pkg/logger"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks,omitempty"`
	Metadata  HealthMetadata    `json:"metadata"`
}

// HealthMetadata contains additional health information
type HealthMetadata struct {
	GoVersion   string `json:"go_version"`
	NumGoroutines int  `json:"num_goroutines"`
	MemStats    MemoryStats `json:"memory"`
}

// MemoryStats contains memory usage information
type MemoryStats struct {
	Alloc      uint64 `json:"alloc_mb"`
	TotalAlloc uint64 `json:"total_alloc_mb"`
	Sys        uint64 `json:"sys_mb"`
	NumGC      uint32 `json:"num_gc"`
}

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// Health handles the main health check endpoint
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Check system health
	status := "healthy"
	checks := make(map[string]string)

	// Add application-specific health checks here
	if err := h.checkDatabase(ctx); err != nil {
		status = "unhealthy"
		checks["database"] = fmt.Sprintf("error: %v", err)
	} else {
		checks["database"] = "ok"
	}

	// Get memory statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Version:   h.version,
		Checks:    checks,
		Metadata: HealthMetadata{
			GoVersion:    runtime.Version(),
			NumGoroutines: runtime.NumGoroutine(),
			MemStats: MemoryStats{
				Alloc:      bToMb(m.Alloc),
				TotalAlloc: bToMb(m.TotalAlloc),
				Sys:        bToMb(m.Sys),
				NumGC:      m.NumGC,
			},
		},
	}

	// Log health check access
	logger.Info("Health check accessed",
		"status", status,
		"ip", c.IP(),
		"user_agent", c.Get("User-Agent"),
	)

	// Set appropriate status code
	if status == "healthy" {
		return c.Status(fiber.StatusOK).JSON(response)
	}
	return c.Status(fiber.StatusServiceUnavailable).JSON(response)
}

// Liveness handles the liveness probe (Kubernetes style)
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	response := map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now(),
		"uptime":    time.Since(h.startTime).String(),
	}

	logger.Debug("Liveness probe accessed", "ip", c.IP())
	return c.JSON(response)
}

// Readiness handles the readiness probe (Kubernetes style)
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Check if application is ready to serve traffic
	ready := true
	checks := make(map[string]string)

	// Check database connectivity
	if err := h.checkDatabase(ctx); err != nil {
		ready = false
		checks["database"] = fmt.Sprintf("error: %v", err)
	} else {
		checks["database"] = "ok"
	}

	response := map[string]interface{}{
		"ready":     ready,
		"timestamp": time.Now(),
		"checks":    checks,
	}

	statusCode := fiber.StatusOK
	if !ready {
		statusCode = fiber.StatusServiceUnavailable
	}

	logger.Debug("Readiness probe accessed", "ready", ready, "ip", c.IP())
	return c.Status(statusCode).JSON(response)
}

// checkDatabase is a placeholder for database health checking
// Replace with actual database connectivity check
func (h *HealthHandler) checkDatabase(ctx context.Context) error {
	// TODO: Implement actual database health check
	// For now, we'll return nil to indicate everything is okay
	return nil
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// HealthCheckConfig represents configuration for health checks
type HealthCheckConfig struct {
	EnableDetailedChecks bool   `json:"enable_detailed_checks"`
	DatabaseURL         string `json:"database_url,omitempty"`
	CheckInterval       string `json:"check_interval,omitempty"`
}