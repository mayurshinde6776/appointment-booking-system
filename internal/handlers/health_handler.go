package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler exposes the /health endpoint.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler constructs a HealthHandler that pings the database.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse is the JSON shape returned by the health check.
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// Check handles GET /health.
//
//	@Summary      Health check
//	@Description  Returns the operational status of the API and its dependencies.
//	@Tags         system
//	@Produce      json
//	@Success      200  {object}  HealthResponse
//	@Failure      503  {object}  HealthResponse
//	@Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	services := map[string]string{}
	httpStatus := http.StatusOK
	apiStatus := "healthy"

	// Ping the database using the underlying *sql.DB.
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		services["database"] = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
		apiStatus = "degraded"
	} else {
		services["database"] = "healthy"
	}

	c.JSON(httpStatus, HealthResponse{
		Status:    apiStatus,
		Timestamp: time.Now().UTC(),
		Services:  services,
	})
}
