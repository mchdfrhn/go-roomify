# go-roomify

A REST API for booking meeting rooms inside an organisation. Employees reserve rooms, General Affairs staff approve or decline, and admins manage rooms, facilities, users, and divisions. Built in Go with Gin and PostgreSQL, following a clean architecture layout.

Originally written during the Enigma Camp backend bootcamp (2024), later documented as a step-by-step guide to growing a Go service from a single file into clean architecture (see [`docs/`](docs) and [`GO_CLEAN_ARCHITECTURE_GUIDE.md`](GO_CLEAN_ARCHITECTURE_GUIDE.md), written in Indonesian).

## Features

- **JWT authentication** with role-based access: `admin`, `ga` (General Affairs), and `employee`
- **Reservations** with a status flow: pending → accepted / declined, or cancelled by the requester
- **Rooms, room types, and facilities** management, including a room availability endpoint
- **Users, roles, divisions, and credentials** management, with passwords hashed using bcrypt
- **Yearly reservation report** exported as CSV
- Paginated list endpoints and request validation

## Tech stack

| Layer | Tools |
| --- | --- |
| HTTP | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL via `database/sql` + `lib/pq` |
| Auth | JWT, bcrypt |
| Config | `.env` via godotenv |
| Packaging | Docker (multi-stage build) |

## Project structure

```text
config/        environment config and database connection
delivery/      Gin server bootstrap and HTTP controllers
middleware/    JWT + role guard
model/         domain models and request/response DTOs
repository/    SQL data access
usecase/       business logic
utils/         JWT helpers, hashing, pagination, query builders, validation
docs/          step-by-step build guide (Indonesian)
```

Requests flow `controller → usecase → repository`, and every layer depends on interfaces, wired together in `delivery/server.go`.

## API overview

All routes are served under `/api` on port `8085`.

| Resource | Base path | Notes |
| --- | --- | --- |
| Auth | `/api/auth` | `POST /login` returns a JWT; user credential management is admin-only |
| Rooms | `/api/room` | create/update/delete: admin · status: admin, GA · read and `/available`: all roles |
| Room types | `/api/type/room` | CRUD |
| Facilities | `/api/facility` | read: public · write: admin |
| Reservations | `/api/reservation` | create, list your own, get by id, `PUT /status` |
| Users | `/api/user` | CRUD with pagination |
| Roles | `/api/roles` | CRUD |
| Divisions | `/api/divisions` | CRUD |
| Report | `/api/report?s=<startYear>&e=<endYear>` | CSV download |

Send the token as `Authorization: Bearer <token>`.

## Running locally

Requirements: Go 1.25+ and a PostgreSQL database.

```bash
cp .env.example .env   # fill in the values below
go mod download
go run .
```

| Variable | Example | Purpose |
| --- | --- | --- |
| `DB_HOST` / `DB_PORT` | `localhost` / `5432` | PostgreSQL address |
| `DB_NAME` / `DB_USER` / `DB_PASS` | | database credentials |
| `DB_DRIVER` | `postgres` | `database/sql` driver name |
| `TOKEN_LIFE_TIME` | `24` | JWT lifetime in hours |
| `ISSUER_NAME` | `go-roomify` | JWT issuer |
| `SIGNATURE` | a long random string | JWT signing key |

### With Docker

The app reads its config from a `.env` file in its working directory, so mount yours into the container:

```bash
docker build -t go-roomify .
docker run -v "$(pwd)/.env:/app/.env" -p 8085:8085 go-roomify
```

The database schema is not included in this repository yet, so you need to create the tables before the API can serve data.

## Author

[Mochammad Farhan Ali](https://www.mochamadfarhanali.my.id) · [GitHub](https://github.com/mchdfrhn)
