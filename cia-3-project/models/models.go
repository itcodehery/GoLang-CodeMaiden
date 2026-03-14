package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered passenger/user.
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	FullName  string         `gorm:"size:100;not null" json:"full_name"`
	Phone     string         `gorm:"size:15" json:"phone"`
	Bookings  []Booking      `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
}

// Train represents a train with its route and schedule.
type Train struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	TrainNumber   string         `gorm:"uniqueIndex;size:10;not null" json:"train_number"`
	TrainName     string         `gorm:"size:100;not null" json:"train_name"`
	Source        string         `gorm:"size:50;not null;index" json:"source"`
	Destination   string         `gorm:"size:50;not null;index" json:"destination"`
	DepartureTime string         `gorm:"size:10;not null" json:"departure_time"`
	ArrivalTime   string         `gorm:"size:10;not null" json:"arrival_time"`
	TotalSeats    int            `gorm:"not null" json:"total_seats"`
	AvailableSeats int           `gorm:"not null" json:"available_seats"`
	FarePerSeat   float64        `gorm:"not null" json:"fare_per_seat"`
	TrainClass    string         `gorm:"size:10;not null;default:'SL'" json:"train_class"` // SL, 3A, 2A, 1A
	DaysOfWeek    string         `gorm:"size:20;not null;default:'ALL'" json:"days_of_week"`
	Bookings      []Booking      `gorm:"foreignKey:TrainID" json:"bookings,omitempty"`
}

// Booking represents a ticket booking made by a user.
type Booking struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	PNR           string         `gorm:"uniqueIndex;size:10;not null" json:"pnr"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	TrainID       uint           `gorm:"not null;index" json:"train_id"`
	PassengerName string         `gorm:"size:100;not null" json:"passenger_name"`
	NumSeats      int            `gorm:"not null;default:1" json:"num_seats"`
	TotalFare     float64        `gorm:"not null" json:"total_fare"`
	Status        string         `gorm:"size:20;not null;default:'PENDING'" json:"status"` // PENDING, CONFIRMED, CANCELLED, WAITLISTED
	JourneyDate   string         `gorm:"size:10;not null" json:"journey_date"`
	SeatNumbers   string         `gorm:"size:100" json:"seat_numbers"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Train         Train          `gorm:"foreignKey:TrainID" json:"train,omitempty"`
}

// --- Request/Response DTOs ---

// RegisterRequest is the signup payload.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone"`
}

// LoginRequest is the login payload.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

// BookingRequest is the ticket booking payload.
type BookingRequest struct {
	TrainID       uint   `json:"train_id" binding:"required"`
	PassengerName string `json:"passenger_name" binding:"required"`
	NumSeats      int    `json:"num_seats" binding:"required,min=1,max=6"`
	JourneyDate   string `json:"journey_date" binding:"required"`
}

// BookingResponse wraps a booking with a status message.
type BookingResponse struct {
	Message string  `json:"message"`
	Booking Booking `json:"booking"`
}

// TrainSearchRequest is the search query.
type TrainSearchRequest struct {
	Source      string `form:"source" binding:"required"`
	Destination string `form:"destination" binding:"required"`
	Date        string `form:"date"`
}

// ErrorResponse is a generic error payload.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// HealthResponse is the health check payload.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Version   string    `json:"version"`
}
