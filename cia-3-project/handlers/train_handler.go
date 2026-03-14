package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itcodehery/irctc-simulator/models"
	"github.com/itcodehery/irctc-simulator/services"
)

// TrainHandler handles train-related requests.
type TrainHandler struct{}

// NewTrainHandler creates a new TrainHandler.
func NewTrainHandler() *TrainHandler {
	return &TrainHandler{}
}

// ListTrains returns all available trains.
// GET /api/v1/trains
func (h *TrainHandler) ListTrains(c *gin.Context) {
	trains, err := services.GetAllTrains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to fetch trains",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":  len(trains),
		"trains": trains,
	})
}

// SearchTrains finds trains by source and destination.
// GET /api/v1/trains/search?source=...&destination=...
func (h *TrainHandler) SearchTrains(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")

	if source == "" || destination == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "source and destination are required",
			Code:    http.StatusBadRequest,
			Details: "Use query params: ?source=NEW DELHI&destination=MUMBAI CENTRAL",
		})
		return
	}

	trains, err := services.SearchTrains(source, destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "search failed",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source":      source,
		"destination": destination,
		"count":       len(trains),
		"trains":      trains,
	})
}

// GetTrain returns a specific train with its details.
// GET /api/v1/trains/:id
func (h *TrainHandler) GetTrain(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid train ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	train, err := services.GetTrainByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "train not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, train)
}

// GetAvailability returns seat availability for a specific train.
// GET /api/v1/trains/:id/availability
func (h *TrainHandler) GetAvailability(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid train ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	availability, err := services.GetTrainAvailability(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "train not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, availability)
}
