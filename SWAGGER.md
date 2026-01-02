# Swagger API Documentation

## Overview

This project uses Swagger/OpenAPI for API documentation. The Swagger UI provides an interactive interface to explore and test all API endpoints.

## Accessing Swagger UI

Once the server is running, access Swagger UI at:

```http
http://localhost:8080/swagger/index.html
```

## Generating Swagger Documentation

Swagger documentation is generated from code annotations using the `swag` CLI tool.

### Manual Generation

```bash
make swagger
```

Or directly:

```bash
swag init -g cmd/api/main.go
```

This will generate:

- `docs/docs.go` - Go code for embedding docs
- `docs/swagger.json` - OpenAPI JSON specification
- `docs/swagger.yaml` - OpenAPI YAML specification

## Adding Swagger Annotations

### Package Level (main.go)

```go
// @title Pakistani ERP System API
// @version 1.0
// @description Production-ready ERP system with JWT authentication
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Handler Functions

```go
// @Summary Register a new user
// @Description Create a new user account with email verification
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "User registration data"
// @Success 201 {object} utils.Response{data=domain.User}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // ...
}
```

## Swagger Annotation Tags

| Tag | Description | Example |
| ----- | ------------- | --------- |
| `@Summary` | Short description | `@Summary Login user` |
| `@Description` | Detailed description | `@Description Authenticate user with email/phone and password` |
| `@Tags` | Group endpoints | `@Tags Auth` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path string true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} utils.Response` |
| `@Failure` | Error response | `@Failure 400 {object} utils.Response` |
| `@Router` | Endpoint path | `@Router /auth/login [post]` |
| `@Security` | Auth requirement | `@Security BearerAuth` |

## Testing with Swagger UI

1. **Start the server**:

   ```bash
   make run
   # or for hot reload:
   make dev
   ```

2. **Open Swagger UI** in your browser:

   ```http
   http://localhost:8080/swagger/index.html
   ```

3. **Test authentication flow**:
   - Use `/auth/register` to create an account
   - Use `/auth/verify-otp` to verify your account (check console for OTP)
   - Use `/auth/login` to get JWT token
   - Click "Authorize" button at top of Swagger UI
   - Enter token in format: `Bearer <your-token>`
   - Test protected endpoints like `/auth/profile`

## Parameter Types

### Path Parameters

```go
// @Param id path string true "User ID"
```

### Query Parameters

```go
// @Param page query int false "Page number" default(1)
```

### Body Parameters

```go
// @Param request body domain.LoginRequest true "Login credentials"
```

### Header Parameters

```go
// @Param Authorization header string true "Bearer token"
```

## Response Objects

Define response structures in your domain models:

```go
type User struct {
    ID          string    `json:"id" example:"507f1f77bcf86cd799439011"`
    Email       string    `json:"email" example:"user@example.com"`
    PhoneNumber string    `json:"phone_number" example:"+923001234567"`
    Role        Role      `json:"role" example:"customer"`
    IsVerified  bool      `json:"is_verified" example:"true"`
    CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}
```

## Workflow

1. **Write handler function** with business logic
2. **Add Swagger annotations** above the handler
3. **Run `make swagger`** to regenerate docs
4. **Restart server** to load new documentation
5. **Test in Swagger UI** at `/swagger/index.html`

## Best Practices

1. **Keep annotations up to date**: Regenerate docs after API changes
2. **Use examples**: Add `example` tags to struct fields for better documentation
3. **Document all responses**: Include success and all error cases
4. **Group endpoints**: Use consistent `@Tags` for logical grouping
5. **Add descriptions**: Provide clear descriptions for complex operations
6. **Security annotations**: Add `@Security BearerAuth` for protected routes

## Common Issues

### Swagger UI not loading

- Ensure `docs` directory exists with generated files
- Check that `_ "github.com/muzammil-cyber/gin-erp/docs"` is imported in main.go
- Verify Swagger route is registered: `/swagger/*any`

### Documentation not updating

- Run `make swagger` to regenerate docs
- Restart the server to reload generated files
- Clear browser cache

### Missing endpoints

- Check that handler has proper Swagger annotations
- Verify `@Router` path matches actual route
- Ensure HTTP method in annotation matches route

## References

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [Gin-Swagger](https://github.com/swaggo/gin-swagger)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
