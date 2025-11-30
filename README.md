# Go Fiber Production Boilerplate

A production-ready, scalable Golang backend boilerplate using Fiber framework with clean architecture principles.

## 🚀 Features

- **Clean Architecture**: Modular design with clear separation of concerns
- **Repository Pattern**: Data access abstraction layer
- **JWT Authentication**: Secure access and refresh token implementation with rotation
- **RBAC (Role-Based Access Control)**: Three-tier role system (SUPER_USER, ADMIN, USER)
- **Refresh Token Security**: UUID-based refresh tokens stored in KeyDB with expiration
- **Token Rotation**: Automatic refresh token rotation for enhanced security
- **Multiple Databases**: MongoDB, PostgreSQL, and KeyDB/Redis support
- **Validation**: Comprehensive request validation with custom validators
- **Pagination**: Laravel-style pagination helper
- **camelCase API**: Modern JSON responses with camelCase field names
- **Error Handling**: Centralized error management
- **Middleware**: Authentication, logging, recovery, CORS, and RBAC
- **Docker Support**: Modern multi-container setup with health checks
- **Testing**: Unit tests with mock repositories
- **Makefile**: Convenient development commands
- **Go 1.24.10**: Latest Go version with modern features

## 📁 Project Structure

```
.
├── config/                 # Configuration management
├── database/              # Database helpers (MongoDB, PostgreSQL, KeyDB)
├── internal/              # Application modules
│   ├── auth/             # Authentication module
│   │   ├── dto/          # Data transfer objects
│   │   ├── entity/       # Domain entities
│   │   ├── handler/      # HTTP handlers
│   │   ├── helper/       # JWT utilities
│   │   ├── repository/   # Data access layer
│   │   ├── service/      # Business logic
│   │   └── routes.go     # Route definitions
│   └── user/             # User module (same structure)
├── middleware/            # HTTP middleware
├── shared/               # Shared components
│   ├── constants/        # Application constants
│   ├── dto/             # Shared DTOs
│   ├── entity/          # Base entities
│   ├── errs/            # Error definitions
│   ├── helper/          # Utility functions
│   └── validator/       # Custom validators
├── main.go              # Application entry point
├── Makefile             # Development commands
└── docker-compose.yml   # Docker orchestration
```

## 🛠️ Tech Stack

- **Framework**: Fiber v2
- **Databases**:
  - MongoDB (User data)
  - PostgreSQL (Relational data)
  - KeyDB/Redis (Cache & sessions)
- **Authentication**: JWT with bcrypt
- **Validation**: go-playground/validator
- **Config**: Viper
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (optional)
- MongoDB, PostgreSQL, KeyDB/Redis (if running locally)

## 🔧 Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/go-fiber-boilerplate.git
cd go-fiber-boilerplate
```

2. **Copy environment file**

```bash
cp .env.example .env
```

3. **Update `.env` with your configuration**

4. **Install dependencies**

```bash
make deps
```

## 🚀 Running the Application

### Using Docker (Recommended)

```bash
# Start all services
make docker-up

# View logs
docker-compose logs -f app

# Stop all services
make docker-down
```

### Running Locally

```bash
# Run the application
make run

# Run with hot reload (requires air)
make dev
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover
```

## 📚 API Documentation

### Authentication Endpoints

#### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "Password123!"
}
```

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "Password123!"
}
```

#### Refresh Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "your-refresh-token-uuid",
  "userId": "user-id-here"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "user": {...},
    "accessToken": "new-access-token",
    "refreshToken": "new-refresh-token-uuid",
    "expiresAt": 1234567890,
    "tokenType": "Bearer"
  }
}
```

**Note:** The old refresh token is automatically revoked (token rotation).

#### Logout

```http
POST /api/v1/auth/logout
Authorization: Bearer your-access-token
Content-Type: application/json

{
  "refreshToken": "refresh-token-to-revoke"
}
```

#### Logout from All Devices

```http
POST /api/v1/auth/logout-all
Authorization: Bearer your-access-token
```

This endpoint revokes all refresh tokens for the user, logging them out from all devices.

### User Endpoints

All user endpoints require authentication (Bearer token).

#### Get Own Profile

```http
GET /api/v1/users/me
Authorization: Bearer your-access-token
```

**Accessible by:** USER, ADMIN, SUPER_USER

#### Create User

```http
POST /api/v1/users
Authorization: Bearer your-access-token
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "Password123!",
  "role": "USER"
}
```

**Accessible by:** ADMIN, SUPER_USER

#### Get Users (Paginated)

```http
GET /api/v1/users?page=1&per_page=15
Authorization: Bearer your-access-token
```

**Accessible by:** ADMIN, SUPER_USER

#### Get User by ID

```http
GET /api/v1/users/:id
Authorization: Bearer your-access-token
```

**Accessible by:** ADMIN, SUPER_USER

#### Update User

```http
PATCH /api/v1/users/:id
Authorization: Bearer your-access-token
Content-Type: application/json

{
  "name": "Jane Smith",
  "email": "jane.smith@example.com"
}
```

**Accessible by:** ADMIN, SUPER_USER

#### Delete User

```http
DELETE /api/v1/users/:id
Authorization: Bearer your-access-token
```

**Accessible by:** SUPER_USER only

## 🔐 Refresh Token Security

This boilerplate implements **industry best practices** for refresh token management:

### Token Strategy

