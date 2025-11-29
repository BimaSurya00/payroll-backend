# Go Fiber Production Boilerplate

A production-ready, scalable Golang backend boilerplate using Fiber framework with clean architecture principles.

## 🚀 Features

- **Clean Architecture**: Modular design with clear separation of concerns
- **Repository Pattern**: Data access abstraction layer
- **JWT Authentication**: Secure access and refresh token implementation with rotation
- **Refresh Token Security**: UUID-based refresh tokens stored in KeyDB with expiration
- **Token Rotation**: Automatic refresh token rotation for enhanced security
- **Multiple Databases**: MongoDB, PostgreSQL, and KeyDB/Redis support
- **Validation**: Comprehensive request validation with custom validators
- **Pagination**: Laravel-style pagination helper
- **Error Handling**: Centralized error management
- **Middleware**: Authentication, logging, recovery, and CORS
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
git clone https://github.com/itsahyarr/go-fiber-boilerplate.git
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
  "refresh_token": "your-refresh-token-uuid",
  "user_id": "user-id-here"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "user": {...},
    "access_token": "new-access-token",
    "refresh_token": "new-refresh-token-uuid",
    "expires_at": 1234567890,
    "token_type": "Bearer"
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
  "refresh_token": "refresh-token-to-revoke"
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

#### Create User

```http
POST /api/v1/users
Authorization: Bearer your-access-token
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "Password123!",
  "role": "user"
}
```

#### Get Users (Paginated)

```http
GET /api/v1/users?page=1&per_page=15
Authorization: Bearer your-access-token
```

#### Get User by ID

```http
GET /api/v1/users/:id
Authorization: Bearer your-access-token
```

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

#### Delete User

```http
DELETE /api/v1/users/:id
Authorization: Bearer your-access-token
```

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

Pagination follows Laravel 11 format:

```json
{
  "current_page": 1,
  "data": [...],
  "first_page_url": "http://localhost:8080/api/v1/users?page=1&per_page=15",
  "from": 1,
  "last_page": 5,
  "last_page_url": "http://localhost:8080/api/v1/users?page=5&per_page=15",
  "links": [...],
  "next_page_url": "http://localhost:8080/api/v1/users?page=2&per_page=15",
  "path": "http://localhost:8080/api/v1/users",
  "per_page": 15,
  "prev_page_url": null,
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

Your Name - [@itsahyarr](https://github.com/itsahyarr)

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Fast HTTP framework
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [go-redis](https://github.com/redis/go-redis) - Redis client
