# Go Appointment Booking System

A scalable, production-ready appointment scheduling REST API built in **Go** using the **Gin** HTTP framework and **GORM** atop **PostgreSQL**. 

---

## 📖 Project Overview

This appointment scheduling API is designed around real-world dynamics: Coaches assign complex weekly scheduling windows, while Users request dynamically calculated mathematical interval slots bounded exactly by those windows. 

It handles cross-timezone safety, mathematical slot boundary generation, active collision protections against exact-millisecond double booking vulnerabilities, and uses a strict, decoupled layered architecture.

---

##  Architecture

The backend strictly implements **Clean Architecture** patterns to fully decouple business logic from framework routing and underlying data transmission dependencies.

```text
appointment-booking-system/
├── cmd/server/main.go            # Bootstraps dependencies, handles graceful exits
├── internal/
│   ├── handlers/                 # HTTP ingress: Gin context, JSON binding, error surfacing
│   ├── services/                 # Core domain logic, boundary scaling, and slot math
│   ├── repositories/             # Persistence integration handling complex GORM transactions 
│   └── models/                   # GORM Entity configurations and abstracted JSON DTOs
└── pkg/database/                 # Global Postgres pooling and automated schema migration logic
```
* **Handlers** know nothing about databases; they strictly handle HTTP statuses, payloads, and parameter conversions.
* **Services** perform high-value logic, including generating overlapping mathematical boundaries, calculating timezone deltas, and enforcing constraints. 
* **Repositories** encapsulate SQL transactions and cleanly abstract raw PostgreSQL errors (like `SQLSTATE 23505`) back up the stack as localized Go error conditions smoothly.

---

##  Database Schema

The system uses GORM's `AutoMigrate` functionality to securely enforce rigid referential architectures inside PostgreSQL. 

**Core Entities:**
1. **User**: Represents a patient or customer seeking bookings (`id`, `name`, `email`).
2. **Coach**: Represents a doctor or scheduler providing the bookings (`id`, `name`, `timezone`).
3. **Availability**: Links to a `Coach` restricting what standard structural windows (`HH:MM` intervals) they accept meetings on categorized by `day_of_week`.
4. **Booking**: The finalized intersecting reservation safely connecting a `User` to a `Coach` at a specific `.SlotTime`.

### Security Guarantees
* **Referential Data Integrity**: Cannot assign availabilities to non-existent coaches.
* **Database Layer Collision Preventions**: The `Booking` table utilizes a dynamic composite constraint on `(coach_id, slot_time)`. This enforces the "one meeting per professional per 30-min block" logic *before* software race conditions can theoretically leak it.

---

##  API Endpoints

Below is a quick overview of primary system functionality. Please refer to [`API_DOCS.md`](./API_DOCS.md) for deeper technical specifications on bodies and response schemas.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Standard network and database liveness ping. |
| `POST` | `/api/v1/coaches/availability` | Creates standard daily recurring working hours for a professional. |
| `GET` | `/api/v1/users/slots?coach_id=1&date=Y-M-D` | Calculates and lists openly available 30-minute availability blocks for a specified day. |
| `POST` | `/api/v1/users/bookings` | Registers a booking. Throws `409 Conflict` dynamically on race conditions. |
| `GET` | `/api/v1/users/bookings?user_id=101` | Retrieves standard active scheduling logs for a user. |
| `DELETE` | `/api/v1/users/bookings/{id}` | Releases a reserved slot back into the active pool. |

---

##  Setup Instructions

1. **Prerequisites**
   * Go runtime (v1.21+)
   * PostgreSQL database service (v14+) running locally or remotely

2. **Clone the Repository**
   ```bash
   git clone https://github.com/your-org/appointment-booking.git
   cd appointment-booking
   ```

3. **Configure Environment**
   Duplicate `.env.example` to `.env` and assign your PostgreSQL database configuration:
   ```dotenv
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=appointment_db
   ```

4. **Initialize the Database**
   Ensure your target database schema exists prior to startup:
   ```bash
   psql -U postgres -c "CREATE DATABASE appointment_db;"
   ```

---

##  Run Instructions

Since all structural modeling dynamically executes on standard Go application launches, you simply need to build/boot natively:

```bash
go run cmd/server/main.go
```

*(Once running, the system will systematically walk all `Models`, executing automated non-destructive creation/migrations across your entire schema table structure immediately prior to binding onto port `:8080`.)*

### Background Seeder Tooling
*If utilizing a starkly empty database environment, run `go run seed.go` to inject one dummy user (`101`) and one dummy coach (`1`) required natively to bypass default strict Postgres Foreign-Key constraints!*

---

##  AI Usage Disclosure

**1. Complete Conversation Logs**
* The full, raw conversation history (prompts + responses) between the human developer and the AI has been exported and is included alongside this submission.

**2. Which parts of the solution were AI-assisted?**
This codebase was a collaborative effort, split roughly **40% human** and **60% AI assistance** utilizing the Deepmind Antigravity Agentic AI acting as a pair programmer. 

**Human Contribution (40%):**
* Overall feature architecture design and database relationship planning.
* Business logic validation.
* Directing the Clean Architecture grouping (Handlers, Services, Repositories).
* Environmental orchestration and execution of Postman validation loops.

**AI Contribution (60%):**
* Direct scaffolding of GORM database schema structs and unique indices.
* Implementation of explicit HTTP logic routing matching the REST specification.
* Authoring mathematical timezone interval isolation logic across the service layer.
* Automatic generation of Swagger-style `README.md` and `API_DOCS.md` documentation.