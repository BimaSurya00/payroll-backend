# Auth Module API Documentation

## Base URL
```
http://localhost:8080/api/v1/auth
```

---

## Endpoints

### 1. Register

Create a new user account.

**Endpoint:** `POST /auth/register`

**Authentication:** Not required (Public endpoint)

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "Password123!"
}
```

**Field Descriptions:**
- `name` (required): User's full name (min 3, max 100 characters)
- `email` (required): User's email address (must be valid email format)
- `password` (required): Account password (min 8 characters, must meet password strength requirements)

**Password Requirements:**
- Minimum 8 characters
- Must contain uppercase and lowercase letters
- Must contain at least one number
- Must contain at least one special character

**Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "USER",
      "isActive": true,
      "createdAt": "2024-01-28T10:30:00Z",
      "updatedAt": "2024-01-28T10:30:00Z"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": 1706464200,
    "tokenType": "Bearer"
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "statusCode": 409,
  "message": "email already exists",
  "error": null
}
```

**Error Response (422 Validation Error):**
```json
{
  "success": false,
  "statusCode": 422,
  "message": "Validation failed",
  "error": {
    "email": "Email is required",
    "password": "Password must be at least 8 characters"
  }
}
```

---

### 2. Login

Authenticate with email and password.

**Endpoint:** `POST /auth/login`

**Authentication:** Not required (Public endpoint)

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "Password123!"
}
```

**Field Descriptions:**
- `email` (required): User's email address
- `password` (required): User's password

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "USER",
      "isActive": true,
      "createdAt": "2024-01-28T10:30:00Z",
      "updatedAt": "2024-01-28T10:30:00Z"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": 1706464200,
    "tokenType": "Bearer"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "statusCode": 401,
  "message": "invalid credentials",
  "error": null
}
```

**Error Response (401 Unauthorized - Deactivated Account):**
```json
{
  "success": false,
  "statusCode": 401,
  "message": "account is deactivated",
  "error": null
}
```

---

### 3. Refresh Token

Get a new access token using a refresh token.

**Endpoint:** `POST /auth/refresh`

**Authentication:** Not required (Public endpoint)

**Request Body:**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "userId": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Field Descriptions:**
- `refreshToken` (required): Valid refresh token from login/register response
- `userId` (required): User ID associated with the refresh token (required for security)

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Token refreshed successfully",
  "data": {
    "user": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "USER",
      "isActive": true,
      "createdAt": "2024-01-28T10:30:00Z",
      "updatedAt": "2024-01-28T10:30:00Z"
    },
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": 1706464200,
    "tokenType": "Bearer"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "statusCode": 401,
  "message": "Token refresh failed",
  "error": "Invalid or expired refresh token"
}
```

---

### 4. Logout

Logout from current device by invalidating the refresh token.

**Endpoint:** `POST /auth/logout`

**Authentication:** Required (Bearer token)

**Request Headers:**
```
Authorization: Bearer <your_access_token>
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Field Descriptions:**
- `refresh_token` (required): Refresh token to invalidate

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Logout successful",
  "data": null
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "statusCode": 400,
  "message": "Invalid request body",
  "error": "refresh_token is required"
}
```

---

### 5. Logout All

Logout from all devices by invalidating all refresh tokens.

**Endpoint:** `POST /auth/logout-all`

**Authentication:** Required (Bearer token)

**Request Headers:**
```
Authorization: Bearer <your_access_token>
```

**Request Body:** None required

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Logged out from all devices successfully",
  "data": null
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "success": false,
  "statusCode": 500,
  "message": "Logout from all devices failed",
  "error": "Database connection error"
}
```

---

## Token Usage

### Access Token
- Used for accessing protected endpoints
- Short-lived expiration (configured in JWT settings)
- Sent in the `Authorization` header: `Bearer <access_token>`

