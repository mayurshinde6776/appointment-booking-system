package handlers

import (
	"net/http"

	"appointment-booking/internal/models"
	"appointment-booking/internal/services"

	"github.com/gin-gonic/gin"
)

// AvailabilityHandler wires HTTP layer to the availability service.
type AvailabilityHandler struct {
	svc services.AvailabilityService
}

// NewAvailabilityHandler constructs an AvailabilityHandler.
func NewAvailabilityHandler(svc services.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{svc: svc}
}

// SetAvailability handles POST /api/v1/coaches/availability.
func (h *AvailabilityHandler) SetAvailability(c *gin.Context) {
	var req models.CreateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svc.SetAvailability(c.Request.Context(), req)
	if err != nil {
		// Return exactly what the DB error is, likely a missing foreign key (coach needs to exist!)
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, "availability set successfully", resp)
}
