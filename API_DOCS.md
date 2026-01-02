# API Documentation - Pakistani ERP System

## Base URL

```
http://localhost:8080/api/v1
```

## Response Format

All API responses follow this standard format:

### Success Response

```json
{
  "success": true,
  "data": { },
  "trace_id": "uuid-here"
}
```

### Error Response

```json
{
  "success": false,
  "error": "error message",
  "trace_id": "uuid-here"
}
```

## Authentication

Most endpoints require authentication using JWT Bearer tokens.

**Header:**

```
Authorization: Bearer <access_token>
```

## Endpoints

### 1. Health Check

Check if the API is running.

**Endpoint:** `GET /health`

**Auth Required:** No

**Response:**

```json
{
  "success": true,
  "data": {
    "status": "ok",
    "service": "Pakistani ERP System"
  },
  "trace_id": "..."
}
```

---

### 2. Register User

Register a new user account.

**Endpoint:** `POST /api/v1/auth/register`

**Auth Required:** No

**Request Body:**

```json
{
  "email": "user@example.com",
  "phone": "+923001234567",
  "password": "Password123",
  "first_name": "John",
  "last_name": "Doe",
  "role": "customer"
}
```

**Fields:**

- `email` (string, required): Valid email address
- `phone` (string, required): Pakistani phone number (+923xxxxxxxxx)
- `password` (string, required): Minimum 8 characters
- `first_name` (string, required): User's first name
- `last_name` (string, required): User's last name
- `role` (string, required): One of: `admin`, `customer`, `finance_manager`, `manager`

**Success Response (201):**

```json
{
  "success": true,
  "data": {
    "message": "Registration successful. Please verify your email with the OTP sent.",
    "user": {
      "id": "...",
      "email": "user@example.com",
      "phone": "+923001234567",
      "first_name": "John",
      "last_name": "Doe",
      "role": "customer",
      "is_verified": false,
      "created_at": "2026-01-02T..."
    }
  },
  "trace_id": "..."
}
```

**Error Responses:**

- `409 Conflict`: User already exists
- `400 Bad Request`: Invalid phone format or role
- `500 Internal Server Error`: Server error

---

### 3. Verify OTP

Verify the OTP sent to user's email.

**Endpoint:** `POST /api/v1/auth/verify-otp`

**Auth Required:** No

**Request Body:**

```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**Fields:**

- `email` (string, required): User's email address
- `code` (string, required): 6-digit OTP code

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Email verified successfully. You can now login."
  },
  "trace_id": "..."
}
```

**Error Responses:**

- `400 Bad Request`: Invalid OTP or OTP not found
- `500 Internal Server Error`: Server error

---

### 4. Resend OTP

Request a new OTP to be sent.

**Endpoint:** `POST /api/v1/auth/resend-otp`

**Auth Required:** No

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "OTP sent successfully. Please check your email."
  },
  "trace_id": "..."
}
```

**Error Responses:**

- `404 Not Found`: User not found
- `429 Too Many Requests`: OTP already sent, wait before requesting again
- `500 Internal Server Error`: Server error

---

### 5. Login

Authenticate and receive access & refresh tokens.

**Endpoint:** `POST /api/v1/auth/login`

**Auth Required:** No

**Rate Limit:** 5 requests per minute

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "Password123"
}
```

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "...",
      "email": "user@example.com",
      "phone": "+923001234567",
      "first_name": "John",
      "last_name": "Doe",
      "role": "customer",
      "is_verified": true,
      "created_at": "2026-01-02T..."
    }
  },
  "trace_id": "..."
}
```

**Token Details:**

- `access_token`: Short-lived token (15 minutes) for API access
- `refresh_token`: Long-lived token (7 days) for refreshing access tokens

**Error Responses:**

- `401 Unauthorized`: Invalid credentials
- `403 Forbidden`: User not verified or inactive
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

### 6. Refresh Token

Get a new access token using refresh token.

**Endpoint:** `POST /api/v1/auth/refresh-token`

**Auth Required:** No

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "...",
      "email": "user@example.com",
      ...
    }
  },
  "trace_id": "..."
}
```

**Error Responses:**

