package services_test

import (
	"testing"
	"time"

	"appointment-booking/internal/services"
)

func TestGenerateSlots(t *testing.T) {
	// Let's create an arbitrary baseline date: 2026-05-01 00:00:00 UTC
	baseDate := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		wantCount int
	}{
		{
			name:      "10 to 15 (5 hours, 10 slots)",
			startTime: baseDate.Add(10 * time.Hour), // 10:00
			endTime:   baseDate.Add(15 * time.Hour), // 15:00
			wantCount: 10,                           // 10:00, 10:30, 11:00, 11:30, 12:00, 12:30, 13:00, 13:30, 14:00, 14:30
		},
		{
			name:      "Same start and end (0 slots)",
			startTime: baseDate.Add(10 * time.Hour), // 10:00
			endTime:   baseDate.Add(10 * time.Hour), // 10:00
			wantCount: 0,
		},
		{
			name:      "Incomplete slot end time (1 slot)",
			startTime: baseDate.Add(10 * time.Hour),                   // 10:00
			endTime:   baseDate.Add(10 * time.Hour + 45*time.Minute),  // 10:45
			wantCount: 1,                                              // 10:00
		},
		{
			name:      "Single perfect slot (1 slot)",
			startTime: baseDate.Add(10 * time.Hour),                   // 10:00
			endTime:   baseDate.Add(10 * time.Hour + 30*time.Minute),  // 10:30
			wantCount: 1,                                              // 10:00
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			slots := services.GenerateSlots(tc.startTime, tc.endTime)
			if len(slots) != tc.wantCount {
				t.Errorf("GenerateSlots() returned %d slots, want %d", len(slots), tc.wantCount)
			}

			// Validate slot continuity and duration lengths if slots exist
			if len(slots) > 0 {
				if !slots[0].Equal(tc.startTime) {
					t.Errorf("GenerateSlots() first slot %v, want %v", slots[0], tc.startTime)
				}
				for i := 1; i < len(slots); i++ {
					diff := slots[i].Sub(slots[i-1])
					if diff != 30*time.Minute {
						t.Errorf("GenerateSlots() gap between slot %d and %d is %v, want 30m", i-1, i, diff)
					}
				}
				// Guarantee the last slot fits within end time
				lastSlotEnd := slots[len(slots)-1].Add(30 * time.Minute)
				if lastSlotEnd.After(tc.endTime) {
					t.Errorf("GenerateSlots() last slot ends at %v, which is AFTER the allowed end time %v", lastSlotEnd, tc.endTime)
				}
			}
		})
	}
}
