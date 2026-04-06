package handlers

import (
	"net/http"
	"strconv"

	"appointment-booking/internal/models"
	"appointment-booking/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	slotSvc    services.SlotService
	bookingSvc services.BookingService
}

func NewUserHandler(slotSvc services.SlotService, bookingSvc services.BookingService) *UserHandler {
	return &UserHandler{slotSvc: slotSvc, bookingSvc: bookingSvc}
}

// GetAvailableSlots handles GET /api/v1/users/slots
// Requires: coach_id and date (YYYY-MM-DD)
func (h *UserHandler) GetAvailableSlots(c *gin.Context) {
	coachIDStr := c.Query("coach_id")
	dateStr := c.Query("date")

	coachID, err := strconv.ParseUint(coachIDStr, 10, 32)
	if err != nil || coachID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid or missing coach_id",
		})
		return
	}

	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date is required (format: YYYY-MM-DD)",
		})
		return
	}

	slots, err := h.slotSvc.GetAvailableSlots(c.Request.Context(), uint(coachID), dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// The problem instructions specify returning the array directly
	c.JSON(http.StatusOK, slots)
}

// CreateBooking handles POST /api/v1/users/bookings
func (h *UserHandler) CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.bookingSvc.CreateBooking(c.Request.Context(), req)
	if err != nil {
		// Specific mapping for duplicate key HTTP 409 Conflict
		if err.Error() == "slot already booked" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
