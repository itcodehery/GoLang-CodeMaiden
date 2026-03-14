package database

import (
	"fmt"
	"log/slog"

	"github.com/itcodehery/irctc-simulator/config"
	"github.com/itcodehery/irctc-simulator/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance.
var DB *gorm.DB

// Initialize sets up the database connection with connection pooling and auto-migration.
func Initialize(cfg *config.Config) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pooling via underlying sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxDBConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxDBConnections / 2)

	// Enable WAL mode for better concurrent read/write performance in SQLite
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA busy_timeout=5000")

	slog.Info("Database connection established",
		"path", cfg.DBPath,
		"max_connections", cfg.MaxDBConnections,
	)

	// Auto-migrate models
	if err := DB.AutoMigrate(&models.User{}, &models.Train{}, &models.Booking{}); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	slog.Info("Database migration completed successfully")
	return nil
}

// SeedTrains populates the database with sample trains if the table is empty.
func SeedTrains() error {
	var count int64
	DB.Model(&models.Train{}).Count(&count)
	if count > 0 {
		slog.Info("Trains already seeded", "count", count)
		return nil
	}

	trains := []models.Train{
		{
			TrainNumber: "12301", TrainName: "Howrah Rajdhani Express",
			Source: "NEW DELHI", Destination: "HOWRAH",
			DepartureTime: "16:55", ArrivalTime: "10:00",
			TotalSeats: 500, AvailableSeats: 500,
			FarePerSeat: 1545.00, TrainClass: "3A", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12302", TrainName: "New Delhi Rajdhani Express",
			Source: "HOWRAH", Destination: "NEW DELHI",
			DepartureTime: "14:05", ArrivalTime: "10:25",
			TotalSeats: 500, AvailableSeats: 500,
			FarePerSeat: 1545.00, TrainClass: "3A", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12951", TrainName: "Mumbai Rajdhani Express",
			Source: "NEW DELHI", Destination: "MUMBAI CENTRAL",
			DepartureTime: "16:35", ArrivalTime: "08:35",
			TotalSeats: 600, AvailableSeats: 600,
			FarePerSeat: 1875.00, TrainClass: "2A", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12952", TrainName: "New Delhi Rajdhani Express",
			Source: "MUMBAI CENTRAL", Destination: "NEW DELHI",
			DepartureTime: "17:00", ArrivalTime: "08:35",
			TotalSeats: 600, AvailableSeats: 600,
			FarePerSeat: 1875.00, TrainClass: "2A", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12259", TrainName: "Sealdah Duronto Express",
			Source: "NEW DELHI", Destination: "SEALDAH",
			DepartureTime: "20:15", ArrivalTime: "10:45",
			TotalSeats: 400, AvailableSeats: 400,
			FarePerSeat: 1320.00, TrainClass: "3A", DaysOfWeek: "MON,WED,FRI",
		},
		{
			TrainNumber: "12627", TrainName: "Karnataka Express",
			Source: "NEW DELHI", Destination: "BANGALORE",
			DepartureTime: "21:20", ArrivalTime: "06:40",
			TotalSeats: 800, AvailableSeats: 800,
			FarePerSeat: 1055.00, TrainClass: "SL", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12622", TrainName: "Tamil Nadu Express",
			Source: "NEW DELHI", Destination: "CHENNAI CENTRAL",
			DepartureTime: "22:30", ArrivalTime: "07:10",
			TotalSeats: 750, AvailableSeats: 750,
			FarePerSeat: 985.00, TrainClass: "SL", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12839", TrainName: "Chennai Mail",
			Source: "HOWRAH", Destination: "CHENNAI CENTRAL",
			DepartureTime: "23:50", ArrivalTime: "04:45",
			TotalSeats: 650, AvailableSeats: 650,
			FarePerSeat: 890.00, TrainClass: "SL", DaysOfWeek: "ALL",
		},
		{
			TrainNumber: "12431", TrainName: "Trivandrum Rajdhani Express",
			Source: "NEW DELHI", Destination: "TRIVANDRUM",
			DepartureTime: "10:55", ArrivalTime: "05:15",
			TotalSeats: 350, AvailableSeats: 350,
			FarePerSeat: 2680.00, TrainClass: "2A", DaysOfWeek: "TUE,FRI",
		},
		{
			TrainNumber: "12723", TrainName: "Telangana Express",
			Source: "NEW DELHI", Destination: "HYDERABAD",
			DepartureTime: "06:50", ArrivalTime: "08:55",
			TotalSeats: 700, AvailableSeats: 700,
			FarePerSeat: 770.00, TrainClass: "SL", DaysOfWeek: "ALL",
		},
	}

	result := DB.Create(&trains)
	if result.Error != nil {
		return fmt.Errorf("failed to seed trains: %w", result.Error)
	}

	slog.Info("Seeded sample trains", "count", len(trains))
	return nil
}

// Close cleanly shuts down the database connection.
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
