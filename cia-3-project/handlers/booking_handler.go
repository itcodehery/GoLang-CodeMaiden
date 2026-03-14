package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itcodehery/irctc-simulator/database"
	"github.com/itcodehery/irctc-simulator/models"
	"github.com/itcodehery/irctc-simulator/services"
)

// BookingHandler handles ticket booking requests.
type BookingHandler struct {
	bookingService *services.BookingService
}

// NewBookingHandler creates a new BookingHandler.
func NewBookingHandler(bs *services.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bs}
}

// CreateBooking processes a new ticket booking request.
// POST /api/v1/bookings
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req models.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid booking request",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	userID := c.GetUint("user_id")

	// Submit to worker pool
	booking, err := h.bookingService.SubmitBooking(req, userID)
	if err != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "booking failed",
			Code:    http.StatusConflict,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.BookingResponse{
		Message: "Ticket booked successfully!",
		Booking: *booking,
	})
}

// GetBooking retrieves a booking by PNR.
// GET /api/v1/bookings/:pnr
func (h *BookingHandler) GetBooking(c *gin.Context) {
	pnr := c.Param("pnr")
	userID := c.GetUint("user_id")

	var booking models.Booking
	if err := database.DB.Preload("Train").
		Where("pnr = ? AND user_id = ?", pnr, userID).
		First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "booking not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// ListBookings returns all bookings for the authenticated user.
// GET /api/v1/bookings
func (h *BookingHandler) ListBookings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var bookings []models.Booking
	if err := database.DB.Preload("Train").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to fetch bookings",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":    len(bookings),
		"bookings": bookings,
	})
}

// CancelBooking cancels an existing booking.
// DELETE /api/v1/bookings/:pnr
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	pnr := c.Param("pnr")
	userID := c.GetUint("user_id")

	booking, err := services.CancelBooking(pnr, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "cancellation failed",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.BookingResponse{
		Message: "Booking cancelled successfully. Refund will be processed.",
		Booking: *booking,
	})
}

// GetQueueStatus returns the current booking queue status.
// GET /api/v1/bookings/queue/status
func (h *BookingHandler) GetQueueStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"queue_length": h.bookingService.QueueLength(),
		"status":       "operational",
	})
}
