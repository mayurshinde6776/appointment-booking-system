package services

import (
	"context"
	"fmt"
	"time"

	"appointment-booking/internal/models"
	"appointment-booking/internal/repositories"
)

// BookingService integrates orchestration for verifying valid slots before inserting reservations.
type BookingService interface {
	CreateBooking(ctx context.Context, req models.CreateBookingRequest) (*models.BookingResponse, error)
}

type bookingService struct {
	bookingRepo repositories.BookingRepository
	slotSvc     SlotService
}

// NewBookingService constructs a BookingService.
func NewBookingService(bookingRepo repositories.BookingRepository, slotSvc SlotService) BookingService {
	return &bookingService{bookingRepo: bookingRepo, slotSvc: slotSvc}
}

func (s *bookingService) CreateBooking(
	ctx context.Context, req models.CreateBookingRequest,
) (*models.BookingResponse, error) {

	// 1. Requirements: Check slot availability.
	// We extract the date component and directly query our pre-built availability logic.
	reqDateStr := req.DateTime.Format("2006-01-02")
	availableSlots, err := s.slotSvc.GetAvailableSlots(ctx, req.CoachID, reqDateStr)
	if err != nil {
		return nil, fmt.Errorf("service: failed validating existing availability: %w", err)
	}

	// 2. Validate the requested time precisely matches an openly generated slot.
	requestedISO := req.DateTime.UTC().Format(time.RFC3339)
	isSlotAvailable := false
	for _, slot := range availableSlots {
		if slot == requestedISO {
			isSlotAvailable = true
			break
		}
	}

	if !isSlotAvailable {
		// Stop before hitting the database if the slot simply doesn't exist or isn't on a valid 30min scale boundary
		return nil, fmt.Errorf("slot %s is either outside coach availability bounds or already occupied", requestedISO)
	}

	// 3. Craft Domain Booking
	booking := &models.Booking{
		UserID:   req.UserID,
		CoachID:  req.CoachID,
		SlotTime: req.DateTime,
		Status:   models.BookingStatusConfirmed, // Presume auto-confirmed for simple reservations
	}

	// 4. Create Bookings (with internal transactional insert + duplicate safety handling)
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err // Can expose ErrSlotAlreadyBooked upwards
	}

	resp := models.ToBookingResponse(booking)
	return &resp, nil
}
