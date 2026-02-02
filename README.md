# MyBankSoal API

This is a Question Bank Backend API built with Go, Echo, GORM, and Casbin.

## Features

- **RESTful API**: Built with Echo framework.
- **Database**: Supports SQLite (default) and PostgreSQL using GORM.
- **Authentication**: JWT-based authentication.
- **RBAC**: Role-Based Access Control using Casbin.
- **Swagger**: API Documentation (OpenAPI 3.0).
- **Clean Architecture**: Structured into Handler, UseCase, Repository, and Entity layers.

## Setup

1.  **Clone the repository**
2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```
3.  **Generate Documentation**:
    ```bash
    ./gen-docs.bat
    ```
4.  **Run the application**:
    ```bash
    go run cmd/api/main.go
    ```
    The server will start on port `8080`.

## Database Migrations

This project uses `golang-migrate` for database versioning.

### Run Migrations
To apply all up migrations:
```bash
go run cmd/migrate/main.go up
```

To rollback the last migration:
```bash
go run cmd/migrate/main.go down
```

To migrate a specific number of steps:
```bash
go run cmd/migrate/main.go step -step 1
```

### Create New Migration
To create a new migration, create two files in `db/migrations`:
- `XXXXXX_name.up.sql`
- `XXXXXX_name.down.sql`
Where `XXXXXX` is a sequential number (e.g., `000003`).

## Configuration

Configuration is managed by `.env` file.
Default configuration uses MySQL.

## API Documentation

The API definition is available in standard OpenAPI 3.0 format at `api.yml` in the project root.
Once the server is running, you can access the Swagger UI at:
http://localhost:8080/docs/index.html

## Usage

### 1. Register a User
```bash
POST /auth/register
{
    "username": "admin",
    "password": "password",
    "role": "admin"
}
```

### 2. Login
```bash
POST /auth/login
{
    "username": "admin",
    "password": "password"
}
```
Response:
```json
{
    "success": true,
    "data": {
        "token": "eyJhbG..."
    }
}
```

### 3. Create a Question (Requires Auth)
```bash
POST /questions
Header: Authorization: Bearer <token>
{
    "title": "Sample Question",
    "content": "What is 2+2?",
    "answer": "4"
}
```

### 4. Update Status (Requires Auth)
```bash
PATCH /questions/1/status
Header: Authorization: Bearer <token>
{
    "status": "publish"
}
```

## Folder Structure

- `cmd/api`: Entry point.
- `config`: Configuration files and RBAC model.
- `internal/entity`: Database models.
- `internal/repository`: Data access layer.
- `internal/usecase`: Business logic.
- `internal/handler`: HTTP handlers.
- `internal/middleware`: Auth and RBAC middleware.
- `pkg`: Shared packages (Database, Response, Utils).