**Example:**
```bash
curl -X GET http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Refresh Token
- Used to obtain new access tokens
- Longer-lived than access tokens
- Stored securely (typically in localStorage or httpOnly cookie)
- Must be sent with `userId` for security

**Token Refresh Flow:**
1. Access token expires (401 Unauthorized response)
2. Client sends refresh token with userId to `/auth/refresh`
3. Server validates and returns new access + refresh tokens
4. Client updates stored tokens
5. Client retries original request with new access token

---

## User Roles

The system supports the following roles:

- `USER`: Regular user with basic permissions
- `ADMIN`: Can manage employees and perform administrative tasks
- `SUPER_USER`: Full system access including delete operations

Role is automatically assigned during registration (default: `USER`) and can be modified by administrators.

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid request body |
| 401 | Unauthorized - Invalid credentials or token |
| 409 | Conflict - Email already exists |
| 422 | Validation Error - Invalid input data |
| 500 | Internal Server Error |

---

## Validation Rules

### Register
- `name`: Required, min 3 characters, max 100 characters
- `email`: Required, valid email format, must be unique
- `password`: Required, min 8 characters, must meet password strength requirements

### Login
- `email`: Required, valid email format
- `password`: Required

### Refresh Token
- `refreshToken`: Required, valid JWT token
- `userId`: Required, valid UUID format

### Logout
- `refresh_token`: Required, valid JWT token

---

## Password Strength Requirements

The password must meet the following criteria:
- At least 8 characters long
- Contains both uppercase (A-Z) and lowercase (a-z) letters
- Contains at least one number (0-9)
- Contains at least one special character (!@#$%^&* etc.)

**Examples of valid passwords:**
- `Password123!`
- `MySecure@Pass2024`
- `Str0ng!Pass`

**Examples of invalid passwords:**
- `password` (no uppercase, number, or special character)
- `Password` (no number or special character)
- `Pass123!` (less than 8 characters)

---

## Postman Collection JSON

Import this JSON into Postman:

```json
{
  "info": {
    "name": "Auth Module API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080/api/v1/auth",
      "type": "string"
    },
    {
      "key": "accessToken",
      "value": "your_access_token_here",
      "type": "string"
    },
    {
      "key": "refreshToken",
      "value": "your_refresh_token_here",
      "type": "string"
    },
    {
      "key": "userId",
      "value": "your_user_id_here",
      "type": "string"
    }
  ],
  "item": [
    {
      "name": "1. Register",
      "request": {
        "method": "POST",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"John Doe\",\n  \"email\": \"john.doe@example.com\",\n  \"password\": \"Password123!\"\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/register",
          "host": ["{{baseUrl}}"],
          "path": ["register"]
        }
      }
    },
    {
      "name": "2. Login",
      "request": {
        "method": "POST",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"john.doe@example.com\",\n  \"password\": \"Password123!\"\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/login",
          "host": ["{{baseUrl}}"],
          "path": ["login"]
        }
      }
    },
    {
      "name": "3. Refresh Token",
      "request": {
        "method": "POST",
        "header": [],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"refreshToken\": \"{{refreshToken}}\",\n  \"userId\": \"{{userId}}\"\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/refresh",
          "host": ["{{baseUrl}}"],
          "path": ["refresh"]
        }
      }
    },
    {
      "name": "4. Logout",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"refresh_token\": \"{{refreshToken}}\"\n}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/logout",
          "host": ["{{baseUrl}}"],
          "path": ["logout"]
        }
      }
    },
    {
      "name": "5. Logout All",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{accessToken}}"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{}",
          "options": {
            "raw": {
              "language": "json"
            }
          }
        },
        "url": {
          "raw": "{{baseUrl}}/logout-all",
          "host": ["{{baseUrl}}"],
          "path": ["logout-all"]
        }
      }
    }
  ]
}
```

---

## Test Flow in Postman

### Recommended Testing Sequence:

1. **Register a new user**
   - Call `POST /auth/register`
   - Save the response tokens to variables
   - Verify user is created with default `USER` role

2. **Login with credentials**
   - Call `POST /auth/login`
   - Save the new tokens to variables
   - Verify access token and refresh token are returned

3. **Access protected endpoint**
   - Use the access token in `Authorization` header
   - Example: `GET /api/v1/employees`
   - Verify successful access

4. **Refresh token when expired**
   - Wait for access token to expire (or test with invalid token)
   - Call `POST /auth/refresh` with refresh token and userId
   - Save new tokens to variables
   - Retry protected endpoint with new access token

5. **Logout from current device**
   - Call `POST /auth/logout` with access token and refresh token
   - Verify tokens are invalidated
   - Attempt to use refresh token (should fail)

6. **Login again and test logout all**
   - Login to get new tokens
   - Login from another device/session (optional)
   - Call `POST /auth/logout-all`
   - Verify all sessions are terminated

---

## Security Best Practices

### For Clients:
1. **Store tokens securely:**
   - Access token: Memory (variable/state)
   - Refresh token: httpOnly cookie or secure storage

2. **Handle token expiration:**
   - Implement automatic token refresh
   - Show user-friendly error messages
   - Redirect to login if refresh fails

3. **Clear tokens on logout:**
   - Remove tokens from storage
   - Clear any cached user data
   - Redirect to login page

### For Developers:
1. **Never expose tokens in URL**
2. **Always use HTTPS in production**
3. **Implement rate limiting on auth endpoints**
4. **Monitor suspicious login activities**
5. **Use environment variables for JWT secrets**

---

## Notes

1. **Account Activation**: New accounts are automatically activated (`isActive: true`)
2. **Token Storage**: Refresh tokens are stored in Redis for validation and revocation
3. **Concurrent Sessions**: Multiple devices can be logged in simultaneously (use Logout All to terminate)
4. **Token Expiration**: Access tokens expire based on JWT configuration (typically 15-60 minutes)
5. **Refresh Token Rotation**: New refresh tokens are issued on each refresh for enhanced security
6. **User ID**: Use the `id` from the user object (not email) for userId in refresh token requests
