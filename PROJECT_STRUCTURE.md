# Cấu Trúc Dự Án Dormitory Management

Dự án này được xây dựng theo nguyên tắc Clean Architecture, tạo ra một hệ thống có khả năng mở rộng, dễ bảo trì và kiểm thử.

## Clean Architecture

Clean Architecture được thực hiện qua việc phân tách mã nguồn thành các lớp (layers) với nguyên tắc phụ thuộc hướng vào trong. Các lớp bên ngoài có thể phụ thuộc vào các lớp bên trong, nhưng không ngược lại.

```
┌─────────────────────────────────────────────────────────┐
│                                                         │
│  ┌─────────────┐      ┌─────────────┐    ┌──────────┐   │
│  │             │      │             │    │          │   │
│  │  Delivery   │◄─────┤   UseCase   │◄───┤  Domain  │   │
│  │             │      │             │    │          │   │
│  └─────────────┘      └─────────────┘    └──────────┘   │
│         ▲                    ▲                 ▲        │
│         │                    │                 │        │
│         └─────────────┐      │                 │        │
│                       │      │                 │        │
│                  ┌────┴──────┴────┐            │        │
│                  │                │            │        │
│                  │   Repository   │◄───────────┘        │
│                  │                │                     │
│                  └────────────────┘                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Cấu Trúc Thư Mục

```
server/
├── cmd/                            # Entry points
│   └── main.go                     # Main application entry point
│
├── internal/                       # Private application code
│   ├── config/                     # Configuration structures
│   │   └── config.go               # Configuration object
│   │
│   ├── database/                   # Database connections
│   │   ├── postgres.go             # PostgreSQL connection
│   │   └── redis.go                # Redis connection
│   │
│   ├── delivery/                   # Delivery Layer
│   │   ├── http/                   # HTTP Handlers
│   │   │   ├── middleware/         # HTTP middleware
│   │   │   │   ├── auth.go         # Authentication middleware
│   │   │   │   ├── error.go        # Error handling middleware
│   │   │   │   ├── logger.go       # Logger middleware
│   │   │   │   └── middleware.go   # Middleware factory
│   │   │   │
│   │   │   ├── auth_handler.go     # Authentication handlers
│   │   │   ├── contract_handler.go # Contract handlers
│   │   │   ├── handlers.go         # Handlers factory
│   │   │   ├── room_handler.go     # Room handlers
│   │   │   ├── room_category_handler.go # Room category handlers
│   │   │   ├── server.go           # HTTP server setup and routes
│   │   │   └── user_handler.go     # User handlers
│   │
│   ├── domain/                     # Domain Layer (Enterprise Business Rules)
│   │   ├── entity/                 # Business objects
│   │   │   ├── contract.go         # Contract entity
│   │   │   ├── room.go             # Room entity
│   │   │   ├── room_category.go    # Room category entity
│   │   │   ├── room_rent.go        # Room rent entity
│   │   │   └── user.go             # User entity
│   │   │
│   │   ├── repository/             # Repository interfaces
│   │   │   ├── contract_repository.go  # Contract repository interface
│   │   │   ├── repositories.go     # Repository factory interface
│   │   │   ├── room_category_repository.go # Room category repository interface
│   │   │   ├── room_repository.go  # Room repository interface
│   │   │   └── user_repository.go  # User repository interface
│   │   │
│   │   └── usecase/                # Use case interfaces
│   │       ├── auth_usecase.go     # Auth use case interface
│   │       ├── contract_usecase.go # Contract use case interface
│   │       ├── room_category_usecase.go # Room category use case interface
│   │       ├── room_usecase.go     # Room use case interface
│   │       ├── usecases.go         # Use case factory interface
│   │       └── user_usecase.go     # User use case interface
│   │
│   ├── repository/                 # Repository Layer (Data Access)
│   │   ├── contract_repository.go  # Contract repository implementation
│   │   ├── repositories.go         # Repository factory implementation
│   │   ├── room_category_repository.go # Room category repository implementation
│   │   ├── room_repository.go      # Room repository implementation
│   │   └── user_repository.go      # User repository implementation
│   │
│   └── usecase/                    # Use Case Layer (Application Business Rules)
│       ├── auth_usecase.go         # Auth use case implementation
│       ├── contract_usecase.go     # Contract use case implementation
│       ├── room_category_usecase.go # Room category use case implementation
│       ├── room_usecase.go         # Room use case implementation
│       ├── usecases.go             # Use case factory implementation
│       └── user_usecase.go         # User use case implementation
│
└── pkg/                            # Public libraries
    └── logger/                     # Logging package
        └── logger.go               # Logger implementation
