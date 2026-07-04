# IRCTC Ticket Booking Simulator - System Architecture

## Overview
The IRCTC Ticket Booking Simulator is a high-performance backend application built with Go, designed to handle concurrent railway ticket bookings with precision and reliability. It employs a layered architecture and a worker-pool pattern to manage high-throughput booking requests while preventing overbooking.

## 🏗️ Architecture Layers

### 1. Presentation Layer (Handlers)
- **Location:** `handlers/`
- **Responsibility:** Manages HTTP request/response lifecycles using the **Gin** web framework.
- **Key Handlers:**
  - `UserHandler`: Handles registration, login, and profile management.
  - `TrainHandler`: Manages train browsing, searching, and availability checks.
  - `BookingHandler`: Coordinates ticket booking, listing, and cancellation.

### 2. Business Logic Layer (Services)
- **Location:** `services/`
- **Responsibility:** Implements core application logic and ensures data consistency.
- **Key Service:**
  - `BookingService`: The heart of the application. It uses a **Worker Pool** and **Job Queue** to process booking requests asynchronously from the API but synchronously for data integrity. It employs **Pessimistic Locking** (`SELECT ... FOR UPDATE`) in database transactions to prevent race conditions during seat allocation.

### 3. Data Access Layer (Database & Models)
- **Location:** `database/`, `models/`
- **Responsibility:** Data persistence and schema definition.
- **Technology:** Uses **GORM** (Object Relational Mapper) with an **SQLite** database.
- **Models:** Includes `User`, `Train`, and `Booking` entities with appropriate relationships.

### 4. Middleware Layer
- **Location:** `middleware/`
- **Responsibility:** Cross-cutting concerns that apply to multiple endpoints.
- **Key Components:**
  - `AuthMiddleware`: Validates JWT tokens for protected routes.
  - `RateLimiter`: Prevents API abuse by limiting requests per IP address using a token bucket-like algorithm.
  - `StructuredLogger`: Provides JSON-formatted logs for better observability.

## 📦 Key Dependencies & Their Usage

| Dependency | Purpose | Exact Usage in App |
|------------|---------|-------------------|
| **Gin Gonic** (`v1.12.0`) | Web Framework | Used in `router/router.go` to define API endpoints and in `handlers/` for request processing. |
| **GORM** (`v1.25.12`) | ORM | Used in `database/database.go` for connection and in `services/` for all CRUD operations. |
| **JWT-Go** (`v5.3.1`) | Authentication | Used in `handlers/user_handler.go` for issuing tokens and `middleware/auth.go` for verification. |
| **Bcrypt** (`golang.org/x/crypto`) | Security | Used in `models/models.go` and `handlers/user_handler.go` to safely hash and compare user passwords. |
| **Slog** (Standard Library) | Logging | Used throughout the application for structured, high-performance logging. |

## ⚙️ Core Workflows

### 🛤️ Ticket Booking Flow
1. **Request:** User sends a POST request to `/api/v1/bookings`.
2. **Authentication:** `AuthMiddleware` verifies the JWT token.
3. **Queueing:** `BookingHandler` submits a `BookingJob` to the `BookingService`'s internal job queue.
4. **Worker Processing:** A background worker picks up the job.
5. **Transaction:**
   - Starts a database transaction.
   - Locks the specific Train row (`FOR UPDATE`).
   - Verifies seat availability.
   - Generates PNR and seat numbers.
   - Creates the Booking record.
   - Updates the Train's `available_seats`.
   - Commits the transaction.
6. **Response:** The worker sends the result back via a channel, and the handler returns the confirmed booking to the user.

### 🛡️ Concurrency & Integrity
The system is designed to never allow double-booking even if thousands of requests hit the same train simultaneously. This is achieved through:
- **Limited Workers:** Only a fixed number of workers (default: 10) process bookings at once.
- **Database Atomicity:** GORM transactions ensure that "check-then-act" operations are atomic.
- **Row-Level Locking:** Prevents other transactions from reading or modifying the same train's availability until the current booking is complete.

## 🚀 Future Improvements
- **Redis Integration:** For more scalable rate limiting and as a job queue backend.
- **Microservices:** Splitting Train and Booking into separate services.
- **WebSockets:** Real-time updates for seat availability changes.
