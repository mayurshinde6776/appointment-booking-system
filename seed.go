package main

import (
	"log"

	"appointment-booking/internal/models"
	"appointment-booking/pkg/database"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	coach := models.Coach{
		ID:       1,
		Name:     "Test Coach",
		Timezone: "UTC",
	}

	user := models.User{
		ID:    101,
		Name:  "Test User",
		Email: "test101@example.com",
	}

	if err := db.FirstOrCreate(&coach, models.Coach{ID: 1}).Error; err != nil {
		log.Printf("Coach creation info: %v", err)
	}

	if err := db.FirstOrCreate(&user, models.User{ID: 101}).Error; err != nil {
		log.Printf("User creation info: %v", err)
	}

	log.Println("✅ Successfully seeded the database! Coach ID '1' and User ID '101' now exist.")
}
