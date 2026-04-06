package models

import "time"

// BookingStatus represents the lifecycle state of a booking.
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

// Booking represents a scheduled appointment between a User and a Coach.
//
// Unique constraint on (coach_id, slot_time) prevents double-booking:
// a coach can only have one booking per time slot.
type Booking struct {
	ID        uint          `gorm:"primaryKey;autoIncrement"                                                   json:"id"`
	UserID    uint          `gorm:"not null;index"                                                             json:"user_id"   validate:"required"`
	CoachID   uint          `gorm:"not null;uniqueIndex:idx_coach_slot_time"                                   json:"coach_id"  validate:"required"`
	SlotTime  time.Time     `gorm:"not null;uniqueIndex:idx_coach_slot_time"                                   json:"slot_time" validate:"required"`
	Status    BookingStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','confirmed','cancelled','completed')" json:"status"`
	CreatedAt time.Time     `gorm:"autoCreateTime"                                                             json:"created_at"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"  json:"user,omitempty"`
	Coach Coach `gorm:"foreignKey:CoachID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"coach,omitempty"`
}
