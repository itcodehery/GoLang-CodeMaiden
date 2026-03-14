# IRCTC Ticket Booking Simulator

> Built by Hari Prasad B K | 2547120
**CIA-3 Project**: Technical Evaluation of a High-Traffic Indian Public Service Portal with a High-Concurrency Go REST API

A high-concurrency REST API service built in **Go** that simulates the IRCTC (Indian Railway Catering and Tourism Corporation) ticket booking workflow. This project demonstrates concurrency handling, structured logging, rate limiting, JWT authentication, database connection pooling, and graceful shutdown.

---

## Project Structure

```
cia-3-project/
├── main.go                    # Application entry point with graceful shutdown
├── config/
│   └── config.go              # Environment-based configuration
├── models/
│   └── models.go              # GORM models and DTOs
├── database/
│   └── database.go            # SQLite + GORM with connection pooling
├── middleware/
│   ├── auth.go                # JWT authentication middleware
│   ├── logger.go              # Structured logging (slog)
│   └── ratelimiter.go         # Token bucket rate limiter
├── handlers/
│   ├── user_handler.go        # Register, login, profile
│   ├── train_handler.go       # List, search, availability
│   └── booking_handler.go     # Book, cancel, list bookings
├── services/
│   ├── booking_service.go     # Worker pool for concurrent bookings
│   └── train_service.go       # Train search and availability
├── router/
│   └── router.go              # Gin router with middleware chain
├── tests/
│   └── api_test.go            # Integration tests
├── docs/
│   └── ANALYSIS.md            # Technical evaluation report
└── go.mod / go.sum
```

## Key Features Demonstrated

| Feature | Implementation |
|---------|---------------|
| **Concurrency Handling** | Worker pool pattern with 10 goroutines processing bookings from a buffered channel |
| **Structured Logging** | Go's `slog` package with JSON output, log-level based on HTTP status |
| **Rate Limiting** | Token bucket algorithm per IP with automatic refill and stale entry cleanup |
| **Database + Connection Pooling** | SQLite via GORM with WAL mode, configurable max connections |
| **Authentication** | JWT (HS256) with Bearer token, protected route groups |
| **Graceful Shutdown** | Signal handling (SIGINT/SIGTERM) → stop server → drain workers → close DB |
| **Double-Booking Prevention** | Pessimistic row locking (`SELECT ... FOR UPDATE`) within transactions |

## Quick Start

### Prerequisites
- Go 1.21+ installed
- GCC (for SQLite CGo compilation)

### Run the Server

```bash
# Clone and navigate
cd cia-3-project

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

The server starts at `http://localhost:8080`.

### Environment Variables (Optional)

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DB_PATH` | `irctc.db` | SQLite database file path |
| `JWT_SECRET` | (built-in) | JWT signing secret |
| `RATE_LIMIT_REQUESTS` | `100` | Max requests per window per IP |
| `BOOKING_WORKERS` | `10` | Number of booking worker goroutines |
| `BOOKING_QUEUE_SIZE` | `1000` | Booking request queue buffer size |

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check with uptime |
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and get JWT token |
| `GET` | `/api/v1/trains` | List all trains |
| `GET` | `/api/v1/trains/search?source=...&destination=...` | Search trains |
| `GET` | `/api/v1/trains/:id` | Get train details |
| `GET` | `/api/v1/trains/:id/availability` | Check seat availability |

### Protected Endpoints (Require `Authorization: Bearer <token>`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/user/profile` | Get user profile |
| `POST` | `/api/v1/bookings` | Book tickets |
| `GET` | `/api/v1/bookings` | List user's bookings |
| `GET` | `/api/v1/bookings/:pnr` | Get booking by PNR |
| `DELETE` | `/api/v1/bookings/:pnr` | Cancel booking |
| `GET` | `/api/v1/bookings/queue/status` | Check booking queue |

## Usage Examples

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User",
    "phone": "9876543210"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
# Save the "token" from the response
```

### 3. Search Trains

```bash
curl "http://localhost:8080/api/v1/trains/search?source=NEW%20DELHI&destination=MUMBAI%20CENTRAL"
```

### 4. Book a Ticket

```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "train_id": 3,
    "passenger_name": "Test User",
    "num_seats": 2,
    "journey_date": "2026-03-15"
  }'
```

### 5. Check Booking

```bash
curl http://localhost:8080/api/v1/bookings/YOUR_PNR \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 6. Cancel Booking

```bash
curl -X DELETE http://localhost:8080/api/v1/bookings/YOUR_PNR \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Concurrency Testing

To test concurrent booking (simulating Tatkal rush), you can fire multiple requests simultaneously:

```bash
# Book 20 tickets concurrently for the same train
for i in $(seq 1 20); do
  curl -s -X POST http://localhost:8080/api/v1/bookings \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_TOKEN" \
    -d "{\"train_id\":1,\"passenger_name\":\"Passenger $i\",\"num_seats\":1,\"journey_date\":\"2026-03-15\"}" &
done
wait
```

The worker pool ensures no double-booking occurs even under extreme concurrency.

## Testing

```bash
# Run all tests
go test ./... -v

# Run with race detector
go test -race ./...

# Build verification
go build -o irctc-simulator .
```