- **Access Tokens**: Short-lived JWT (15 minutes), contains user claims
- **Refresh Tokens**: Long-lived UUID v4 (7 days), stored in KeyDB
- **Separation**: Refresh tokens are NOT JWTs - they're random UUIDs for security

### Storage Pattern

```
KeyDB Key: refresh_token:{user_id}:{token_id}
Value: JSON with { token_id, user_id, expires_at, created_at }
TTL: Automatically expires after 7 days
```

### Security Features

1. **Token Rotation**: When a refresh token is used, a new pair is issued and the old refresh token is deleted
2. **Secure Format**: Refresh tokens are UUIDs, not JWTs (no data exposure)
3. **Expiration**: Automatic TTL-based expiration in KeyDB
4. **Revocation**: Can revoke specific tokens or all user tokens
5. **No Local Storage**: Refresh tokens never stored in memory - always in KeyDB
6. **User ID Required**: Refresh endpoint requires both token and user ID for extra validation

### Token Flow

```
1. Login/Register → Generate access + refresh tokens
2. Store refresh token in KeyDB with TTL
3. Client uses access token for API calls
4. Access token expires → Client sends refresh token + user ID
5. Validate refresh token in KeyDB
6. Generate NEW access + refresh tokens
7. Delete OLD refresh token (rotation)
8. Store NEW refresh token in KeyDB
```

### Logout Options

- **Single Device**: `POST /auth/logout` with specific refresh token
- **All Devices**: `POST /auth/logout-all` revokes all user's refresh tokens

## 🛡️ Role-Based Access Control (RBAC)

This boilerplate implements production-grade RBAC with three roles:

### Role Hierarchy

```
SUPER_USER → Full system access (automatically bypasses all role checks)
ADMIN      → Administrative access (manage users and resources)
USER       → Basic access (own resources only)
```

### User Module Access Matrix

| Endpoint | Method | Allowed Roles | Description |
|----------|--------|---------------|-------------|
| `/users/me` | GET | USER, ADMIN, SUPER_USER | Get own profile |
| `/users` | POST | ADMIN, SUPER_USER | Create new user |
| `/users` | GET | ADMIN, SUPER_USER | List all users |
| `/users/:id` | GET | ADMIN, SUPER_USER | Get user by ID |
| `/users/:id` | PATCH | ADMIN, SUPER_USER | Update user |
| `/users/:id` | DELETE | SUPER_USER | Delete user |

### RBAC Examples

```bash
# USER accessing own profile - ✅ Allowed
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <USER_TOKEN>"

# USER trying to list all users - ❌ 403 Forbidden
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <USER_TOKEN>"

# ADMIN creating a user - ✅ Allowed
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -d '{"name":"New User","email":"user@example.com","password":"Pass123!","role":"USER"}'

# SUPER_USER deleting a user - ✅ Allowed
curl -X DELETE http://localhost:8080/api/v1/users/123 \
  -H "Authorization: Bearer <SUPER_USER_TOKEN>"
```

**Key Feature**: SUPER_USER automatically has access to ALL endpoints, regardless of role requirements.

See [RBAC_IMPLEMENTATION.md](RBAC_IMPLEMENTATION.md) for complete documentation.

## 🎨 camelCase API Responses

All JSON responses use camelCase naming convention for seamless frontend integration:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "123",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "USER",
      "isActive": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    "accessToken": "eyJhbGci...",
    "refreshToken": "uuid-here",
    "expiresAt": 1701345678,
    "tokenType": "Bearer"
  }
}
```

**TypeScript interfaces match perfectly** - no transformation needed!

```typescript
interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}
```

See [CAMELCASE_MIGRATION.md](CAMELCASE_MIGRATION.md) for complete field mapping.

## 🔐 Custom Validators

The boilerplate includes custom validation tags:

- `password_strength`: Validates password contains uppercase, lowercase, number, and special character
- `trimmed_string`: Ensures no leading/trailing whitespace

Example usage:

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=3,trimmed_string"`
    Password string `json:"password" validate:"required,password_strength"`
}
```

## 📊 Pagination Format

Pagination follows Laravel 11 format with camelCase:

```json
{
  "currentPage": 1,
  "data": [...],
  "firstPageUrl": "http://localhost:8080/api/v1/users?page=1&per_page=15",
  "from": 1,
  "lastPage": 5,
  "lastPageUrl": "http://localhost:8080/api/v1/users?page=5&per_page=15",
  "links": [...],
  "nextPageUrl": "http://localhost:8080/api/v1/users?page=2&per_page=15",
  "path": "http://localhost:8080/api/v1/users",
  "perPage": 15,
  "prevPageUrl": null,
  "to": 15,
  "total": 67
}
```

## 🏗️ Adding New Modules

To add a new module, follow this structure:

1. Create module directory: `internal/mymodule/`
2. Add subdirectories: `dto/`, `entity/`, `handler/`, `helper/`, `repository/`, `service/`
3. Create `routes.go` for route registration
4. Implement the repository pattern
5. Register routes in `main.go`

## 🔑 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Application port | `8080` |
| `APP_ENV` | Environment | `development` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `KEYDB_HOST` | KeyDB host | `localhost` |
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_ACCESS_EXPIRY` | Access token expiry | `15m` |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry | `168h` |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License.

## 👨‍💻 Author

Your Name - [@yourusername](https://github.com/yourusername)

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Fast HTTP framework
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [go-redis](https://github.com/redis/go-redis) - Redis client
