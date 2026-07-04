# Technical Evaluation Report: IRCTC (Indian Railway Catering and Tourism Corporation)

> **Portal**: [https://www.irctc.co.in](https://www.irctc.co.in)  
> **Domain**: Online Railway Ticket Booking & Tourism Services  
> **Author**: CIA-3 Project Submission  
> **Date**: March 2026

---

## Table of Contents

1. [Introduction](#1-introduction)
2. [Observable Behavior Analysis](#2-observable-behavior-analysis)
3. [System Architecture Analysis](#3-system-architecture-analysis)
4. [Performance Characteristics](#4-performance-characteristics)
5. [Scalability Limitations](#5-scalability-limitations)
6. [Reliability Risks](#6-reliability-risks)
7. [Proposed Architectural Improvements](#7-proposed-architectural-improvements)
8. [Technical Justification: Why Go for Public Service Systems](#8-technical-justification-why-go-for-public-service-systems)
9. [Conclusion](#9-conclusion)
10. [References](#10-references)

---

## 1. Introduction

IRCTC is India's sole authorized online platform for booking Indian Railways tickets. It serves as one of the highest-traffic transactional websites globally, processing over **9 lakh (900,000) ticket bookings daily** and serving approximately **3 crore (30 million) users** across India. During peak periods such as Tatkal booking windows (10:00 AM – 12:00 PM IST), the system experiences extreme traffic surges with millions of concurrent users competing for a limited number of seats.

This report conducts a structured technical evaluation of IRCTC based on observable behavior, publicly available information, and industry analysis. The evaluation covers load time analysis, request patterns, session handling, API structure, client–server interaction, and proposes concrete architectural improvements.

---

## 2. Observable Behavior Analysis

### 2.1 Load Time and Page Performance

| Metric | Observed Value | Industry Benchmark |
|--------|---------------|-------------------|
| Initial Page Load (Desktop) | 3–6 seconds | < 2 seconds |
| Time to Interactive (TTI) | 5–8 seconds | < 3.5 seconds |
| First Contentful Paint (FCP) | 1.5–3 seconds | < 1.8 seconds |
| Largest Contentful Paint (LCP) | 4–7 seconds | < 2.5 seconds |
| Total Page Weight | ~3–5 MB | < 1.5 MB |
| Number of HTTP Requests | 80–120+ | < 50 |

**Key Observations:**
- The homepage is heavily loaded with promotional banners, tourism advertisements, and multiple third-party tracking scripts that inflate page weight.
- Static assets (CSS, JS, images) are served via CDN, with **87% of static content** delivered through CDN infrastructure.
- JavaScript bundles are large and not optimally code-split, leading to high TTI.
- Multiple render-blocking resources delay First Contentful Paint.

### 2.2 Request Patterns

**Login Flow:**
1. `POST /authprovider/v1/login` — Credential submission with CAPTCHA validation
2. Server responds with session token (cookie-based) + CSRF token
3. Subsequent requests include both session cookie and CSRF header

**Train Search Flow:**
1. `POST /eticketing/protected/mapps1/altabordstnmap` — Station autocomplete
2. `POST /eticketing/protected/mapps1/avlfaremapps1` — Availability and fare query
3. Response includes seat matrix, class-wise availability, and dynamic pricing

**Booking Flow:**
1. `POST /eticketing/protected/mapps1/bookingBoardingGetDetails` — Initiate booking
2. `POST /eticketing/protected/mapps1/bookingconfirm` — Confirm and pay
3. Redirect to payment gateway → callback to confirm transaction

**Observed API Characteristics:**
- APIs follow a **REST-like** pattern but with non-standard URL naming conventions
- JSON payloads are used throughout
- Heavy reliance on server-side session state
- Anti-bot CAPTCHA is enforced at login and booking confirmation
- API responses include rate-limiting headers during peak hours

### 2.3 Session Handling

| Aspect | Implementation |
|--------|---------------|
| Session Type | Server-side sessions with cookie tokens |
| Session Duration | ~15–20 minutes idle timeout |
| CSRF Protection | Double-submit cookie pattern |
| Concurrency Control | Single active session per user enforced |
| Booking Lock | Timed lock (~10 minutes) on selected seats during payment |

**Issues Identified:**
- Server-side sessions create high memory pressure during peak loads
- Session stickiness requirements limit horizontal scaling
- Abrupt session expiry during payment flow causes failed transactions
- No WebSocket or Server-Sent Events for real-time seat status updates

### 2.4 Client–Server Interaction Model

```
┌─────────┐     HTTPS      ┌──────────┐     ┌───────────────┐
│  Client  │ ──────────────>│   CDN    │────>│  Load Balancer │
│ (Browser)│                │ (Static) │     │   (NGINX/F5)   │
└─────────┘                └──────────┘     └───────┬───────┘
                                                     │
                                            ┌────────┴────────┐
                                            │   API Gateway    │
                                            │  (Rate Limiting) │
                                            └────────┬────────┘
                                                     │
                                    ┌────────────────┼────────────────┐
                                    │                │                │
                             ┌──────┴──────┐  ┌──────┴──────┐ ┌──────┴──────┐
                             │   Auth      │  │   Search    │ │  Booking    │
                             │  Service    │  │  Service    │ │  Service    │
                             └──────┬──────┘  └──────┬──────┘ └──────┬──────┘
                                    │                │                │
                             ┌──────┴────────────────┴────────────────┴──────┐
                             │              Cache Layer (Redis/Varnish)      │
                             └──────────────────────┬───────────────────────┘
                                                    │
                             ┌──────────────────────┴───────────────────────┐
                             │         Database (Oracle / PostgreSQL)        │
                             │     Primary-Replica with Auto-Failover       │
                             └──────────────────────────────────────────────┘
```

---

## 3. System Architecture Analysis

### 3.1 Technology Stack (Inferred)

| Layer | Technology |
|-------|-----------|
| **Frontend** | Angular / React, HTML5, CSS3, JavaScript |
| **Backend** | Java (Spring Boot), Python (auxiliary services) |
| **Web Server** | NGINX (reverse proxy / load balancer) |
| **Database** | Oracle (primary transactional), PostgreSQL (auxiliary) |
| **Cache** | Redis (distributed cache), Varnish (HTTP cache) |
| **CDN** | Enterprise CDN for static asset delivery |
| **Session Store** | Server-side (likely Redis-backed) |
| **Payment Gateway** | Multiple integrations (SBI, HDFC, Paytm, UPI) |
| **Anti-Bot** | AI-powered bot detection with CAPTCHA |

### 3.2 Microservices Architecture

IRCTC has transitioned from a monolithic J2EE application to a **microservices architecture**. Key service boundaries include:

1. **Authentication Service** — Login, OTP, CAPTCHA verification
2. **User Profile Service** — Passenger management, master lists
3. **Train Search Service** — Schedule lookup, seat availability (read-heavy)
4. **Booking Service** — Seat allocation, reservation creation (write-heavy, critical path)
5. **Payment Service** — Gateway integration, transaction management
6. **Notification Service** — SMS, email, push notifications
7. **Waitlist/RAC Service** — Queue management and chart preparation

### 3.3 Database Architecture

The database layer is segmented into three logical domains:

| Domain | Data Type | Access Pattern |
|--------|-----------|---------------|
| **Inquiry DB** | Train schedules, routes, stations | Read-heavy, cacheable |
| **User Profile DB** | Account details, saved passengers | Read/Write, moderate |
| **Transactional DB** | Bookings, payments, cancellations, refunds | Write-heavy, ACID-critical |

**Data Replication**: Primary-replica setup with automatic failover for disaster recovery. The transactional database requires strong consistency guarantees, ruling out eventual consistency models for the booking critical path.

### 3.4 Caching Strategy

```
Request Flow with Caching:

Client Request
    │
    ├── Static Content? ──── YES ──> CDN (87% hit rate)
    │
    ├── Train Schedule? ──── YES ──> Varnish HTTP Cache (TTL: 15min)
    │
    ├── Seat Availability? ── YES ──> Redis Cache (TTL: 30s-60s)
    │                                 (stale reads acceptable for display)
    │
    └── Booking/Payment? ─── YES ──> Direct to Database
                                     (no caching, strong consistency)
```

---

## 4. Performance Characteristics

### 4.1 Key Metrics (FY 2024-25)

| Metric | Value |
|--------|-------|
| System Uptime | 99.86% (FY24-25), improving to 99.98% (Apr–Oct 2025) |
| Daily Bookings | 9,00,000+ tickets |
| Peak Booking Rate | 31,814 tickets/minute (record: May 22, 2025) |
| Best Single-Day Volume | 16.17 lakh tickets (March 6, 2025) |
| Target Capacity (post-upgrade) | 2,25,000 transactions/minute |
| User Login Growth | ~20% YoY increase (FY24-25 vs FY23-24) |
| Bot Traffic Mitigated | 50%+ via CDN + anti-bot |

### 4.2 Tatkal Booking Analysis (Peak Load Pattern)

The **Tatkal window** (10:00 AM for AC, 11:00 AM for Non-AC) creates the most extreme traffic spike:

```
Traffic Pattern (Tatkal Window):

Requests/sec
    │
12K ┤                    ╭──╮
    │                   ╭╯  ╰╮
 8K ┤                  ╭╯    ╰╮
    │                 ╭╯      ╰──╮
 4K ┤           ╭────╯           ╰────╮
    │     ╭────╯                      ╰────╮
 1K ┤────╯                                 ╰────────
    └──────────────────────────────────────────────
     9:30  9:45  10:00  10:15  10:30  10:45  11:00
```

**Characteristics of this spike:**
- **Thundering herd**: Millions of users refresh simultaneously at 10:00:00.000
- **Write-heavy**: Every user wants to book, not just browse
- **Time-critical**: A 100ms delay can mean the difference between confirmed and waitlisted
- **Zero tolerance for double-booking**: Must maintain ACID guarantees under extreme concurrency

### 4.3 Bottleneck Analysis

| Bottleneck | Impact | Severity |
|-----------|--------|----------|
| Database write contention | Seat allocation serialization | **Critical** |
| Server-side sessions | Memory pressure, session stickiness | **High** |
| Payment gateway latency | 3–10 second round-trip for payment confirmation | **High** |
| CAPTCHA rendering | Adds 2–5 seconds to booking flow | **Medium** |
| Large JS bundles | Slow TTI on low-end devices/networks | **Medium** |

---

## 5. Scalability Limitations

### 5.1 Vertical Scaling Ceiling

The transition from 25,000 to 2,25,000 TPS capacity required significant hardware upgrades. However, the Oracle database licensing model creates a **cost-proportional scaling ceiling** — each additional CPU core requires expensive licensing, making vertical scaling economically unsustainable beyond a point.

### 5.2 Horizontal Scaling Constraints

| Constraint | Cause |
|-----------|-------|
| Session stickiness | Server-side sessions require load balancer affinity |
| Database connection limits | Pooled connections have a hard upper bound |
| Write serialization | Seat booking requires sequential consistency |
| Vendor lock-in | Oracle-specific features limit database portability |

### 5.3 Geographic Distribution Limitations

Currently, IRCTC operates from centralized data centers. This means:
- Users in Northeast India, rural areas experience 100–300ms additional latency
- No multi-region active-active deployment for the transactional layer
- CDN helps with static content but not with API latency

---

## 6. Reliability Risks

### 6.1 Single Points of Failure

1. **Central Booking Database**: Despite replication, the write-primary is a single point for all booking transactions nationwide.
2. **Payment Gateway Dependency**: External gateway failures (SBI ePay, HDFC) cascade into booking failures.
3. **Session Store**: If Redis cluster fails, all active sessions are lost, forcing millions of re-logins during peak time.

### 6.2 Failure Modes

| Failure Mode | Probability | Impact |
|-------------|-------------|--------|
| Database deadlock during peak | Medium | Booking timeouts for thousands |
| Payment gateway timeout | Medium-High | Money deducted, ticket not issued |
| CDN outage | Low | Site effectively unreachable |
| DDoS / bot surge | Medium | Genuine users locked out |
| Session store failure | Low | Mass logout, re-authentication storm |

### 6.3 Observed User-Facing Issues

From public complaints and government reports:
- **Failed transactions with money deducted**: Payment confirmed at gateway but booking not confirmed at IRCTC's end
- **Delayed refunds**: Refund processing takes 5–7 business days
- **Cryptic error messages**: Users see "Service Unavailable" without actionable information
- **Booking confirmation delays**: During peak, confirmation emails/SMS delayed by 10–30 minutes

---

## 7. Proposed Architectural Improvements

### 7.1 Event-Driven Booking Pipeline

Replace the synchronous booking flow with an event-driven architecture:

```
Current (Synchronous):
  User → API → Reserve Seat → Process Payment → Confirm → Response
  (All in one blocking request, 10-30 second timeout)

Proposed (Event-Driven):
  User → API → Enqueue Booking Request → Immediate Acknowledgment (Booking ID)
       ↓
  Booking Worker → Reserve Seat → Emit PaymentRequired Event
       ↓
  Payment Worker → Process Payment → Emit BookingConfirmed Event
       ↓
  Notification Worker → SMS/Email/Push with final status
```

**Benefits:**
- Decouples seat reservation from payment processing
- Absorbs traffic spikes via message queue buffering
- Enables retry logic for transient payment failures
- Reduces API response time from ~30s to < 500ms

**Technology**: Apache Kafka or RabbitMQ for event streaming, with consumer groups for parallel processing.

### 7.2 CQRS (Command Query Responsibility Segregation)

Separate the read path (train search, availability) from the write path (booking):

| Path | Optimization |
|------|-------------|
| **Query (Read)** | Denormalized read replicas, aggressive caching, eventual consistency acceptable |
| **Command (Write)** | Dedicated write database, strong consistency, optimistic locking |

This allows the read path to scale independently to millions of read QPS using Redis/Elasticsearch, while the write path maintains ACID guarantees on a smaller, optimized database.

### 7.3 Circuit Breaker Pattern for Payment Gateways

```go
// Pseudo-code for circuit breaker
type PaymentCircuitBreaker struct {
    failureCount   int
    threshold      int
    state          string // "closed", "open", "half-open"
    lastFailure    time.Time
    cooldownPeriod time.Duration
}

// If SBI gateway fails 5 times in 60s → circuit opens
// Route traffic to HDFC/Paytm → retry SBI after cooldown
```

**Benefits:**
- Prevents cascading failures when one payment gateway is down
- Automatic failover to healthy gateways
- Reduces "money deducted but ticket not issued" scenarios

### 7.4 Distributed Rate Limiting with Fair Queuing

Instead of simply rejecting excess requests (current approach), implement a **fair queuing system**:

```
Traditional Rate Limiting:
  User A (fast network) → Gets tickets ✓
  User B (slow network) → Rate limited, rejected ✗

Fair Queue Rate Limiting:
  All users → Virtual Queue → Process in FIFO order
  User gets a queue position and estimated wait time
  "You are #4521 in line. Estimated wait: 45 seconds"
```

This is more equitable for users on slower networks (rural India) and reduces the "fastest finger first" problem.

### 7.5 Stateless Authentication with JWT

Replace server-side sessions with **stateless JWT tokens**:

| Aspect | Current (Sessions) | Proposed (JWT) |
|--------|-------------------|----------------|
| Storage | Server-side Redis | Client-side token |
| Scaling | Requires session affinity | Any server can handle any request |
| Memory | O(n) per active user | O(1) server-side |
| Invalidation | Instant | Requires token blocklist |

### 7.6 Progressive Web App (PWA) Optimization

- **Code splitting**: Lazy-load booking module only when needed
- **Service workers**: Cache train schedule data offline
- **Reduce JS bundle**: Target < 200KB gzipped for core bundle
- **Image optimization**: WebP format, responsive images, lazy loading

---

## 8. Technical Justification: Why Go for Public Service Systems

### 8.1 Comparison Matrix

| Feature | Go | Java (Spring) | Node.js | Python (Django) | PHP (Laravel) |
|---------|-----|---------------|---------|----------------|--------------|
| **Concurrency Model** | Goroutines (lightweight, M:N scheduling) | OS Threads (heavy, 1MB+ stack) | Event loop (single-threaded) | GIL-limited threads | Process-per-request |
| **Memory per Concurrent Unit** | ~2-8 KB (goroutine) | ~1 MB (thread) | N/A (single-threaded) | ~8 MB (process) | ~2-10 MB (process) |
| **Max Concurrent Connections** | 100K+ easily | 10K–50K (thread pool limited) | 10K+ (I/O bound only) | 1K–5K | 1K–2K |
| **Cold Start Time** | ~10ms | ~2-10 seconds (JVM warmup) | ~100-500ms | ~200ms-1s | ~50-200ms |
| **Binary Size** | ~10-20 MB (single binary) | ~200MB+ (JRE + dependencies) | N/A (requires runtime) | N/A (requires runtime) | N/A (requires runtime) |
| **CPU-bound Performance** | Excellent (compiled, no GC pauses) | Excellent (JIT, but GC pauses) | Poor (single-threaded CPU) | Poor (interpreted + GIL) | Poor (interpreted) |
| **Deployment** | Single static binary | JAR + JRE + config | node_modules + runtime | virtualenv + runtime | Composer + runtime |

### 8.2 Why Go Excels for IRCTC-like Systems

#### 1. Goroutine-based Concurrency
Go's goroutines are the killer feature for high-concurrency public services. During Tatkal booking, IRCTC handles millions of simultaneous connections. Each goroutine costs only **2-8 KB of stack space** (vs 1 MB for a Java thread), meaning a single server with 8 GB RAM can handle **1 million concurrent goroutines**.

```go
// Each booking request runs in its own goroutine
// No thread pool sizing, no callback hell
go func() {
    result := processBooking(request)
    responseChan <- result
}()
```

Java's thread-per-request model (even with virtual threads in Project Loom) still carries higher overhead. Node.js's event loop cannot parallelize CPU-bound tasks like seat allocation algorithms.

#### 2. Compiled Performance with Low Latency
Go compiles to native machine code with no JVM warmup or interpreter overhead. For a system where **100ms matters** (Tatkal booking), Go's consistent sub-millisecond response times from cold start are critical.

- No JIT compilation pauses (unlike Java)
- No Global Interpreter Lock (unlike Python)
- Predictable garbage collection with sub-millisecond STW pauses

#### 3. Single Binary Deployment
Go compiles to a **single static binary** with zero external dependencies. This dramatically simplifies deployment in government infrastructure:

```bash
# Build
GOOS=linux GOARCH=amd64 go build -o irctc-service .

# Deploy — just copy the binary
scp irctc-service prod-server:/opt/services/
```

Compare this with Java (JRE + WAR + Tomcat + config) or Python (virtualenv + pip dependencies + WSGI server).

#### 4. Built-in HTTP/2 and Networking
Go's `net/http` package is production-grade out of the box, powering systems like **Cloudflare**, **Docker**, and **Kubernetes**. No need for external web servers like Tomcat (Java) or Gunicorn (Python).

#### 5. Resource Efficiency for Government Scale
Indian government IT infrastructure often operates under budget constraints. Go's low memory footprint means:
- **Fewer servers needed**: 10 Go servers can replace 50 PHP/Python servers
- **Lower cloud costs**: Smaller instances, less RAM
- **Faster scaling**: Sub-second container startup (vs 10+ seconds for Java)

### 8.3 Real-World Go Adoption in Similar Systems

| Organization | System | Scale |
|-------------|--------|-------|
| Cloudflare | Edge proxy | Millions of requests/second |
| Uber | Microservices | ~2000+ Go microservices |
| Twitch | Chat/Video | Millions of concurrent viewers |
| Google | Core infrastructure | YouTube, Kubernetes |
| Zerodha (India) | Stock trading platform | Millions of trades/day |

**Zerodha** is particularly relevant — an Indian fintech company handling similar traffic patterns (market opening surge at 9:15 AM IST) with Go-based infrastructure, demonstrating Go's suitability for Indian high-traffic financial systems.

---

## 9. Conclusion

IRCTC has made significant strides in modernizing its infrastructure, achieving 99.98% uptime and record booking volumes. However, the analysis reveals several areas for improvement:

1. **Architectural**: Moving to event-driven booking and CQRS would better handle the extreme read/write ratio disparity during Tatkal windows.
2. **Performance**: Frontend optimization (code splitting, reduced JS payload) would significantly improve the experience for users on low-bandwidth connections.
3. **Reliability**: Circuit breakers for payment gateways and fair queuing would reduce the most common user complaints.
4. **Technology**: Adopting Go for performance-critical microservices (booking, payment processing) would provide better concurrency handling and resource efficiency compared to the current Java/Spring Boot stack.

The accompanying Go REST API implementation demonstrates these principles in practice, simulating IRCTC's ticket booking workflow with proper concurrency handling, rate limiting, and graceful shutdown.

---

## 10. References

1. Press Information Bureau, Government of India — IRCTC Technology Modernization Updates (2025)
2. Indian Railways Official Statistics — Booking Volume and Uptime Data (FY 2024-25)
3. GeeksforGeeks — IRCTC System Design Analysis
4. Millennium Post — IRCTC CDN Implementation and Anti-Bot Measures (2025)
5. Pratt Institute — IRCTC Usability Design Critique (2025)
6. GoodReturns — Indian Railways Booking Capacity Upgrade (2025)
7. Economic Times — IRCTC Technology Infrastructure Reports (2024-25)