```

## Chi Tiết Các Layer

### 1. Domain Layer (internal/domain/)

Domain layer là layer trung tâm và không phụ thuộc vào bất kỳ layer nào khác. Layer này định nghĩa:

- **Entities (internal/domain/entity/)**: Cấu trúc dữ liệu chứa business rules của hệ thống:

  - `user.go`: Đối tượng người dùng (admin, manager, student)
  - `room.go`: Đối tượng phòng ở
  - `room_category.go`: Loại phòng và giá cả
  - `contract.go`: Hợp đồng thuê phòng
  - `room_rent.go`: Thông tin về việc thuê phòng

- **Repository Interfaces (internal/domain/repository/)**: Định nghĩa các interface để thao tác với dữ liệu:

  - `user_repository.go`: Interface cho thao tác với dữ liệu người dùng
  - `room_repository.go`: Interface cho thao tác với dữ liệu phòng
  - `room_category_repository.go`: Interface cho thao tác với dữ liệu loại phòng
  - `contract_repository.go`: Interface cho thao tác với dữ liệu hợp đồng
  - `repositories.go`: Factory cho các repository

- **Use Case Interfaces (internal/domain/usecase/)**: Định nghĩa business logic interface:
  - `auth_usecase.go`: Interface cho logic xác thực
  - `user_usecase.go`: Interface cho logic quản lý người dùng
  - `room_usecase.go`: Interface cho logic quản lý phòng
  - `room_category_usecase.go`: Interface cho logic quản lý loại phòng
  - `contract_usecase.go`: Interface cho logic quản lý hợp đồng
  - `usecases.go`: Factory cho các use case

### 2. Repository Layer (internal/repository/)

Repository layer thực hiện các interface được định nghĩa trong domain/repository. Layer này phụ thuộc vào domain layer và chịu trách nhiệm cho tương tác với cơ sở dữ liệu (PostgreSQL):

- `user_repository.go`: Thực hiện các thao tác CRUD cho User
- `room_repository.go`: Thực hiện các thao tác CRUD cho Room
- `room_category_repository.go`: Thực hiện các thao tác CRUD cho RoomCategory
- `contract_repository.go`: Thực hiện các thao tác CRUD cho Contract
- `repositories.go`: Factory để tạo tất cả các repository instances

### 3. Use Case Layer (internal/usecase/)

Use case layer thực hiện các interface được định nghĩa trong domain/usecase. Layer này chứa business logic của ứng dụng và phụ thuộc vào domain layer:

- `auth_usecase.go`: Xử lý đăng ký, đăng nhập, refresh token, logout
- `user_usecase.go`: Xử lý các tác vụ liên quan đến người dùng
- `room_usecase.go`: Xử lý các tác vụ liên quan đến phòng
- `room_category_usecase.go`: Xử lý các tác vụ liên quan đến loại phòng
- `contract_usecase.go`: Xử lý các tác vụ liên quan đến hợp đồng
- `usecases.go`: Factory để tạo tất cả các use case instances

### 4. Delivery Layer (internal/delivery/)

Delivery layer xử lý việc giao tiếp với client thông qua HTTP API. Layer này phụ thuộc vào use case layer:

- **HTTP Handlers (internal/delivery/http/)**: Xử lý HTTP requests/responses:

  - `auth_handler.go`: Xử lý authentication requests
  - `user_handler.go`: Xử lý user requests
  - `room_handler.go`: Xử lý room requests
  - `room_category_handler.go`: Xử lý room category requests
  - `contract_handler.go`: Xử lý contract requests
  - `handlers.go`: Factory tạo tất cả handler instances
  - `server.go`: Khởi tạo HTTP server và cấu hình routes

- **Middleware (internal/delivery/http/middleware/)**: Cung cấp middleware cho HTTP server:
  - `auth.go`: Xác thực JWT
  - `error.go`: Xử lý lỗi
  - `logger.go`: Ghi log request/response
  - `middleware.go`: Factory cho middleware

### 5. Configuration (internal/config/)

Configuration layer chứa cấu trúc cấu hình và logic để đọc cấu hình từ biến môi trường:

- `config.go`: Định nghĩa cấu trúc cấu hình và đọc cấu hình từ .env hoặc biến môi trường

### 6. Database (internal/database/)

Database layer xử lý kết nối và tương tác với databases:

- `postgres.go`: Khởi tạo và quản lý kết nối PostgreSQL
- `redis.go`: Khởi tạo và quản lý kết nối Redis

### 7. Shared Packages (pkg/)

Shared packages chứa code có thể sử dụng bởi nhiều projects:

- **Logger (pkg/logger/)**: Package xử lý logging:
  - `logger.go`: Interface và implementation của logger

## Luồng Dữ Liệu

1. HTTP Request đến Gin router
2. Middleware xử lý (authentication, logging, error handling)
3. Handler nhận request và gọi use case method tương ứng
4. Use case thực thi business logic và gọi repository khi cần
5. Repository thực hiện thao tác dữ liệu với database
6. Dữ liệu được trả về qua use case đến handler
7. Handler chuyển đổi dữ liệu thành HTTP response phù hợp

## Dependency Injection

Dự án sử dụng Dependency Injection để kết nối các layer:

1. `main.go` khởi tạo configuration và database connections
2. Repository instances được tạo với database connections
3. Use case instances được tạo với repository instances
4. Handler instances được tạo với use case instances
5. HTTP server được cấu hình với handlers và middleware

## Ví Dụ Luồng Request

Ví dụ luồng xử lý cho request "Get User by ID":

1. HTTP GET request đến `/api/v1/users/:id`
2. AuthMiddleware xác thực JWT token
3. UserHandler.GetUserByID xử lý request
4. UserHandler gọi UserUseCase.GetUserByID
5. UserUseCase gọi UserRepository.GetByID
6. UserRepository thực hiện query PostgreSQL
7. Dữ liệu user được trả về qua Repository -> UseCase -> Handler
8. Handler tạo HTTP response JSON với dữ liệu user

## Dịch Vụ Bổ Sung

- **Redis**: Được sử dụng để lưu trữ JWT refresh tokens và cache
- **JWT**: Được sử dụng cho authentication và authorization
- **Logger**: Logger sử dụng Zap để ghi log hiệu suất cao

## Tài Liệu API

Tất cả API endpoints được phiên bản hóa dưới tiền tố `/api/v1`. Chi tiết về các endpoints có thể tìm thấy trong file README.md.
