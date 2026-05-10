# judgenot0 Backend — Agent Guide

This document is written for AI coding agents working on the `judgenot0/server` repository. The reader is assumed to know nothing about the project.

---

## Project Overview

`judgenot0` is a competitive programming judge backend written in Go. It provides:

- User authentication and role-based access control (RBAC) with `user`, `setter`, and `admin` roles.
- Contest lifecycle management (create, list, update contests).
- Problem and testcase management.
- Code submission handling with asynchronous judging via RabbitMQ.
- Real-time contest standings with penalty tracking and CSV export.
- Communication with remote code-execution clusters ("engines") via HMAC-signed webhooks.
- Bulk user generation via CSV.

The project currently has **99 `.go` files** and **zero test files**.

---

## Technology Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.25 |
| HTTP Router | `net/http` `ServeMux` (stdlib) |
| Database | PostgreSQL |
| ORM | GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`) |
| Authentication | JWT (`github.com/golang-jwt/jwt/v5`) |
| Password Hashing | `golang.org/x/crypto/bcrypt` |
| Message Queue | RabbitMQ (`github.com/rabbitmq/amqp091-go`) |
| Configuration | Environment variables (loaded via `github.com/joho/godotenv`) |

---

## Build and Run Commands

### Local Development

Requires Go 1.25+, a running PostgreSQL instance, and optionally RabbitMQ.

```bash
# 1. Install dependencies
go mod tidy

# 2. Create a .env file (see Configuration section below)
cp .env.example .env

# 3. Run the application
go run main.go

# 4. Or build a binary
go build -o judgeserver .
./judgeserver
```

The server listens on `0.0.0.0:8000` by default.

### Docker Deployment

A multi-stage `Dockerfile` is provided. The final image is Alpine-based and exposes port `8000`.

```bash
# Create the external network expected by docker-compose.yml
docker network create shared-net

# Build and run
docker compose up -d --build
```

**Important Docker notes:**
- The `docker-compose.yml` mounts `./generated_csv` into the container and expects an external network named `shared-net`.
- When running inside Docker, the `config/config.go` file uses `godotenv.Load()`. If env vars are injected by the orchestrator instead of a `.env` file, comment out the `env.Load()` call in `config/config.go`.
- Set `DB_HOST` to the database service name (e.g., `db`) rather than `localhost` when running in Docker.

---

## Configuration

All configuration is read from environment variables. The `.env.example` file in the project root shows the required variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `HTTP_PORT` | Yes | Server port (e.g., `8000`) |
| `JWT_SECRET` | Yes | Secret key for signing JWTs |
| `ENGINE_KEY` | Yes | Shared secret for HMAC verification of engine callbacks |
| `ENGINE_URL` | Yes | URL of the remote execution engine |
| `DB_URL` | Yes | PostgreSQL connection string |
| `QUEUE_NAME` | No | RabbitMQ queue name (default: `judge_queue`) |
| `RABBITMQ_URL` | No | RabbitMQ AMQP URL (default: `amqp://guest:guest@localhost:5672/`) |

Configuration is loaded once at startup via `config.GetConfig()` and stored in a singleton. Missing required variables cause `log.Fatalln` and process exit.

---

## Code Organization

