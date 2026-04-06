package services

import (
	"context"
	"fmt"

	"appointment-booking/internal/models"
	"appointment-booking/internal/repositories"
)

// AvailabilityService defines business logic for coach availability.
type AvailabilityService interface {
	SetAvailability(ctx context.Context, req models.CreateAvailabilityRequest) (*models.AvailabilityResponse, error)
}

type availabilityService struct {
	repo repositories.AvailabilityRepository
}

// NewAvailabilityService constructs an AvailabilityService.
func NewAvailabilityService(repo repositories.AvailabilityRepository) AvailabilityService {
	return &availabilityService{repo: repo}
}

func (s *availabilityService) SetAvailability(
	ctx context.Context, req models.CreateAvailabilityRequest,
) (*models.AvailabilityResponse, error) {

	avail := &models.Availability{
		CoachID:   req.CoachID,
		DayOfWeek: models.ParseDayOfWeek(req.DayOfWeek),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := s.repo.Create(ctx, avail); err != nil {
		return nil, fmt.Errorf("service: set availability: %w", err)
	}

	resp := models.ToAvailabilityResponse(avail)
	return &resp, nil
}
