# judgenot0 Backend

A robust, high-performance competitive programming judge backend written in Go. This system handles user authentication, problem setting, test case management, code submissions, automated standings, and communication with remote code-execution clusters.

## ✨ Features

* **Code-First ORM Database:** Fully modernized purely in Go using [GORM](https://gorm.io). No manual `.sql` migration files to maintain!
* **Automated Migrations & Seeding:** The database structure syncs automatically on startup and self-seeds a default admin account.
* **Role-Based Authentication (RBAC):** Distinct roles for `user`, `setter`, `admin`, and secure tokens for `engine` execution clusters.
* **Real-time Standings:** Advanced aggregations, denormalization, and penalty tracking for fast contest leaderboards.
* **Problem Management:** Pluggable problem execution, checker precision types, resource limits, and test case processing.

## 🛠️ Technology Stack

* **Language:** Go (1.24+)
* **Routing:** `net/http` ServeMux (Go Standard Library)
* **Database:** PostgreSQL
* **ORM:** `gorm.io/gorm` & `gorm.io/driver/postgres`
* **Authentication:** JWT (JSON Web Tokens)

## 📋 Prerequisites

* **Go 1.24** or higher installed on your machine.
* **PostgreSQL** instance running locally or remotely.
* **Docker / Docker Compose** (Optional, for containerized deployments).

## ⚙️ Configuration & Setup

1. **Clone the repository:**
   ```bash
   git clone <your-repo-url>
   cd judgenot0/server
   ```

2. **Setup your environment variables:**
   Create a `.env` file in the root of the project with the following required variables:
   ```env
   # Server Configuration
   HTTP_PORT=8000
   JWT_SECRET=your_super_secret_jwt_key
   
   # Execution Engine Configuration
   ENGINE_KEY=secret_key_to_communicate_with_execution_cluster
   ENGINE_URL=http://localhost:8080/api/internal
   
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_postgres_password
   DB_NAME=judgenot0
   DB_SSL_MODE=disable # set to 'true' in production
   ```

3. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the Application:**
   ```bash
   go run main.go
   # Or alternatively:
   # go build . && ./judge-backend
   ```

## 🗄️ Database Architecture & Migrations

This project heavily leans on the **Code-First** approach. 

* **No Schema Directory:** You do not need to configure tools like `sql-migrate` or apply raw SQL scripts. 
* **Auto-Migration:** Upon startup, GORM invokes `AutoMigrate()` in `infra/db/connection.go`. This automatically synchronizes all the structures defined inside the `models/` directory natively with PostgreSQL (including composite indexes, primary keys, and relations).
* **Default Seeding:** When a migration runs for the completely fresh database, the backend checks for the presence of an `admin` account. If absent, it automatically seeds the initial administrator profile:
   * **Username:** `admin`
   * **Role:** `admin`
   * *(Password matches the original schema hash you injected).*

## 🐳 Docker Deployment

The application includes a `Dockerfile` and `docker-compose.yml` for straightforward deployments.

1. **Ensure your `.env` file exists** in the root directory.
2. **Build and run the container:**
   ```bash
   docker-compose up -d --build
   ```
   *Note: Ensure your `DB_HOST` in the `.env` file points to a reachable PostgreSQL instance from inside the Docker container (e.g. your host network IP rather than `localhost`).*

## 📁 Project Structure

```text
├── cmd/                # Entrypoint orchestrator (serve.go)
├── config/             # Application environment configurations (.env parser)
├── handlers/           # Core API modules (grouped by domain)
│   ├── cluster/        # Execution node registry
│   ├── compile_run/    # Direct compiler bridging 
│   ├── contest/        # Contest lifecycle (CRUD)
│   ├── contest_problems/
│   ├── problem/        # Statement & limit settings
│   ├── setter/         # Setter dashboards
│   ├── standings/      # Leaderboards and CSV exports
│   ├── submissions/    # Asynchronous scoring and queues
│   ├── user_csv/       # Bulk user generation
│   └── users/          # Authentication & RBAC maps
├── infra/              # Core infrastructure modules
│   └── db/             # GORM connection and migration initializers
├── middlewares/        # Request intercepts (CORS, Auth, Logs, Preflight)
├── models/             # GORM code-first Database Schema declarations
├── utils/              # Helper functions (HTTP responses, JWT generation)
├── Dockerfile          # Container environment specification
├── docker-compose.yml  # Local deploy configuration
└── main.go             # Application bootstrap 
```

## 🤝 Contributing
For further iterations, modify the database schemas strictly within the `models/*.go` structs. Restart the backend to automatically propagate table updates. Avoid modifying fields mapped out manually by the underlying relational driver directly via `psql`.