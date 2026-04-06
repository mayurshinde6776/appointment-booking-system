package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"appointment-booking/internal/handlers"
	"appointment-booking/internal/repositories"
	"appointment-booking/internal/services"
	"appointment-booking/pkg/database"
)

func main() {
	// ── Database ──────────────────────────────────────────────────────────────
	// Reads DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME from environment.
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("startup: database connection failed: %v", err)
	}
	log.Println("startup: database connected ✓")

	// ── Auto-migrate models ───────────────────────────────────────────────────
	// Runs on every startup — creates/updates tables for:
	// User, Coach, Availability, Booking, Appointment
	if err := database.Migrate(db); err != nil {
		log.Fatalf("startup: migration failed: %v", err)
	}
	log.Println("startup: migrations applied ✓")

	// ── Wire dependencies ─────────────────────────────────────────────────────
	appointmentRepo := repositories.NewAppointmentRepository(db)
	appointmentSvc := services.NewAppointmentService(appointmentRepo)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentSvc)

	availabilityRepo := repositories.NewAvailabilityRepository(db)
	availabilitySvc := services.NewAvailabilityService(availabilityRepo)
	availabilityHandler := handlers.NewAvailabilityHandler(availabilitySvc)

	bookingRepo := repositories.NewBookingRepository(db)
	slotSvc := services.NewSlotService(availabilityRepo, bookingRepo)
	
	bookingSvc := services.NewBookingService(bookingRepo, slotSvc)
	userHandler := handlers.NewUserHandler(slotSvc, bookingSvc)

	healthHandler := handlers.NewHealthHandler(db)

	// ── Router ────────────────────────────────────────────────────────────────
	router := handlers.SetupRouter(healthHandler, appointmentHandler, availabilityHandler, userHandler)

	// ── HTTP server ───────────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start in a goroutine so we can listen for shutdown signals.
	go func() {
		log.Printf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
