# NailzByDardo API

A production-style REST API for managing a nail salon business — clients, appointments, services, expenses, and financial reporting. Built from scratch in Go with a strict layered architecture, type-safe SQL, and session-based authentication.

Built as the backend for a small business management app used by a real nail tech. The operator UI is a separate React PWA: [nailzbydardo-frontend](https://github.com/anagarcia3174/nailzbydardo-frontend). This API owns business rules, money math (integer cents), and persistence; the frontend is a cookie-authenticated client that points at this server via `VITE_API_URL`.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL 18 |
| DB Driver | [pgx/v5](https://github.com/jackc/pgx) |
| SQL → Go codegen | [sqlc](https://sqlc.dev) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Password hashing | bcrypt |
| Testing | Go standard `testing` package, integration tests against a real Postgres instance |

## Architecture

The API follows a strict layered architecture — each layer only talks to the one directly beneath it:
```
HTTP Request
↓
Handler (decode requests, encode responses, map errors → HTTP status)
↓
Service (business rules, validation, money calculations)
↓
Repository (Postgres access via sqlc, type conversion)
↓
PostgreSQL
```

- **Handlers** never contain SQL or business logic — they translate HTTP ↔ Go types and delegate.
- **Services** own all business rules (e.g. a discount can't exceed 100%, an appointment can't be created for a nonexistent client) and all money math (revenue, discounts, tips) — computed in Go rather than SQL to correctly handle percentage-based discounts.
- **Repositories** are the only layer aware of SQL/pgx — they translate between sqlc's generated database types and the application's domain models, and never leak database-specific errors upward (e.g. `pgx.ErrNoRows` becomes a domain-level `ErrNotFound`).

## Project Structure
```
backend/
├── cmd/server/ entrypoint
├── internal/
│ ├── app/ dependency wiring (shared by main + tests)
│ ├── config/ environment-based configuration
│ ├── db/ connection pool + sqlc-generated code
│ ├── handler/ HTTP handlers
│ ├── middleware/ auth middleware
│ ├── model/ domain types
│ ├── repository/ data access layer
│ ├── router/ route registration
│ ├── service/ business logic layer
│ └── testutil/ shared test helpers & fixtures
├── migrations/ versioned SQL schema migrations
└── db/query/ hand-written SQL powering sqlc codegen
```

## Features

- **Session-based authentication** with HTTP-only, secure cookies (no third-party auth — manually provisioned users, since this is a single-business internal tool)
- **Client management** with soft-delete (client history is preserved even after "deletion")
- **Appointment scheduling** with per-appointment services and discounts (amount- or percentage-based)
- **Automatic revenue calculation** — subtotal, discounts, tips, and grand total computed server-side per appointment
- **Dashboard & financial reporting** — revenue, expenses, appointment counts, and tips over any date range
- **Consistent error handling** — every error maps to a predictable HTTP status (`400`, `401`, `404`, `500`) with a uniform JSON error shape

## API Overview

All routes except `/health` and `/auth/login` require a valid session cookie.

### Auth
| Method | Route | Description |
|---|---|---|
| POST | `/auth/login` | Authenticate and receive a session cookie |
| POST | `/auth/logout` | Invalidate the current session |
| GET | `/auth/me` | Check current authentication status |

### Clients
| Method | Route | Description |
|---|---|---|
| GET | `/clients` | List all active clients |
| POST | `/clients` | Create a client |
| GET | `/clients/{id}` | Get a client |
| PATCH | `/clients/{id}` | Update a client |
| DELETE | `/clients/{id}` | Soft-delete a client |
| GET | `/clients/{id}/appointments` | List a client's appointment history |
| GET | `/clients/{id}/spent` | Get a client's total spent

### Services
| Method | Route | Description |
|---|---|---|
| GET | `/services` | List services offered |
| POST | `/services` | Create a service |
| PATCH | `/services/{id}` | Update a service |
| DELETE | `/services/{id}` | Delete a service |

### Expenses
| Method | Route | Description |
|---|---|---|
| GET | `/expenses` | List expenses |
| POST | `/expenses` | Log an expense |
| DELETE | `/expenses/{id}` | Delete an expense |

### Appointments
| Method | Route | Description |
|---|---|---|
| GET | `/appointments` | List all appointments |
| POST | `/appointments` | Create an appointment |
| GET | `/appointments/upcoming` | List upcoming appointments |
| GET | `/appointments/{id}` | Get full appointment detail (services + discounts) |
| PATCH | `/appointments/{id}` | Update an appointment |
| DELETE | `/appointments/{id}` | Delete an appointment |
| GET | `/appointments/{id}/total` | Calculate the appointment's total (subtotal, discounts, tip, grand total) |
| POST | `/appointments/{id}/services` | Add a service to an appointment |
| DELETE | `/appointments/{id}/services/{serviceId}` | Remove a service from an appointment |
| POST | `/appointments/{id}/discounts` | Add a discount to an appointment |
| DELETE | `/appointments/{id}/discounts/{discountId}` | Remove a discount from an appointment |

### Reporting
| Method | Route | Description |
|---|---|---|
| GET | `/dashboard` | Current-month snapshot: revenue, tips, appointment count, expenses, upcoming appointments |
| GET | `/financials?start=...&end=...` | Revenue, expenses, appointment count, and tips over an arbitrary date range |

## Database Schema

8 tables, UUIDv7 primary keys (time-ordered, index-friendly), `TIMESTAMPTZ` throughout, and money stored as integer cents in application code to avoid floating-point rounding errors:

`users` · `sessions` · `clients` · `services` · `expenses` · `appointments` · `appointment_services` · `appointment_discounts`

## Testing

Integration tests exercise the full stack — HTTP request → handler → service → repository → a real PostgreSQL test database — covering both success and failure paths (validation errors, not-found, unauthorized, malformed requests) for each resource.

```bash
go test ./...
```

## Running Locally

### Prerequisites
- Go 1.2x+
- Docker (for PostgreSQL)

## Running Locally

Requires Go and Docker.

```bash
docker compose up -d
migrate -path migrations -database "<connection-string>" up
go run ./cmd/server
```

The server starts on the port specified in `.env` (default `8080`). Check `/health` to confirm the database connection.

To run the UI against this API, clone [nailzbydardo-frontend](https://github.com/anagarcia3174/nailzbydardo-frontend) and set `VITE_API_URL` to this origin (include credentials / CORS must allow that frontend origin).

## License

MIT