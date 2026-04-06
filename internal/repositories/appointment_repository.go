package repositories

import (
	"context"
	"errors"
	"fmt"

	"appointment-booking/internal/models"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a record does not exist.
var ErrNotFound = errors.New("record not found")

// AppointmentRepository defines the data-access contract.
type AppointmentRepository interface {
	Create(ctx context.Context, a *models.Appointment) error
	GetByID(ctx context.Context, id uint) (*models.Appointment, error)
	List(ctx context.Context, filter ListFilter) ([]models.Appointment, int64, error)
	Update(ctx context.Context, a *models.Appointment) error
	Delete(ctx context.Context, id uint) error
}

// ListFilter carries pagination + optional filters for List queries.
type ListFilter struct {
	Page       int
	PageSize   int
	DoctorName string
	Status     string
}

// appointmentRepository is the concrete GORM-backed implementation.
type appointmentRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository constructs an AppointmentRepository backed by GORM.
func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) Create(ctx context.Context, a *models.Appointment) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return fmt.Errorf("repository: create appointment: %w", err)
	}
	return nil
}

func (r *appointmentRepository) GetByID(ctx context.Context, id uint) (*models.Appointment, error) {
	var a models.Appointment
	err := r.db.WithContext(ctx).First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repository: get appointment %d: %w", id, err)
	}
	return &a, nil
}

func (r *appointmentRepository) List(ctx context.Context, f ListFilter) ([]models.Appointment, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Appointment{})

	if f.DoctorName != "" {
		query = query.Where("doctor_name ILIKE ?", "%"+f.DoctorName+"%")
	}
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repository: count appointments: %w", err)
	}

	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}
	offset := (f.Page - 1) * f.PageSize

	var appointments []models.Appointment
	if err := query.
		Order("date ASC").
		Offset(offset).
		Limit(f.PageSize).
		Find(&appointments).Error; err != nil {
		return nil, 0, fmt.Errorf("repository: list appointments: %w", err)
	}

	return appointments, total, nil
}

func (r *appointmentRepository) Update(ctx context.Context, a *models.Appointment) error {
	result := r.db.WithContext(ctx).Save(a)
	if result.Error != nil {
		return fmt.Errorf("repository: update appointment %d: %w", a.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *appointmentRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Appointment{}, id)
	if result.Error != nil {
		return fmt.Errorf("repository: delete appointment %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
