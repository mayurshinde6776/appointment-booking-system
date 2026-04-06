package models

import "time"

// ─── Request DTOs ─────────────────────────────────────────────────────────────

// CreateAppointmentRequest is the payload accepted by POST /appointments.
type CreateAppointmentRequest struct {
	PatientName string    `json:"patient_name" binding:"required,min=2,max=255"`
	DoctorName  string    `json:"doctor_name"  binding:"required,min=2,max=255"`
	Date        time.Time `json:"date"         binding:"required"`
	Duration    int       `json:"duration_mins" binding:"required,min=10,max=480"`
	Notes       string    `json:"notes"`
}

// UpdateAppointmentRequest is the payload accepted by PUT /appointments/:id.
// All fields are optional; only provided fields are updated.
type UpdateAppointmentRequest struct {
	PatientName *string            `json:"patient_name" binding:"omitempty,min=2,max=255"`
	DoctorName  *string            `json:"doctor_name"  binding:"omitempty,min=2,max=255"`
	Date        *time.Time         `json:"date"`
	Duration    *int               `json:"duration_mins" binding:"omitempty,min=10,max=480"`
	Status      *AppointmentStatus `json:"status"       binding:"omitempty,oneof=pending confirmed cancelled completed"`
	Notes       *string            `json:"notes"`
}

// ─── Response DTOs ────────────────────────────────────────────────────────────

// AppointmentResponse is the shape returned to API consumers.
type AppointmentResponse struct {
	ID          uint              `json:"id"`
	PatientName string            `json:"patient_name"`
	DoctorName  string            `json:"doctor_name"`
	Date        time.Time         `json:"date"`
	Duration    int               `json:"duration_mins"`
	Status      AppointmentStatus `json:"status"`
	Notes       string            `json:"notes"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ToResponse converts a domain Appointment to its API response shape.
func ToResponse(a *Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:          a.ID,
		PatientName: a.PatientName,
		DoctorName:  a.DoctorName,
		Date:        a.Date,
		Duration:    a.Duration,
		Status:      a.Status,
		Notes:       a.Notes,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// ToResponseList converts a slice of domain Appointments.
func ToResponseList(appointments []Appointment) []AppointmentResponse {
	result := make([]AppointmentResponse, len(appointments))
	for i := range appointments {
		result[i] = ToResponse(&appointments[i])
	}
	return result
}

// ─── Generic envelope ─────────────────────────────────────────────────────────

// APIResponse is a standard JSON envelope used across all endpoints.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
