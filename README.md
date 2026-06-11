# Appointment Booking System

A backend system for managing appointments, built in Go. The focus throughout was on correctness, observability, and making deliberate decisions about things that tend to go wrong in real systems - concurrent writes, unprotected routes, opaque errors, and ungraceful deploys.

## Design Decisions

### Authentication & Authorization
Token-based auth with a full user activation flow - users receive an emailed link on sign-up and cannot access resources until activated. Authorization is permission-based (RBAC) rather than role-only, giving finer control over what each user can do. Auth middleware is applied selectively so public routes aren't burdened unnecessarily.

### Data Integrity
Updates use optimistic locking via a DB version field, so concurrent writes fail loudly rather than silently overwriting each other. Duplicate creation is caught and returns a clear error rather than a generic 500. UUIDs are used for staff identifiers instead of auto-increment integers.

### Resilience
The server handles panics gracefully, returns structured errors, and shuts down cleanly - draining in-flight requests before exiting. Rate limiting is applied globally to protect against abuse.

### Observability
Prometheus instrumentation tracks HTTP request counts, error rates, and in-flight requests out of the box. All logs are structured JSON, making them queryable in any log aggregation setup.

### Configuration
All environment-specific values - database URL, frontend origin, SMTP credentials - are loaded from environment variables. Nothing is hardcoded.

## API Overview

| Resource        | Endpoints                        |
|----------------|----------------------------------|
| Users           | POST /users, PUT /users/activate |
| Auth Tokens     | POST /tokens/authentication      |
| Staff           | CRUD /staff                      |
| Service Types   | CRUD /service-types              |
| Services        | CRUD + search/filter /services   |
| Appointments    | CRUD /appointments               |
| Metrics         | GET /metrics (Prometheus)        |

Services support full-text search, dynamic filtering, sorting, and pagination.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Auth:** Token-based + bcrypt password hashing (`x/crypto/bcrypt`)
- **Metrics:** Prometheus (`prometheus/client_golang`)
- **Rate Limiting:** `golang.org/x/time/rate`

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL

### Setup

```bash
git clone https://github.com/Jane-Mwangi/nailit-api
cd nailit
```

Set your environment variables:

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname
FRONTEND_URL=http://localhost:5173

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_username
SMTP_PASSWORD=your_password
```

Run migrations and start the server:

```bash
make db/migrations/up
make run/api
```

