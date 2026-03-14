package services

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/itcodehery/irctc-simulator/database"
	"github.com/itcodehery/irctc-simulator/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BookingJob represents a booking request queued for processing.
type BookingJob struct {
	Request    models.BookingRequest
	UserID     uint
	ResultChan chan BookingResult
}

// BookingResult holds the outcome of processing a booking job.
type BookingResult struct {
	Booking *models.Booking
	Error   error
}

// BookingService manages concurrent ticket bookings using a worker pool.
type BookingService struct {
	jobQueue   chan BookingJob
	workers    int
	wg         sync.WaitGroup
	mu         sync.Mutex // Protects seat allocation
	stopChan   chan struct{}
	isRunning  bool
}

// NewBookingService creates and starts a booking service with the given worker count and queue size.
func NewBookingService(workers, queueSize int) *BookingService {
	bs := &BookingService{
		jobQueue: make(chan BookingJob, queueSize),
		workers:  workers,
		stopChan: make(chan struct{}),
	}
	bs.Start()
	return bs
}

// Start launches the worker pool goroutines.
func (bs *BookingService) Start() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.isRunning {
		return
	}

	slog.Info("Starting booking service worker pool",
		"workers", bs.workers,
		"queue_size", cap(bs.jobQueue),
	)

	for i := 0; i < bs.workers; i++ {
		bs.wg.Add(1)
		go bs.worker(i)
	}

	bs.isRunning = true
}

// Stop gracefully shuts down the worker pool, waiting for in-flight jobs.
func (bs *BookingService) Stop() {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if !bs.isRunning {
		return
	}

	slog.Info("Stopping booking service worker pool...")
	close(bs.stopChan)
	bs.wg.Wait()
	bs.isRunning = false
	slog.Info("Booking service worker pool stopped")
}

// SubmitBooking enqueues a booking job and returns the result.
func (bs *BookingService) SubmitBooking(req models.BookingRequest, userID uint) (*models.Booking, error) {
	resultChan := make(chan BookingResult, 1)

	job := BookingJob{
		Request:    req,
		UserID:     userID,
		ResultChan: resultChan,
	}

	// Try to enqueue with a timeout
	select {
	case bs.jobQueue <- job:
		slog.Info("Booking job enqueued",
			"user_id", userID,
			"train_id", req.TrainID,
			"num_seats", req.NumSeats,
			"queue_length", len(bs.jobQueue),
		)
	default:
		return nil, fmt.Errorf("booking queue is full, please try again later")
	}

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		return result.Booking, result.Error
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("booking request timed out")
	}
}

// worker processes booking jobs from the queue.
func (bs *BookingService) worker(id int) {
	defer bs.wg.Done()

	slog.Info("Booking worker started", "worker_id", id)

	for {
		select {
		case job := <-bs.jobQueue:
			slog.Info("Worker processing booking",
				"worker_id", id,
				"user_id", job.UserID,
				"train_id", job.Request.TrainID,
			)

			booking, err := bs.processBooking(job)
			job.ResultChan <- BookingResult{
				Booking: booking,
				Error:   err,
			}

		case <-bs.stopChan:
			slog.Info("Booking worker shutting down", "worker_id", id)
			return
		}
	}
}

// processBooking handles the actual seat reservation logic with mutex protection.
func (bs *BookingService) processBooking(job BookingJob) (*models.Booking, error) {
	// Use database-level locking for seat allocation to prevent double-booking
	var booking *models.Booking

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Lock the train row for update (pessimistic locking)
		var train models.Train
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&train, job.Request.TrainID).Error; err != nil {
			return fmt.Errorf("train not found: %w", err)
		}

		// Check seat availability
		if train.AvailableSeats < job.Request.NumSeats {
			if train.AvailableSeats == 0 {
				return fmt.Errorf("no seats available on train %s", train.TrainNumber)
			}
			return fmt.Errorf("only %d seats available, requested %d",
				train.AvailableSeats, job.Request.NumSeats)
		}

		// Generate PNR
		pnr := generatePNR()

		// Generate seat numbers
		startSeat := train.TotalSeats - train.AvailableSeats + 1
		seatNumbers := make([]string, job.Request.NumSeats)
		for i := 0; i < job.Request.NumSeats; i++ {
			seatNumbers[i] = fmt.Sprintf("%s-%d", train.TrainClass, startSeat+i)
		}

		// Calculate fare
		totalFare := float64(job.Request.NumSeats) * train.FarePerSeat

		// Create booking
		booking = &models.Booking{
			PNR:           pnr,
			UserID:        job.UserID,
			TrainID:       job.Request.TrainID,
			PassengerName: job.Request.PassengerName,
			NumSeats:      job.Request.NumSeats,
			TotalFare:     totalFare,
			Status:        "CONFIRMED",
			JourneyDate:   job.Request.JourneyDate,
			SeatNumbers:   strings.Join(seatNumbers, ","),
		}

		if err := tx.Create(booking).Error; err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}

		// Update available seats
		train.AvailableSeats -= job.Request.NumSeats
		if err := tx.Save(&train).Error; err != nil {
			return fmt.Errorf("failed to update seat availability: %w", err)
		}

		slog.Info("Booking confirmed",
			"pnr", pnr,
			"train", train.TrainNumber,
			"seats", strings.Join(seatNumbers, ","),
			"fare", totalFare,
		)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load related data
	database.DB.Preload("Train").First(booking, booking.ID)

	return booking, nil
}

// QueueLength returns the current number of jobs waiting in the queue.
func (bs *BookingService) QueueLength() int {
	return len(bs.jobQueue)
}

// generatePNR creates a random 10-character PNR number.
func generatePNR() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	pnr := make([]byte, 10)
	for i := range pnr {
		pnr[i] = charset[rand.Intn(len(charset))]
	}
	return string(pnr)
}

// CancelBooking cancels a booking and releases seats back to the train.
func CancelBooking(pnr string, userID uint) (*models.Booking, error) {
	var booking models.Booking

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("pnr = ? AND user_id = ?", pnr, userID).First(&booking).Error; err != nil {
			return fmt.Errorf("booking not found: %w", err)
		}

		if booking.Status == "CANCELLED" {
			return fmt.Errorf("booking is already cancelled")
		}

		// Update booking status
		booking.Status = "CANCELLED"
		if err := tx.Save(&booking).Error; err != nil {
			return fmt.Errorf("failed to cancel booking: %w", err)
		}

		// Release seats back
		if err := tx.Model(&models.Train{}).
			Where("id = ?", booking.TrainID).
			Update("available_seats", gorm.Expr("available_seats + ?", booking.NumSeats)).
			Error; err != nil {
			return fmt.Errorf("failed to release seats: %w", err)
		}

		slog.Info("Booking cancelled",
			"pnr", pnr,
			"seats_released", booking.NumSeats,
		)

		return nil
	})

	if err != nil {
		return nil, err
	}

	database.DB.Preload("Train").First(&booking, booking.ID)
	return &booking, nil
}
