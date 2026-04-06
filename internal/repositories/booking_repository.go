package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"appointment-booking/internal/models"

	"gorm.io/gorm"
)

// ErrSlotAlreadyBooked indicates a unique constraint violation when double booking.
var ErrSlotAlreadyBooked = fmt.Errorf("slot already booked")

type BookingRepository interface {
	GetBookingsByCoachAndDateRange(ctx context.Context, coachID uint, start, end time.Time) ([]models.Booking, error)
	Create(ctx context.Context, b *models.Booking) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) GetBookingsByCoachAndDateRange(ctx context.Context, coachID uint, start, end time.Time) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.WithContext(ctx).
		Where("coach_id = ? AND slot_time >= ? AND slot_time < ?", coachID, start, end).
		Where("status != ?", models.BookingStatusCancelled).
		Find(&bookings).Error
	if err != nil {
		return nil, fmt.Errorf("repository: get bookings by date range: %w", err)
	}
	return bookings, nil
}

func (r *bookingRepository) Create(ctx context.Context, b *models.Booking) error {
	// Requirements: "Use database transaction". GORM automatically provides transacted 
	// writes, but doing it explicitly guarantees the lifecycle requested.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(b).Error; err != nil {
			// Requirements: "Handle duplicate booking using database unique constraint"
			// PostgreSQL state 23505 is unique_violation
			if strings.Contains(err.Error(), "SQLSTATE 23505") || strings.Contains(err.Error(), "duplicate key value") {
				return ErrSlotAlreadyBooked
			}
			return fmt.Errorf("repository: failed to create booking: %w", err)
		}
		return nil
	})
}
