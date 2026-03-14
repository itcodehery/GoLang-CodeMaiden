package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itcodehery/irctc-simulator/config"
	"github.com/itcodehery/irctc-simulator/database"
	"github.com/itcodehery/irctc-simulator/middleware"
	"github.com/itcodehery/irctc-simulator/router"
	"github.com/itcodehery/irctc-simulator/services"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("=== IRCTC Ticket Booking Simulator ===")
	slog.Info("Starting application...")

	// Load configuration
	cfg := config.Load()
	slog.Info("Configuration loaded",
		"port", cfg.ServerPort,
		"db_path", cfg.DBPath,
		"booking_workers", cfg.BookingWorkers,
		"rate_limit", fmt.Sprintf("%d req/%s", cfg.RateLimitRequests, cfg.RateLimitWindow),
	)

	// Set JWT secret
	middleware.JWTSecret = []byte(cfg.JWTSecret)

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Seed sample data
	if err := database.SeedTrains(); err != nil {
		slog.Error("Failed to seed trains", "error", err)
		os.Exit(1)
	}

	// Start booking service worker pool
	bookingService := services.NewBookingService(cfg.BookingWorkers, cfg.BookingQueueSize)
	defer bookingService.Stop()

	// Setup router
	r := router.Setup(cfg, bookingService)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting", "address", srv.Addr)
		fmt.Printf("\n")
		fmt.Printf("  ╔══════════════════════════════════════════════════╗\n")
		fmt.Printf("  ║       IRCTC Ticket Booking Simulator v1.0       ║\n")
		fmt.Printf("  ╠══════════════════════════════════════════════════╣\n")
		fmt.Printf("  ║  Server:   http://localhost:%s                 ║\n", cfg.ServerPort)
		fmt.Printf("  ║  Health:   http://localhost:%s/health           ║\n", cfg.ServerPort)
		fmt.Printf("  ║  API:      http://localhost:%s/api/v1           ║\n", cfg.ServerPort)
		fmt.Printf("  ║  Workers:  %d booking workers                    ║\n", cfg.BookingWorkers)
		fmt.Printf("  ╚══════════════════════════════════════════════════╝\n")
		fmt.Printf("\n")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	// Wait for interrupt signal (SIGINT or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("Shutdown signal received", "signal", sig.String())
	slog.Info("Initiating graceful shutdown...",
		"timeout", cfg.ShutdownTimeout,
	)

	// Create a deadline context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Step 1: Stop accepting new connections
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// Step 2: Stop booking worker pool (waits for in-flight bookings)
	slog.Info("Waiting for in-flight bookings to complete...")
	bookingService.Stop()

	// Step 3: Close database connections
	slog.Info("Closing database connections...")
	database.Close()

	slog.Info("Server gracefully stopped")
}
