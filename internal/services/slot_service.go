package services

import (
	"context"
	"fmt"
	"time"

	"appointment-booking/internal/models"
	"appointment-booking/internal/repositories"
)

type SlotService interface {
	GetAvailableSlots(ctx context.Context, coachID uint, dateStr string) ([]string, error)
}

type slotService struct {
	availRepo   repositories.AvailabilityRepository
	bookingRepo repositories.BookingRepository
}

func NewSlotService(availRepo repositories.AvailabilityRepository, bookingRepo repositories.BookingRepository) SlotService {
	return &slotService{availRepo: availRepo, bookingRepo: bookingRepo}
}

func (s *slotService) GetAvailableSlots(ctx context.Context, coachID uint, dateStr string) ([]string, error) {
	// Parse date string (assumed format: "2006-01-02")
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// 1. Determine day_of_week from date
	dayOfWeek := models.DayOfWeek(parsedDate.Weekday())

	// 2. Fetch coach availability from database
	availabilities, err := s.availRepo.GetByCoachAndDay(ctx, coachID, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch availability: %w", err)
	}

	if len(availabilities) == 0 {
		return []string{}, nil
	}

	// Load timezone from the coach (we have it preloaded in availability)
	coachTZ := availabilities[0].Coach.Timezone
	loc, err := time.LoadLocation(coachTZ)
	if err != nil {
		loc = time.UTC
	}

	// Calculate bounds for the requested date inside the coach's timezone
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 4. Fetch existing active bookings for that coach and date
	bookings, err := s.bookingRepo.GetBookingsByCoachAndDateRange(ctx, coachID, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings: %w", err)
	}

	// Store booked slots in a map for constant time lookout
	bookedSlots := make(map[time.Time]bool)
	for _, b := range bookings {
		bookedSlots[b.SlotTime.In(loc)] = true
	}

	var availableSlotsISO []string

	for _, avail := range availabilities {
		// Build combined DateTime strings
		startAvailStr := fmt.Sprintf("%sT%s:00", dateStr, avail.StartTime)
		endAvailStr := fmt.Sprintf("%sT%s:00", dateStr, avail.EndTime)

		startTime, err := time.ParseInLocation("2006-01-02T15:04:05", startAvailStr, loc)
		if err != nil {
			continue
		}
		endTime, err := time.ParseInLocation("2006-01-02T15:04:05", endAvailStr, loc)
		if err != nil {
			continue
		}

		// 3. Generate 30-minute slots
		slots := GenerateSlots(startTime, endTime)

		// 5. Remove booked slots
		for _, slot := range slots {
			if !bookedSlots[slot] {
				// 6. Return available slots as ISO datetime array (Z)
				availableSlotsISO = append(availableSlotsISO, slot.UTC().Format(time.RFC3339))
			}
		}
	}

	// Ensure empty array is safely returned over nil
	if availableSlotsISO == nil {
		return []string{}, nil
	}

	return availableSlotsISO, nil
}

// GenerateSlots creates sequential 30-minute time slots between a start and end time.
// It ensures that every generated slot fits entirely within the provided time boundary.
// For example, if start is 10:00 and end is 15:00, it returns [10:00, 10:30, ... 14:30].
func GenerateSlots(startTime, endTime time.Time) []time.Time {
	// 30-minute slot duration
	slotDuration := 30 * time.Minute
	var slots []time.Time

	current := startTime
	// Ensure the full 30-minute slot fits within or exactly ending at the endTime limit.
	for current.Add(slotDuration).Before(endTime) || current.Add(slotDuration).Equal(endTime) {
		slots = append(slots, current)
		current = current.Add(slotDuration)
	}

	return slots
}
