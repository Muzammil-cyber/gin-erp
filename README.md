# Pakistani ERP System - Gin Framework Boilerplate

A production-ready, scalable ERP system boilerplate built with **Gin Web Framework**, following **Domain-Driven Design (DDD)** principles.

## 🚀 Features

### Core Architecture

- **Domain-Driven Design (DDD)** with clear separation of concerns
- **Clean Architecture** with layers: domain, usecase, repository, delivery
- **Dependency Injection** using a container pattern
- **MongoDB** for primary data storage
- **Redis** for caching, rate limiting, and OTP storage

### Security & Authentication

- **JWT-based Authentication** with Token Rotation (Access + Refresh Tokens)
- **Role-Based Access Control (RBAC)** with roles:
  - Admin
  - Customer
  - Finance Manager
  - Manager
- **Password Hashing** using bcrypt
- **Redis-based Rate Limiting** to prevent brute-force attacks

### Registration & Verification

- **Pakistani Phone Number Validation** (+923xxxxxxxxx format)
- **Email OTP Verification** with 5-minute TTL stored in Redis
- **Mock SMTP Service** (ready for production SMTP integration)

### Mobile-Ready API

- **Standardized JSON Response Format**:

  ```json
  {
    "success": true/false,
    "data": {},
    "error": "error message",
    "trace_id": "unique-trace-id"
  }
  ```

- **Request Tracing** with unique trace IDs
- **CORS Support** for cross-origin requests

### Additional Features

- **Graceful Shutdown** for the HTTP server
- **Comprehensive Logging** with request details
- **Environment-based Configuration** using Viper
- **Docker Compose** setup for MongoDB and Redis
- **Health Check Endpoint**
- **Swagger/OpenAPI Documentation** - Interactive API documentation at `/swagger/index.html`

## 📁 Project Structure

```tree
gin-erp/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── container/
│   │   └── container.go            # Dependency injection container
│   ├── domain/
│   │   └── auth/
│   │       ├── entity.go           # Domain entities (User, Role, etc.)
│   │       ├── repository.go       # Repository interfaces
│   │       ├── usecase.go          # Use case interfaces
│   │       └── errors.go           # Domain-specific errors
│   ├── usecase/
│   │   └── auth_usecase.go         # Auth business logic
│   ├── repository/
│   │   ├── mongodb/
│   │   │   ├── user_repository.go
│   │   │   └── refresh_token_repository.go
│   │   └── redis/
│   │       ├── otp_repository.go
│   │       └── rate_limiter_repository.go
│   ├── service/
│   │   └── email_service.go        # Email/OTP service
│   └── delivery/
│       └── http/
│           ├── handler/
│           │   └── auth_handler.go # HTTP handlers
│           ├── middleware/
│           │   ├── auth_middleware.go
│           │   ├── rate_limiter_middleware.go
│           │   ├── cors_middleware.go
│           │   └── logger_middleware.go
│           └── routes/
│               └── routes.go       # Route definitions
├── pkg/
│   ├── database/
│   │   ├── mongodb/
│   │   │   └── mongodb.go          # MongoDB client
│   │   └── redis/
│   │       └── redis.go            # Redis client
│   └── utils/
│       ├── jwt.go                  # JWT utilities
│       ├── password.go             # Password hashing
│       ├── validator.go            # Phone number validation
│       ├── otp.go                  # OTP generation
│       └── response.go             # Standardized responses
├── tests/
│   ├── unit/
│   │   └── auth_usecase_test.go    # Unit tests for business logic
│   ├── http/
│   │   └── routes_test.go          # HTTP route tests
│   └── integration/
│       └── auth_integration_test.go # End-to-end integration tests
├── docs/                           # Swagger documentation (generated)
├── .env.example                    # Environment variables template
├── .gitignore
├── docker-compose.yml              # Docker services
├── Makefile                        # Common commands
├── go.mod
└── README.md

```

## 🛠️ Prerequisites

- **Go 1.25+**
- **MongoDB 6.0+**
- **Redis 7+**
- **Docker & Docker Compose** (optional, for running dependencies)

## 🚦 Getting Started

### 1. Clone the Repository

