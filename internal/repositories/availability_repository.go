package repositories

import (
	"context"
	"fmt"

	"appointment-booking/internal/models"

	"gorm.io/gorm"
)

// AvailabilityRepository defines the data-access contract.
type AvailabilityRepository interface {
	Create(ctx context.Context, a *models.Availability) error
	GetByCoachAndDay(ctx context.Context, coachID uint, day models.DayOfWeek) ([]models.Availability, error)
}

type availabilityRepository struct {
	db *gorm.DB
}

// NewAvailabilityRepository constructs an AvailabilityRepository backed by GORM.
func NewAvailabilityRepository(db *gorm.DB) AvailabilityRepository {
	return &availabilityRepository{db: db}
}

func (r *availabilityRepository) Create(ctx context.Context, a *models.Availability) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return fmt.Errorf("repository: create availability: %w", err)
	}
	return nil
}

func (r *availabilityRepository) GetByCoachAndDay(ctx context.Context, coachID uint, day models.DayOfWeek) ([]models.Availability, error) {
	var availabilities []models.Availability
	err := r.db.WithContext(ctx).Preload("Coach").Where("coach_id = ? AND day_of_week = ?", coachID, day).Find(&availabilities).Error
	if err != nil {
		return nil, fmt.Errorf("repository: get availability: %w", err)
	}
	return availabilities, nil
}
