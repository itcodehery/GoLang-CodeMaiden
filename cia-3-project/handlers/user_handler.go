package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itcodehery/irctc-simulator/config"
	"github.com/itcodehery/irctc-simulator/database"
	"github.com/itcodehery/irctc-simulator/middleware"
	"github.com/itcodehery/irctc-simulator/models"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user registration, login, and profile requests.
type UserHandler struct {
	Config *config.Config
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{Config: cfg}
}

// Register creates a new user account.
// POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request body",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to process registration",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "registration failed",
			Code:    http.StatusConflict,
			Details: "Username or email already exists",
		})
		return
	}

	// Generate JWT
	token, expiresAt, err := middleware.GenerateToken(&user, h.Config.JWTExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to generate token",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	})
}

// Login authenticates a user and returns a JWT.
// POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request body",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "invalid username or password",
			Code:  http.StatusUnauthorized,
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "invalid username or password",
			Code:  http.StatusUnauthorized,
		})
		return
	}

	// Generate JWT
	token, expiresAt, err := middleware.GenerateToken(&user, h.Config.JWTExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "failed to generate token",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	})
}

// GetProfile returns the authenticated user's profile.
// GET /api/v1/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := database.DB.Preload("Bookings").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "user not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