```
cmd/
  serve.go              # Application bootstrap: wires DB, queue, middlewares, routes, and starts HTTP server

config/
  config.go             # Singleton config loader from env vars

handlers/               # Domain-driven HTTP handlers (one package per domain)
  cluster/              # Execution node registry
  compile_run/          # Direct compiler bridging
  contest/              # Contest CRUD
  contest_problems/     # Assigning problems to contests
  problem/              # Problem and testcase CRUD
  setter/               # Setter-specific dashboards
  standings/            # Leaderboards and CSV exports
  submissions/          # Submission creation, listing, engine callback updates
  user_csv/             # Bulk user generation from CSV
  users/                # Authentication and user management

infra/                  # Infrastructure concerns
  db/
    connection.go       # GORM PostgreSQL connection initializer
  queue/
    queue.go            # RabbitMQ queue struct and lifecycle
    connection.go       # Connection, reconnection, DLX/DLQ setup
    publisher.go        # Message publishing with retry
    dlq.go              # Dead-letter queue background processor

middlewares/
  manager.go            # Middleware chaining manager
  middlewares.go        # Middlewares struct factory
  auth.go               # JWT bearer-token authentication
  auth_admin.go         # Admin role guard
  auth_setter.go        # Setter role guard
  auth_engine.go        # HMAC engine webhook verification
  cors.go               # CORS headers
  logger.go             # Request logging middleware
  preflight.go          # OPTIONS preflight handling
  types.go              # JWT payload definitions

models/                 # GORM code-first schema definitions
  user.go
  contest.go
  problem.go
  testcase.go
  submission.go
  contest_problems.go
  contest_problem_result.go
  usercreds.go

utils/
  sendResponse.go       # Standard JSON API response helper

generated_csv/          # Runtime output directory for generated CSV files
```

### Handler Pattern

Each handler package follows a consistent pattern:

1. **`types.go`** — Defines request/response DTOs and the `Handler` struct (which holds `*gorm.DB`, `*config.Config`, and optionally `*queue.Queue`).
2. **`handler.go`** — Factory function `NewHandler(...)` that constructs the `Handler`.
3. **`routes.go`** — `RegisterRoutes(mux, manager, middlewares)` method wiring endpoints to handler methods with per-route middleware chains.
4. **Individual endpoint files** — Named after the endpoint (e.g., `login.go`, `create_submission.go`).

Routes are registered in `cmd/serve.go` using the middleware manager:

```go
mux.Handle("GET /api/users/{contestId}", manager.With(h.GetUsers, middlewares.Authenticate, middlewares.AuthenticateAdmin))
```

### Middleware Chaining

`middlewares.Manager` supports two levels:
- **Global middlewares** (`manager.Use(...)`) — Applied to all requests via `manager.WrapMux(mux)` (CORS, preflight, logger).
- **Per-route middlewares** (`manager.With(handler, mw1, mw2, ...)`) — Applied in reverse order around the specific handler.

---

## Database Architecture

This project uses a **code-first** approach with GORM:

- **No SQL migration files.** Schema changes are made by editing structs in `models/`.
- **Auto-migration** runs on every startup in `infra/db/connection.go` via `db.AutoMigrate(...)`.
- The migration also enables the `pgcrypto` PostgreSQL extension and seeds a default `admin` user if none exists:
  - Username: `admin`
  - Password hash: `$2a$12$Ncde3vjx7AbBXwyDlzgN5ue8PKgD1XexbvWdityKLbQHsHJAi1jKG`
  - Role: `admin`

### Key Models

- `User` — Supports `user`, `setter`, `admin` roles. `AllowedContest` restricts a user to a single contest.
- `Contest` — Has `StartTime`, `EndTime`, `DurationSeconds`, and a unique `UserPrefix`.
- `Problem` — Contains statement, time/memory limits, and checker configuration.
- `Testcase` — Linked to a problem. `IsSample` flag distinguishes sample testcases.
- `Submission` — Links user + problem + optional contest. Tracks `status`, `exec_time`, `exec_memory`.
- `ContestProblem` — Junction table with an `Index` column for ordering.
- `ContestProblemResult` — Denormalized standings data (solved, wrong attempts, penalty, first blood).
- `UserCreds` — Stores plaintext passwords for contest-generated users (used for credential distribution).

---

## Message Queue (RabbitMQ)

Submissions are pushed to a RabbitMQ queue for asynchronous judging by remote execution engines.

### Queue Setup

- Queue type: **quorum** queue.
- DLX (dead-letter exchange) and DLQ (dead-letter queue) are declared automatically on connection.
- A background goroutine (`StartDLQProcessor`) runs a 30-second ticker to requeue failed messages from the DLQ.
- Messages are retried up to **5 times** (tracked via `x-retry-count` header) before being permanently dropped.
- Publisher reconnects with exponential backoff (max 30 seconds) if the channel is closed.

### Submission Flow

