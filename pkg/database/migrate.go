package database

import (
	"fmt"
	"log"

	"appointment-booking/internal/models"

	"gorm.io/gorm"
)

// Migrate runs GORM AutoMigrate for every domain model when the server starts.
//
// Migration order is intentional — parent tables must exist before child
// tables that carry foreign keys:
//
//	1. users         (no FK deps)
//	2. coaches       (no FK deps)
//	3. availabilities  (FK → coaches)
//	4. bookings        (FK → users, coaches)
//	5. appointments    (legacy, no FK deps)
//
// AutoMigrate creates missing tables, adds missing columns, and creates
// missing indexes — it never drops columns or tables.
func Migrate(db *gorm.DB) error {
	log.Println("migration: starting auto-migration…")

	migrations := []struct {
		name  string
		model interface{}
	}{
		{"users", &models.User{}},
		{"coaches", &models.Coach{}},
		{"availabilities", &models.Availability{}},
		{"bookings", &models.Booking{}},
		{"appointments", &models.Appointment{}},
	}

	for _, m := range migrations {
		log.Printf("migration: migrating table %q…", m.name)
		if err := db.AutoMigrate(m.model); err != nil {
			return fmt.Errorf("migration: failed to migrate %q: %w", m.name, err)
		}
		log.Printf("migration: table %q OK", m.name)
	}

	log.Println("migration: all tables up to date ✓")
	return nil
}
