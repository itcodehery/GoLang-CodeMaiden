package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	ServerPort        string
	DBPath            string
	JWTSecret         string
	JWTExpiry         time.Duration
	RateLimitRequests int
	RateLimitWindow   time.Duration
	MaxDBConnections  int
	BookingWorkers    int
	BookingQueueSize  int
	ShutdownTimeout   time.Duration
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8081"),
		DBPath:            getEnv("DB_PATH", "irctc.db"),
		JWTSecret:         getEnv("JWT_SECRET", "irctc-simulator-secret-key-change-in-production"),
		JWTExpiry:         getDurationEnv("JWT_EXPIRY", 24*time.Hour),
		RateLimitRequests: getIntEnv("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getDurationEnv("RATE_LIMIT_WINDOW", 1*time.Minute),
		MaxDBConnections:  getIntEnv("MAX_DB_CONNECTIONS", 25),
		BookingWorkers:    getIntEnv("BOOKING_WORKERS", 10),
		BookingQueueSize:  getIntEnv("BOOKING_QUEUE_SIZE", 1000),
		ShutdownTimeout:   getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
