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

// CreateAvailabilityRequest is the payload accepted by POST /coaches/availability.
type CreateAvailabilityRequest struct {
	CoachID   uint   `json:"coach_id" binding:"required"`
	DayOfWeek string `json:"day_of_week" binding:"required,oneof=Sunday Monday Tuesday Wednesday Thursday Friday Saturday"`
	StartTime string `json:"start_time" binding:"required"` // Optional validation for HH:MM format could be added
	EndTime   string `json:"end_time" binding:"required"`
}

// CreateBookingRequest is the payload accepted by POST /users/bookings.
type CreateBookingRequest struct {
	UserID   uint      `json:"user_id" binding:"required"`
	CoachID  uint      `json:"coach_id" binding:"required"`
	DateTime time.Time `json:"datetime" binding:"required"`
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

// AvailabilityResponse is the API response shape for Availability.
type AvailabilityResponse struct {
	ID        uint   `json:"id"`
	CoachID   uint   `json:"coach_id"`
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// parseDayOfWeek converts a string day to a DayOfWeek type.
func ParseDayOfWeek(day string) DayOfWeek {
	switch day {
	case "Sunday":
		return Sunday
	case "Monday":
		return Monday
	case "Tuesday":
		return Tuesday
	case "Wednesday":
		return Wednesday
	case "Thursday":
		return Thursday
	case "Friday":
		return Friday
	case "Saturday":
		return Saturday
	default:
		return Sunday // Should be handled by validation beforehand
	}
}

// formatDayOfWeek converts DayOfWeek type to string.
func FormatDayOfWeek(day DayOfWeek) string {
	switch day {
	case Sunday:
		return "Sunday"
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	default:
		return "Unknown"
	}
}

// ToAvailabilityResponse converts a domain Availability to API response shape.
func ToAvailabilityResponse(a *Availability) AvailabilityResponse {
	return AvailabilityResponse{
		ID:        a.ID,
		CoachID:   a.CoachID,
		DayOfWeek: FormatDayOfWeek(a.DayOfWeek),
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
	}
}

// BookingResponse represents the formatted structure returned upon successful booking.
type BookingResponse struct {
	ID        uint          `json:"id"`
	UserID    uint          `json:"user_id"`
	CoachID   uint          `json:"coach_id"`
	SlotTime  time.Time     `json:"slot_time"`
	Status    BookingStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

// ToBookingResponse converts a domain Booking over to the API DTO format.
func ToBookingResponse(b *Booking) BookingResponse {
	return BookingResponse{
		ID:        b.ID,
		UserID:    b.UserID,
		CoachID:   b.CoachID,
		SlotTime:  b.SlotTime.UTC(),
		Status:    b.Status,
		CreatedAt: b.CreatedAt.UTC(),
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
