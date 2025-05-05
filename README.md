# Dormitory Management System

Backend API for Dormitory Management

## Requirements

- Go 1.16+
- PostgreSQL 12+
- Redis 6+

## Environment Configuration

Create a `.env` file in the project root with the following content:

```
# Database
DB_CONN_STR=postgresql://username:password@localhost:5432/dormitory_db?sslmode=disable

# JWT
JWT_ACCESS_SECRET=your_access_token_secret
JWT_REFRESH_SECRET=your_refresh_token_secret

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/dormitory-management.git
cd dormitory-management/server
```

2. Install dependencies:

```bash
go mod download
```

3. Create PostgreSQL database:

```sql
CREATE DATABASE dormitory_db;
```

## Running the Project

1. Build and run the server:

```bash
go run cmd/main.go
```

Or build a binary:

```bash
go build -o dormitory-server cmd/main.go
./dormitory-server
```

The server will run on port 8080 by default.

## Project Structure

```
server/
├── cmd/                # Entry point
│   └── main.go
├── internal/
│   ├── database/       # Database and Redis connections
│   ├── handlers/       # API handlers
│   ├── middleware/     # Middleware (auth, logging)
│   ├── models/         # Models and structs
│   ├── server/         # Server setup and routes
│   ├── types/          # Type definitions
│   └── utils/          # Utilities (JWT, response)
└── pkg/                # Shared packages
```

## API Endpoints

### Authentication

- **POST /api/auth/register** - Register a new account
- **POST /api/auth/login** - Login
- **GET /api/auth/me** - Get current user information

For other API details, refer to the source code or API documentation.