- `401 Unauthorized`: Invalid, expired, or revoked token
- `500 Internal Server Error`: Server error

---

### 7. Get Profile

Get authenticated user's profile.

**Endpoint:** `GET /api/v1/auth/profile`

**Auth Required:** Yes

**Headers:**

```
Authorization: Bearer <access_token>
```

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "id": "...",
    "email": "user@example.com",
    "phone": "+923001234567",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer",
    "is_verified": true,
    "created_at": "2026-01-02T..."
  },
  "trace_id": "..."
}
```

**Error Responses:**

- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

---

## Role-Based Endpoints

### 8. Admin Users Endpoint (Example)

Access to admin-only operations.

**Endpoint:** `GET /api/v1/admin/users`

**Auth Required:** Yes

**Required Role:** `admin`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Admin users endpoint"
  },
  "trace_id": "..."
}
```

**Error Responses:**

- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Insufficient permissions
- `500 Internal Server Error`: Server error

---

### 9. Finance Reports (Example)

Access to finance operations.

**Endpoint:** `GET /api/v1/finance/reports`

**Auth Required:** Yes

**Required Roles:** `admin`, `finance_manager`

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Finance reports endpoint"
  },
  "trace_id": "..."
}
```

---

### 10. Manager Dashboard (Example)

Access to manager operations.

**Endpoint:** `GET /api/v1/manager/dashboard`

**Auth Required:** Yes

**Required Roles:** `admin`, `manager`

**Success Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Manager dashboard endpoint"
  },
  "trace_id": "..."
}
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Bad Request | Invalid request format or parameters |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Global Rate Limit:** 10 requests per minute per IP/user
- **Login Rate Limit:** 5 requests per minute per IP

When rate limit is exceeded, the API returns:

```json
{
  "success": false,
  "error": "rate limit exceeded, please try again later",
  "trace_id": "..."
}
```

**HTTP Status Code:** 429 Too Many Requests

---

## Pakistani Phone Number Validation

Phone numbers must follow the Pakistani format:

**Format:** `+923XXXXXXXXX`

**Examples:**

- ✅ `+923001234567`
- ✅ `+923211234567`
- ✅ `+923451234567`
- ❌ `03001234567` (missing country code)
- ❌ `+92301234567` (missing digit)
- ❌ `+923001234` (too short)

The API automatically normalizes phone numbers:

- Converts `03001234567` → `+923001234567`
- Converts `923001234567` → `+923001234567`

---

## Authentication Flow

### Registration & Login Flow

1. **Register** → Receive OTP via email
2. **Verify OTP** → Account activated
3. **Login** → Receive access_token & refresh_token
4. **Use access_token** → Access protected endpoints
5. **Token expires** → Use refresh_token to get new access_token

### Token Lifecycle

```
Access Token:  [15 minutes lifetime]
              ↓
           Expires
              ↓
Refresh Token: [7 days lifetime] → Get new access_token
              ↓
           Expires
              ↓
         Login again
```

---

## Testing with cURL

### Complete Flow Example

```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "phone": "+923001234567",
    "password": "Test123456",
    "first_name": "Test",
    "last_name": "User",
    "role": "customer"
  }'

# 2. Check terminal for OTP, then verify
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "code": "YOUR_OTP_HERE"
  }'

# 3. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123456"
  }'

# 4. Save the access_token from response, then use it
export TOKEN="your_access_token_here"

# 5. Get profile
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## Postman Collection

You can import these endpoints into Postman for easier testing:

1. Create a new collection
2. Add environment variables:
   - `base_url`: `http://localhost:8080/api/v1`
   - `token`: (will be set after login)
3. Add all endpoints from this documentation
4. Use `{{base_url}}` and `{{token}}` in your requests

---

## WebSocket Support (Future)

WebSocket endpoints for real-time features will be added in future versions.

---

## Pagination (Future)

List endpoints will support pagination with these query parameters:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `sort`: Sort field (e.g., `created_at`)
- `order`: Sort order (`asc` or `desc`)

---

## Changelog

### Version 1.0.0 (Current)

- Initial release
- Authentication & authorization
- User management
- OTP verification
- Rate limiting
- Role-based access control

---

**For more information, see the [README.md](README.md)**