1. User creates a submission → validated against contest time window and problem assignment.
2. Submission is inserted into PostgreSQL.
3. Submission payload (source code, testcases, language) is JSON-marshaled and published to RabbitMQ.
4. The execution engine consumes the message, runs the code, and calls `PATCH /api/submissions` with an HMAC-signed result.
5. `UpdateSubmission` validates the HMAC, updates the DB row, and (in the future) triggers standings recalculation.

---

## Authentication & Authorization

### User Authentication

- Login endpoint: `POST /api/users/login`
- JWT is returned in an **HttpOnly, Secure, SameSite=Lax** cookie named `access_token`.
- Token expiry: **3 hours**.
- The `Authenticate` middleware extracts the token from the `Authorization: Bearer <token>` header and injects a `Payload` into the request context under key `"user"`.

### Role Guards

- `AuthenticateAdmin` — Checks `Payload.Role == "admin"`.
- `AuthenticateSetter` — Checks `Payload.Role == "setter"`.

### Engine Authentication

- Execution engines authenticate via HMAC-SHA256.
- The engine signs a JSON payload containing `submission_id`, `verdict`, `execution_time`, `execution_memory`, and `timestamp`.
- The `AuthEngine` middleware verifies the HMAC against `ENGINE_KEY` and rejects requests older than 5 minutes.
- Valid statuses from engines: `ACCEPTED`, `WRONG_ANSWER`, `TIME_LIMIT_EXCEEDED`, `RUNTIME_ERROR`, `MEMORY_LIMIT_EXCEEDED`, `COMPILATION_ERROR`.

---

## Code Style Guidelines

The codebase does not enforce a linter or formatter, but the following conventions are observed:

- **Go version:** 1.25+ syntax (e.g., `net/http` verb+path route patterns like `GET /api/users/{contestId}`).
- **Handler methods:** Use pointer receivers on `*Handler`.
- **JSON responses:** Always use `utils.SendResponse(w, statusCode, message, data)` for consistency.
- **Request body limits:** Handlers often set `http.MaxBytesReader` to prevent large payloads (e.g., `50 * 1024` for submissions, `1024` for login).
- **Error handling:** Errors are logged with `log.Println(...)` and the user receives a generic message via `SendResponse`.
- **Database access:** A mix of GORM methods (`db.Create`, `db.Model(...).Updates`) and raw SQL (`db.Raw(...)`) is used. When adding new queries, follow the existing style in the handler package.
- **Variable naming:** The codebase uses both camelCase (`current_time`) and snake_case inconsistently. Prefer **camelCase** for new code to align with Go conventions.

---

## Testing Instructions

**There are currently no tests in this repository.**

If you add tests:

- Use the standard Go testing package (`testing`).
- Place test files next to the source files (`*_test.go`).
- For handler tests, you will need to spin up a test PostgreSQL instance or mock the GORM layer.
- There is no existing test database configuration or test helper infrastructure.

---

## Security Considerations

- **Passwords:** Stored as bcrypt hashes. Plaintext passwords are stored in `UserCreds` **only** for contest-generated accounts (for credential distribution).
- **JWT Secret:** Must be strong and kept secret. It is used for both signing and validation.
- **Engine Key:** Shared secret between the backend and execution engines. Compromise allows forged submission results.
- **Cookies:** The `access_token` cookie is `HttpOnly` and `Secure`. Ensure HTTPS is enabled in production.
- **CORS:** Currently allows `*` origin. Restrict this in production deployments.
- **Body Limits:** Several endpoints cap request body size. Verify caps are appropriate when adding new endpoints.
- **SQL Injection:** Raw SQL queries use parameterized placeholders (`$1`, `$2`). Do not introduce string-concatenated SQL.

---

## Deployment Notes

- The binary is built with `CGO_ENABLED=0` and static linking for Alpine compatibility.
- The `generated_csv` directory is mounted as a volume for persistent CSV exports.
- The application expects RabbitMQ to be reachable at the configured `RABBITMQ_URL`.
- Standings are cached in-memory (`Last_standings` map) with a background goroutine that evicts stale entries every hour.
