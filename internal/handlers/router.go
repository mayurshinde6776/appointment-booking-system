package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter wires all routes and middleware onto a new Gin engine.
func SetupRouter(
	health *HealthHandler,
	appointment *AppointmentHandler,
) *gin.Engine {
	r := gin.New()

	// ── Global middleware ─────────────────────────────────────────────────────
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// ── System routes ─────────────────────────────────────────────────────────
	r.GET("/health", health.Check)

	// ── 404 handler ───────────────────────────────────────────────────────────
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "route not found",
		})
	})

	// ── API v1 group ──────────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")
	{
		appointments := v1.Group("/appointments")
		{
			appointments.POST("", appointment.Create)
			appointments.GET("", appointment.List)
			appointments.GET("/:id", appointment.GetByID)
			appointments.PUT("/:id", appointment.Update)
			appointments.DELETE("/:id", appointment.Delete)
		}
	}

	return r
}

// corsMiddleware adds permissive CORS headers.
// Tighten AllowOrigins in production.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
