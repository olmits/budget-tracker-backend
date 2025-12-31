# Go + Next.js Budget Tracker

A high-performance, type-safe Budget Tracker application built with a **Golang REST API** backend and a **Next.js** frontend. This application allows users to track income, expenses, and monitor monthly budget limits in real-time.

## 🛠 Tech Stack

### Backend (The Core)
* **Language:** Golang (1.21+)
* **Architecture:** RESTful API with Layered Architecture (Handler -> Service -> Repository)
* **Router:** [Gin](https://github.com/gin-gonic/gin) (High-performance HTTP web framework)
* **Database:** PostgreSQL 15+
* **DB Driver:** [pgx/v5](https://github.com/jackc/pgx) (Fast, efficient Postgres driver)
* **Auth:** JWT (JSON Web Tokens)
* **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)

### Infrastructure
* **Deployment:** AWS App Runner (Containerized Go Binary)
* **Database Hosting:** AWS RDS (PostgreSQL)

---

## 📂 Project Structure

We follow the **Standard Go Project Layout** to ensure scalability and separation of concerns.

```text
budget-tracker/
├── cmd/
│   └── api/
│       └── main.go           # Entry point: Initializes DB and Router
├── internal/                 # Private application logic
│   ├── models/               # Go structs representing DB tables (User, Transaction)
│   ├── handler/              # HTTP Layer: Parses JSON requests, validation
│   ├── service/              # Business Logic: Budget calculations, rules
│   └── repository/           # Data Layer: Raw SQL queries (pgx)
├── pkg/                      # Public shared utilities
│   ├── database/             # Postgres connection setup
│   └── utils/                # Helper functions (hashing, formatting)
├── migrations/               # SQL files for DB versioning
├── go.mod                    # Dependency manager
└── README.md
