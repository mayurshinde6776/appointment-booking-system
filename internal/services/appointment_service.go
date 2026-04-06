package services

import (
	"context"
	"errors"
	"fmt"

	"appointment-booking/internal/models"
	"appointment-booking/internal/repositories"
)

// ErrNotFound is re-exported so handlers only need one import.
var ErrNotFound = repositories.ErrNotFound

// AppointmentService defines business-logic operations on appointments.
type AppointmentService interface {
	CreateAppointment(ctx context.Context, req models.CreateAppointmentRequest) (*models.AppointmentResponse, error)
	GetAppointment(ctx context.Context, id uint) (*models.AppointmentResponse, error)
	ListAppointments(ctx context.Context, filter repositories.ListFilter) ([]models.AppointmentResponse, int64, error)
	UpdateAppointment(ctx context.Context, id uint, req models.UpdateAppointmentRequest) (*models.AppointmentResponse, error)
	DeleteAppointment(ctx context.Context, id uint) error
}

type appointmentService struct {
	repo repositories.AppointmentRepository
}

// NewAppointmentService constructs an AppointmentService.
func NewAppointmentService(repo repositories.AppointmentRepository) AppointmentService {
	return &appointmentService{repo: repo}
}

func (s *appointmentService) CreateAppointment(
	ctx context.Context, req models.CreateAppointmentRequest,
) (*models.AppointmentResponse, error) {

	appointment := &models.Appointment{
		PatientName: req.PatientName,
		DoctorName:  req.DoctorName,
		Date:        req.Date,
		Duration:    req.Duration,
		Status:      models.StatusPending,
		Notes:       req.Notes,
	}

	if err := s.repo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("service: create appointment: %w", err)
	}

	resp := models.ToResponse(appointment)
	return &resp, nil
}

func (s *appointmentService) GetAppointment(
	ctx context.Context, id uint,
) (*models.AppointmentResponse, error) {

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("service: get appointment: %w", err)
	}

	resp := models.ToResponse(a)
	return &resp, nil
}

func (s *appointmentService) ListAppointments(
	ctx context.Context, filter repositories.ListFilter,
) ([]models.AppointmentResponse, int64, error) {

	appointments, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("service: list appointments: %w", err)
	}

	return models.ToResponseList(appointments), total, nil
}

func (s *appointmentService) UpdateAppointment(
	ctx context.Context, id uint, req models.UpdateAppointmentRequest,
) (*models.AppointmentResponse, error) {

	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("service: fetch for update: %w", err)
	}

	// Apply only the fields that were explicitly set in the request.
	if req.PatientName != nil {
		a.PatientName = *req.PatientName
	}
	if req.DoctorName != nil {
		a.DoctorName = *req.DoctorName
	}
	if req.Date != nil {
		a.Date = *req.Date
	}
	if req.Duration != nil {
		a.Duration = *req.Duration
	}
	if req.Status != nil {
		a.Status = *req.Status
	}
	if req.Notes != nil {
		a.Notes = *req.Notes
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, fmt.Errorf("service: update appointment: %w", err)
	}

	resp := models.ToResponse(a)
	return &resp, nil
}

func (s *appointmentService) DeleteAppointment(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("service: delete appointment: %w", err)
	}
	return nil
}
