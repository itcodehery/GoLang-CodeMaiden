package services

import (
	"fmt"
	"strings"

	"github.com/itcodehery/irctc-simulator/database"
	"github.com/itcodehery/irctc-simulator/models"
)

// SearchTrains finds trains matching the given source and destination.
func SearchTrains(source, destination string) ([]models.Train, error) {
	var trains []models.Train

	result := database.DB.Where(
		"UPPER(source) = ? AND UPPER(destination) = ?",
		strings.ToUpper(source),
		strings.ToUpper(destination),
	).Find(&trains)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to search trains: %w", result.Error)
	}

	return trains, nil
}

// GetAllTrains returns all available trains.
func GetAllTrains() ([]models.Train, error) {
	var trains []models.Train

	result := database.DB.Find(&trains)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch trains: %w", result.Error)
	}

	return trains, nil
}

// GetTrainByID retrieves a specific train by its ID.
func GetTrainByID(id uint) (*models.Train, error) {
	var train models.Train

	result := database.DB.First(&train, id)
	if result.Error != nil {
		return nil, fmt.Errorf("train not found: %w", result.Error)
	}

	return &train, nil
}

// GetTrainAvailability returns seat availability for a specific train.
func GetTrainAvailability(trainID uint) (map[string]interface{}, error) {
	train, err := GetTrainByID(trainID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"train_number":    train.TrainNumber,
		"train_name":      train.TrainName,
		"source":          train.Source,
		"destination":     train.Destination,
		"class":           train.TrainClass,
		"total_seats":     train.TotalSeats,
		"available_seats": train.AvailableSeats,
		"booked_seats":    train.TotalSeats - train.AvailableSeats,
		"fare_per_seat":   train.FarePerSeat,
		"occupancy_pct":   float64(train.TotalSeats-train.AvailableSeats) / float64(train.TotalSeats) * 100,
	}, nil
}
