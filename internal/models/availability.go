package models

// DayOfWeek mirrors Go's time.Weekday (0 = Sunday … 6 = Saturday).
type DayOfWeek int

const (
	Sunday DayOfWeek = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// Availability defines a recurring weekly time window during which a Coach
// accepts bookings.  StartTime and EndTime are stored as "HH:MM" strings in
// the coach's own timezone (see Coach.Timezone).
type Availability struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"                                         json:"id"`
	CoachID   uint      `gorm:"not null;index"                                                   json:"coach_id"     validate:"required"`
	DayOfWeek DayOfWeek `gorm:"not null;check:day_of_week >= 0 AND day_of_week <= 6"             json:"day_of_week"  validate:"min=0,max=6"`
	StartTime string    `gorm:"type:varchar(5);not null"                                         json:"start_time"   validate:"required"`
	EndTime   string    `gorm:"type:varchar(5);not null"                                         json:"end_time"     validate:"required"`

	// Relationships
	Coach Coach `gorm:"foreignKey:CoachID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"coach,omitempty"`
}
