package models

import "time"

// Coach represents a doctor or service provider who offers appointments.
type Coach struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"      validate:"required,min=2,max=255"`
	Timezone  string    `gorm:"type:varchar(100);not null;default:'UTC'" json:"timezone" validate:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime"           json:"created_at"`

	// Relationships
	Availabilities []Availability `gorm:"foreignKey:CoachID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"availabilities,omitempty"`
	Bookings       []Booking      `gorm:"foreignKey:CoachID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bookings,omitempty"`
}
