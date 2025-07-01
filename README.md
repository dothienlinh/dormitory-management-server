# Dormitory Management System

Backend API for Dormitory Management

## Requirements

- Go 1.16+
- PostgreSQL 12+
- Redis 6+

## Environment Configuration

Create a `.env` file in the project root with the following content:

```
# Server configuration
SERVER_PORT=8080
SERVER_MODE=debug

# Database configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dormitory
DB_SSL_MODE=disable

# Redis configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis
REDIS_DB=0

# JWT configuration
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRES_IN=3600
JWT_REFRESH_EXPIRES_IN=604800

# Log level
LOG_LEVEL=info
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
CREATE DATABASE dormitory;
```

## Running the Project

### Using Docker (Recommended)

```bash
docker-compose up -d
```

### Running Locally

1. Run migration

```bash
make migrate-up
```

2. Build and run the server:

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
├── internal/           # Internal packages
│   ├── config/         # Configuration structures
│   ├── database/       # Database connections
│   ├── delivery/       # HTTP delivery layer
│   │   ├── http/       # HTTP handlers and routes
│   │   │   ├── handler/    # HTTP request handlers
│   │   │   ├── middleware/ # Middleware (auth, logging)
│   │   │   └── router/     # Route configuration
│   ├── domain/         # Domain layer (entities, interfaces)
│   │   ├── entity/     # Domain entities
│   │   ├── repository/ # Repository interfaces
│   │   └── usecase/    # Use case interfaces
│   ├── repository/     # Repository implementations
│   └── usecase/        # Use case implementations
└── pkg/                # Shared packages
    └── logger/         # Logging utilities
```

## Architecture

This project follows Clean Architecture principles with the following layers:

1. **Domain Layer**: Contains business entities and interfaces (repository and use case interfaces)
2. **Repository Layer**: Implements data access logic
3. **Use Case Layer**: Implements business logic
4. **Delivery Layer**: Handles HTTP requests/responses
   - **Handler**: Contains the logic for handling HTTP requests
   - **Router**: Configures and organizes API routes and endpoints
     - Routes are organized by entity type (auth, user, room, etc.)
   - **Middleware**: Processes requests before they reach handlers (auth, logging, etc.)

The architecture follows the dependency rule where inner layers (Domain) are independent of outer layers. Dependencies point inward, with adapters connecting the different layers.

## API Endpoints

All API routes are versioned under `/api/v1` prefix.

### Authentication

- **POST /api/v1/auth/register** - Register a new account
- **POST /api/v1/auth/login** - Login to get access and refresh tokens
- **POST /api/v1/auth/refresh** - Refresh access token
- **POST /api/v1/auth/logout** - Logout (requires authentication)

### Users

- **GET /api/v1/users** - Get list of users (requires authentication)
- **GET /api/v1/users/:id** - Get user by ID (requires authentication)
- **PUT /api/v1/users/:id** - Update user (requires authentication)
- **DELETE /api/v1/users/:id** - Delete user (requires admin role)
- **POST /api/v1/users/room** - Assign user to room (requires authentication)
- **POST /api/v1/users/remove-room** - Remove user from room (requires authentication)

### Rooms

- **GET /api/v1/rooms** - Get list of rooms
- **GET /api/v1/rooms/:id** - Get room by ID
- **POST /api/v1/rooms** - Create room (requires admin role)
- **PUT /api/v1/rooms/:id** - Update room (requires admin role)
- **DELETE /api/v1/rooms/:id** - Delete room (requires admin role)

### Room Categories

- **GET /api/v1/room-categories** - Get list of room categories
- **GET /api/v1/room-categories/:id** - Get room category by ID
- **POST /api/v1/room-categories** - Create room category (requires admin role)
- **PUT /api/v1/room-categories/:id** - Update room category (requires admin role)
- **DELETE /api/v1/room-categories/:id** - Delete room category (requires admin role)

### Contracts

- **GET /api/v1/contracts** - Get list of contracts (requires authentication)
- **GET /api/v1/contracts/:id** - Get contract by ID (requires authentication)
- **GET /api/v1/contracts/user/:user_id** - Get contracts by user ID (requires authentication)
- **POST /api/v1/contracts** - Create contract (requires admin role)
- **PUT /api/v1/contracts/:id** - Update contract (requires admin role)
- **DELETE /api/v1/contracts/:id** - Delete contract (requires admin role)

## Request/Response Formats

### User Registration

**Request**:

```json
POST /api/v1/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "1234567890"
}
```

**Response**:

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2023-05-15T10:30:45Z"
  }
}
```

### User Login

**Request**:

```json
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response**:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

### Remove User from Room

**Request**:

```json
POST /api/v1/users/remove-room
{
  "user_id": 1,
  "room_id": 5
}
```

**Response**:

```json
{
  "success": true,
  "message": "User removed from room successfully"
}
```

## Error Handling

All API errors follow a consistent format:

```json
{
  "success": false,
  "message": "Error message",
  "error": "Error type"
}
```

Common HTTP status codes:

- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## Development

For local development, you can use the provided Docker Compose file:

```bash
docker-compose up -d
```

This will start PostgreSQL and Redis services, along with the API server.

## Production Deployment

For production deployment, you can use the provided Docker Compose file with a production environment file:

```bash
docker-compose --env-file .env.production up -d
```
