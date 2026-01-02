# Quick Start Guide - Pakistani ERP System

## 🎯 Prerequisites Check

Before starting, ensure you have:

- ✅ Go 1.23+ installed
- ✅ Docker and Docker Compose installed
- ✅ Git installed

## 🚀 Quick Start (5 minutes)

### Step 1: Start Infrastructure Services

```bash
# Start MongoDB and Redis using Docker Compose
make docker-up

# Verify services are running
docker ps
```

You should see:

- MongoDB on port 27017
- Redis on port 6379  
- Mongo Express (web UI) on port 8081

### Step 2: Install Dependencies

```bash
# Download Go dependencies
make deps

# Or manually:
go mod download
go mod tidy
```

### Step 3: Configure Environment

The `.env` file has been created with default values. For local development, the defaults work fine. For production, update these critical values:

```env
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
MONGODB_URI=mongodb://localhost:27017
REDIS_HOST=localhost
```

### Step 4: Run the Application

```bash
# Run the server
make run

# Or with hot reload (requires air):
make install-tools
make dev
```

The server will start on: **<http://localhost:8080>**

### Step 5: Access Swagger UI (Optional)

Open your browser and navigate to:

**<http://localhost:8080/swagger/index.html>**

This provides an interactive interface to test all API endpoints. See [SWAGGER.md](SWAGGER.md) for details.

## 📡 Test the API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "success": true,
  "data": {
    "status": "ok",
    "service": "Pakistani ERP System"
  },
  "trace_id": "unique-id-here"
}
```

### 2. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "phone": "+923001234567",
    "password": "SecurePassword123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer"
  }'
```

**Important:** Check your terminal output for the OTP code (it's logged in development mode).

### 3. Verify Email with OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "code": "123456"
  }'
```

Replace `123456` with the OTP from your terminal logs.

### 4. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePassword123"
  }'
```

Save the `access_token` from the response!

### 5. Access Protected Route

```bash
# Replace YOUR_ACCESS_TOKEN with the token from login
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🧪 Run Tests

```bash
# Run unit tests
make test

# Run with coverage
make test-coverage

# Integration tests (requires MongoDB and Redis running)
go test -v ./tests/integration
```

## 🗄️ MongoDB Web UI

Access Mongo Express at: **<http://localhost:8081>**

Credentials:

- Username: `admin`
- Password: `admin123`

Here you can:

- View the `users` collection
- View the `refresh_tokens` collection
- Manually verify data

## 🔍 Understanding the Response Format

All API responses follow this standard format:

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "user": {
      "id": "...",
      "email": "john.doe@example.com",
      "role": "customer"
    }
  },
  "trace_id": "unique-request-id"
}
```

For errors:

```json
{
  "success": false,
  "error": "invalid email or password",
  "trace_id": "unique-request-id"
}
```

## 🛠️ Common Commands

```bash
# Start infrastructure
make docker-up

# Stop infrastructure
make docker-down

# View logs
make docker-logs

# Build binary
make build

# Run application
make run

# Run with hot reload
make dev

# Run tests
make test

# Generate coverage report
make test-coverage

# Clean build artifacts
make clean

# See all available commands
make help
```

## 🎨 Available User Roles

The system supports 4 roles with different permissions:

1. **admin** - Full system access
2. **customer** - Basic user access
3. **finance_manager** - Access to financial operations
4. **manager** - Access to management operations

## 📊 Project Structure Overview

```tree
gin-erp/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── container/               # Dependency injection
│   ├── domain/auth/             # Business entities & interfaces
│   ├── usecase/                 # Business logic
│   ├── repository/              # Data access layer
│   ├── service/                 # External services (email)
│   └── delivery/http/           # HTTP handlers & middleware
├── pkg/
│   ├── database/                # DB clients (MongoDB, Redis)
│   └── utils/                   # Utilities (JWT, password, etc.)
└── tests/integration/           # Integration tests
```

## 🐛 Troubleshooting

### MongoDB Connection Error

```bash
# Check if MongoDB is running
docker ps | grep mongo

# Restart MongoDB
docker-compose restart mongodb
```

### Redis Connection Error

```bash
# Check if Redis is running
docker ps | grep redis

# Restart Redis
docker-compose restart redis
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Build Errors

```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

## 📚 Next Steps

1. **Explore the Code**: Start with [cmd/api/main.go](cmd/api/main.go)
2. **Add New Features**: Create new domains in `internal/domain/`
3. **Write Tests**: Add tests in corresponding `_test.go` files
4. **Configure Production**: Update `.env` for production deployment
5. **Add More Endpoints**: Extend routes in `internal/delivery/http/routes/`

## 🔐 Security Notes for Production

Before deploying to production:

1. ✅ Change `JWT_SECRET` to a strong random value
2. ✅ Use a production MongoDB instance (MongoDB Atlas recommended)
3. ✅ Use a production Redis instance (Redis Cloud or AWS ElastiCache)
4. ✅ Configure real SMTP for email sending
5. ✅ Set `APP_ENV=production`
6. ✅ Set `APP_DEBUG=false`
7. ✅ Use HTTPS with proper SSL certificates
8. ✅ Implement rate limiting at API gateway level
9. ✅ Set up proper monitoring and logging
10. ✅ Regular security audits and dependency updates

## 📞 Support

For issues or questions:

- Check the [README.md](README.md) for detailed documentation
- Review the code comments
- Check integration tests for usage examples
