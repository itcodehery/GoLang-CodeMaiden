package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itcodehery/irctc-simulator/config"
	"github.com/itcodehery/irctc-simulator/handlers"
	"github.com/itcodehery/irctc-simulator/middleware"
	"github.com/itcodehery/irctc-simulator/models"
	"github.com/itcodehery/irctc-simulator/services"
)

// Setup creates and configures the Gin router with all routes and middleware.
func Setup(cfg *config.Config, bookingService *services.BookingService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.StructuredLogger())

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRequests, cfg.RateLimitWindow)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(cfg)
	trainHandler := handlers.NewTrainHandler()
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Application start time for health check
	startTime := time.Now()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, models.HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Uptime:    time.Since(startTime).String(),
			Version:   "1.0.0",
		})
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// API Index
		api.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to IRCTC Ticket Booking Simulator API v1",
				"version": "1.0.0",
				"status":  "active",
			})
		})

		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Public train routes (no auth required for browsing)
		trains := api.Group("/trains")
		{
			trains.GET("", trainHandler.ListTrains)
			trains.GET("/search", trainHandler.SearchTrains)
			trains.GET("/:id", trainHandler.GetTrain)
			trains.GET("/:id/availability", trainHandler.GetAvailability)
		}

		// Protected routes (require JWT)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/user/profile", userHandler.GetProfile)

			// Bookings
			bookings := protected.Group("/bookings")
			{
				bookings.POST("", bookingHandler.CreateBooking)
				bookings.GET("", bookingHandler.ListBookings)
				bookings.GET("/queue/status", bookingHandler.GetQueueStatus)
				bookings.GET("/:pnr", bookingHandler.GetBooking)
				bookings.DELETE("/:pnr", bookingHandler.CancelBooking)
			}
		}
	}

	return r
}
