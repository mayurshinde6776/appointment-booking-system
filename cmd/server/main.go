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
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// ── Auto-migrate models ───────────────────────────────────────────────────
	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	// ── Wire dependencies ─────────────────────────────────────────────────────
	appointmentRepo := repositories.NewAppointmentRepository(db)
	appointmentSvc := services.NewAppointmentService(appointmentRepo)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentSvc)

	healthHandler := handlers.NewHealthHandler(db)

	// ── Router ────────────────────────────────────────────────────────────────
	router := handlers.SetupRouter(healthHandler, appointmentHandler)

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