```bash
git clone repo-link
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Setup Environment Variables

```bash
cp .env.example .env
```

### 4. Start Dependencies (Using Docker)

```bash
make docker-up
```

This will start:

- MongoDB on `localhost:27017`
- Redis on `localhost:6379`
- Mongo Express (MongoDB UI) on `localhost:8081`

### 5. Run the Application

```bash
make run
```

Or with hot reload (requires [air](https://github.com/cosmtrek/air)):

```bash
make install-tools  # Install air
make dev           # Run with hot reload
```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### API Documentation

**Interactive Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

For detailed documentation on using Swagger, see [SWAGGER.md](SWAGGER.md).

### Public Endpoints

| Method | Endpoint | Description |
| -------- | ---------- | ------------- |
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login user |
| POST | `/api/v1/auth/verify-otp` | Verify OTP |
| POST | `/api/v1/auth/refresh-token` | Refresh access token |
| POST | `/api/v1/auth/resend-otp` | Resend OTP |

### Protected Endpoints (Require Authentication)

| Method | Endpoint | Description | Required Role |
| -------- | ---------- | ------------- | --------------- |
| GET | `/api/v1/auth/profile` | Get user profile | Any authenticated |
| GET | `/api/v1/admin/users` | Admin endpoint | Admin |
| GET | `/api/v1/finance/reports` | Finance reports | Admin, Finance Manager |
| GET | `/api/v1/manager/dashboard` | Manager dashboard | Admin, Manager |

## 📝 Example API Requests

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "phone": "+923001234567",
    "password": "SecurePass123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer"
  }'
```

### Verify OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Get Profile (with Bearer Token)

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🧪 Testing

### Run Unit Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates a `coverage.html` file that you can open in your browser.

### Run Integration Tests

Integration tests require MongoDB and Redis to be running:

```bash
make docker-up
go test -v ./tests/integration
```

## 🔧 Configuration

Configuration is managed through environment variables. See [.env.example](.env.example) for all available options.

Key configurations:

- **Server**: `APP_PORT`, `APP_ENV`, `APP_DEBUG`
- **MongoDB**: `MONGODB_URI`, `MONGODB_DATABASE`
- **Redis**: `REDIS_HOST`, `REDIS_PORT`
- **JWT**: `JWT_SECRET`, `JWT_ACCESS_TOKEN_EXPIRE_MINUTES`, `JWT_REFRESH_TOKEN_EXPIRE_HOURS`
- **OTP**: `OTP_LENGTH`, `OTP_EXPIRE_MINUTES`
- **Rate Limiting**: `RATE_LIMIT_REQUESTS_PER_MINUTE`, `RATE_LIMIT_LOGIN_REQUESTS_PER_MINUTE`
- **SMTP**: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`

## 🏗️ Architecture Overview

### Domain-Driven Design (DDD)

The project follows DDD principles with clear separation:

1. **Domain Layer** (`internal/domain/`):
   - Entities (User, RefreshToken, OTP)
   - Repository interfaces
   - Use case interfaces
   - Domain-specific errors

2. **Use Case Layer** (`internal/usecase/`):
   - Business logic implementation
   - Orchestrates repositories and services
   - Independent of delivery mechanism

3. **Repository Layer** (`internal/repository/`):
   - Data access implementation
   - MongoDB repositories
   - Redis repositories

4. **Delivery Layer** (`internal/delivery/`):
   - HTTP handlers
   - Middleware
   - Route definitions
   - Request/Response mapping

### Dependency Injection

The `container.go` initializes all dependencies and wires them together, making the system:

- **Testable**: Easy to mock dependencies
- **Maintainable**: Clear dependency graph
- **Scalable**: Easy to add new features

## 🔐 Security Features

1. **Password Security**: Bcrypt hashing with default cost
2. **JWT Security**: HMAC-SHA256 signing with secret rotation support
3. **Rate Limiting**: Redis-based, prevents brute force attacks
4. **OTP Expiration**: 5-minute TTL on Redis
5. **Token Rotation**: Refresh token mechanism prevents token theft
6. **Phone Validation**: Strict Pakistani phone number format (+923xxxxxxxxx)
7. **Role-Based Access**: Middleware enforces role permissions

## 📦 Dependencies

Major dependencies:

- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **MongoDB Driver**: [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)
- **Redis Client**: [go-redis](https://github.com/redis/go-redis)
- **JWT**: [jwt-go](https://github.com/golang-jwt/jwt)
- **Config**: [Viper](https://github.com/spf13/viper)
- **Testing**: [Testify](https://github.com/stretchr/testify)
- **Password**: [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## 🚀 Production Deployment

### Build

```bash
make build
```

This creates a binary in `./bin/api`.

### Environment Variables

Ensure all production environment variables are set, especially:

- `JWT_SECRET`: Use a strong, random secret
- `MONGODB_URI`: Production MongoDB connection string
- `REDIS_HOST` and `REDIS_PORT`: Production Redis instance
- `SMTP_*`: Configure real SMTP server for emails

### Docker Deployment

You can containerize the application:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./api"]
```

## 🤝 Contributing

This is a boilerplate project. Feel free to:

- Fork and customize for your needs
- Add new features (inventory, sales, finance modules)
- Improve existing code
- Add more tests

## 📄 License

This project is open-source and available under the MIT License.

## 👨‍💻 Author

Created by [Muzammil Loya](https://github.com/muzammil-cyber)

## 🙏 Acknowledgments

- Inspired by clean architecture principles
- Built for the Pakistani market with localized features
- Designed for high-scale ERP systems
