package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"appointment-booking/internal/models"
	"appointment-booking/internal/repositories"
	"appointment-booking/internal/services"

	"github.com/gin-gonic/gin"
)

// AppointmentHandler wires HTTP layer to the appointment service.
type AppointmentHandler struct {
	svc services.AppointmentService
}

// NewAppointmentHandler constructs an AppointmentHandler.
func NewAppointmentHandler(svc services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{svc: svc}
}

// ─── Create ───────────────────────────────────────────────────────────────────

// Create handles POST /api/v1/appointments.
func (h *AppointmentHandler) Create(c *gin.Context) {
	var req models.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svc.CreateAppointment(c.Request.Context(), req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to create appointment")
		return
	}

	respondSuccess(c, http.StatusCreated, "appointment created", resp)
}

// ─── Get by ID ───────────────────────────────────────────────────────────────

// GetByID handles GET /api/v1/appointments/:id.
func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid appointment id")
		return
	}

	resp, err := h.svc.GetAppointment(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondError(c, http.StatusNotFound, "appointment not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to fetch appointment")
		return
	}

	respondSuccess(c, http.StatusOK, "", resp)
}

// ─── List ────────────────────────────────────────────────────────────────────

// List handles GET /api/v1/appointments.
func (h *AppointmentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := repositories.ListFilter{
		Page:       page,
		PageSize:   pageSize,
		DoctorName: c.Query("doctor_name"),
		Status:     c.Query("status"),
	}

	appointments, total, err := h.svc.ListAppointments(c.Request.Context(), filter)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to fetch appointments")
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: gin.H{
			"appointments": appointments,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
		},
	})
}

// ─── Update ───────────────────────────────────────────────────────────────────

// Update handles PUT /api/v1/appointments/:id.
func (h *AppointmentHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid appointment id")
		return
	}

	var req models.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svc.UpdateAppointment(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondError(c, http.StatusNotFound, "appointment not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to update appointment")
		return
	}

	respondSuccess(c, http.StatusOK, "appointment updated", resp)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

// Delete handles DELETE /api/v1/appointments/:id.
func (h *AppointmentHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid appointment id")
		return
	}

	if err := h.svc.DeleteAppointment(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondError(c, http.StatusNotFound, "appointment not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to delete appointment")
		return
	}

	respondSuccess(c, http.StatusOK, "appointment deleted", nil)
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 64)
	return uint(val), err
}

func respondSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, models.APIResponse{Success: true, Message: message, Data: data})
}

func respondError(c *gin.Context, code int, errMsg string) {
	c.AbortWithStatusJSON(code, models.APIResponse{Success: false, Error: errMsg})
}
