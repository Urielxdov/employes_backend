# CRUD API Go - Employees Database

**Date:** 2026-05-08  
**Scope:** REST API for mysql-employees database with full CRUD operations  
**Tech Stack:** Go, MySQL, Docker

---

## Overview

REST API server in Go exposing CRUD endpoints for mysql-employees database. No authentication. JSON responses. Stdout logging. Health check endpoint. Single Docker Compose setup with MySQL + API.

---

## Architecture

### Server
- **Framework:** Go `net/http` stdlib
- **Database Driver:** `database/sql` + MySQL driver
- **Port:** 8080
- **Response Format:** JSON with `{success, data, error}` envelope

### Database
- **Engine:** MySQL 8.0
- **Source:** mysql-employees schema (employees, departments, salaries, titles, dept_emp, dept_manager)
- **Connection:** Via Docker Compose network
- **Init:** SQL dump loaded on container startup

### Deployment
- **Docker:** Two-service compose (MySQL + API)
- **Healthcheck:** `/health` endpoint, 10s interval
- **Logging:** Stdout only (structured JSON)

---

## API Endpoints

### Core CRUD (all entities)

| Method | Path | Action |
|--------|------|--------|
| GET | `/api/v1/{entity}` | List all (paginated if needed) |
| POST | `/api/v1/{entity}` | Create |
| GET | `/api/v1/{entity}/:id` | Get by ID |
| PUT | `/api/v1/{entity}/:id` | Update |
| DELETE | `/api/v1/{entity}/:id` | Delete |

**Entities:** `employees`, `departments`, `salaries`, `titles`, `dept_emp`, `dept_manager`

### Health
- **GET** `/health` → `{status: "ok"}` (for Docker healthcheck)

---

## Data Models

Map database tables to Go structs:

```
Employee
├── emp_no (int)
├── birth_date (date)
├── first_name (string)
├── last_name (string)
├── gender (string)
├── hire_date (date)

Department
├── dept_no (string)
├── dept_name (string)

Salary
├── emp_no (int)
├── salary (int)
├── from_date (date)
├── to_date (date)

Title
├── emp_no (int)
├── title (string)
├── from_date (date)
├── to_date (date)

DeptEmp
├── emp_no (int)
├── dept_no (string)
├── from_date (date)
├── to_date (date)

DeptManager
├── emp_no (int)
├── dept_no (string)
├── from_date (date)
├── to_date (date)
```

---

## Project Structure

```
.
├── main.go                    # Entry point, server startup
├── models.go                  # Data structures
├── handlers.go                # HTTP handlers for all CRUD
├── db.go                      # Connection pool, queries
├── errors.go                  # Error types + JSON response envelope
├── go.mod
├── go.sum
├── Dockerfile                 # API container
├── docker-compose.yml         # MySQL + API orchestration
├── .env.example               # DB_HOST, DB_PORT, DB_USER, DB_PASS
├── .dockerignore
└── docs/
    └── superpowers/specs/
        └── 2026-05-08-crud-go-employees-design.md
```

---

## Response Format

All endpoints return JSON:

**Success (2xx):**
```json
{
  "success": true,
  "data": {...}
}
```

**Error (4xx/5xx):**
```json
{
  "success": false,
  "error": "specific error message"
}
```

---

## Error Handling

- **400 Bad Request:** Invalid JSON, missing required fields
- **404 Not Found:** Resource doesn't exist
- **500 Internal Server Error:** Database connection, query failure
- All errors logged to stdout with timestamp

---

## Database Initialization

1. MySQL container starts with `mysql-employees` schema loaded
2. API waits for MySQL readiness (retry loop)
3. Runs schema validation on startup
4. Ready to accept requests

---

## Docker Setup

**docker-compose.yml:**
- **mysql** service: Image `mysql:8.0`, volume mount for employees dump
- **api** service: Builds Dockerfile, depends on mysql, exposes 8080

**Healthcheck:**
- `/health` endpoint checked every 10s
- API logs startup timestamp to stdout

**Environment:**
- `.env` file (git-ignored) with DB credentials
- `.env.example` provided as template

---

## Logging

- **Format:** Timestamp + level + message to stdout
- **Events:** Server start, incoming requests, DB errors, shutdown
- **No file rotation needed** (container-managed)

---

## Testing Scope

Not part of MVP. Can be added post-launch if needed (integration tests via docker-compose).

---

## Success Criteria

✓ All 6 entities have GET/POST/PUT/DELETE working  
✓ JSON responses match envelope format  
✓ Healthcheck endpoint responds  
✓ Docker Compose starts cleanly  
✓ Logs appear in `docker-compose logs api`  
✓ No unhandled panics
