# RBAC (Role-Based Access Control) Implementation

## Overview

This boilerplate implements production-grade Role-Based Access Control (RBAC) with three distinct roles and automatic privilege escalation for super users.

## Role Hierarchy

```
SUPER_USER (highest privilege)
    ↓
ADMIN (administrative access)
    ↓
USER (basic access)
```

### Role Definitions

| Role | Constant | Description | Access Level |
|------|----------|-------------|--------------|
| **SUPER_USER** | `constants.RoleSuperUser` | System administrator | Full access to all resources |
| **ADMIN** | `constants.RoleAdmin` | Organization administrator | Manage users and resources |
| **USER** | `constants.RoleUser` | Regular user | Access own resources |

## Key Features

### 1. Automatic Super User Privilege

**SUPER_USER automatically has access to EVERYTHING** regardless of route-specific role requirements.

```go
// In HasRole middleware
if userRole == constants.RoleSuperUser {
    return c.Next() // Bypass all role checks
}
```

This is a security best practice that ensures system administrators always have emergency access.

### 2. JWT Role Claims

Every access token contains the user's role:

```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "role": "ADMIN",
  "type": "access",
  "exp": 1701345678,
  "iat": 1701344778,
  "jti": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### 3. Reusable Middleware

```go
middleware.HasRole(requiredRoles ...string) fiber.Handler
```

**Features:**

- Accepts multiple roles (OR logic)
- Automatic SUPER_USER privilege
- Context-aware (uses JWT claims)
- Descriptive error messages

## User Module RBAC Rules

### Complete Access Matrix

| Endpoint | Method | Route | Allowed Roles | Description |
|----------|--------|-------|---------------|-------------|
| **Get Own Profile** | GET | `/api/v1/users/me` | USER, ADMIN, SUPER_USER | Any authenticated user can view their own profile |
| **Create User** | POST | `/api/v1/users` | ADMIN, SUPER_USER | Only admins can create new users |
| **Get All Users** | GET | `/api/v1/users` | ADMIN, SUPER_USER | Only admins can list all users |
| **Get User by ID** | GET | `/api/v1/users/:id` | ADMIN, SUPER_USER | Only admins can view other users |
| **Update User** | PATCH | `/api/v1/users/:id` | ADMIN, SUPER_USER | Only admins can modify users |
| **Delete User** | DELETE | `/api/v1/users/:id` | SUPER_USER | Only super users can delete users |

### Implementation Example

```go
// Get own profile - all authenticated users
users.Get("/me",
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    userHandler.GetOwnProfile,
)

// Delete user - SUPER_USER only
users.Delete("/:id",
    middleware.HasRole(constants.RoleSuperUser),
    userHandler.DeleteUser,
)

// Update user - ADMIN and SUPER_USER
users.Patch("/:id",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    userHandler.UpdateUser,
)
```

## API Examples

### 1. Login as Different Roles

#### USER Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!"
  }'
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "...",
      "role": "USER"
    },
    "accessToken": "eyJhbGci...",
    "refreshToken": "...",
    "expiresAt": 1701345678,
    "tokenType": "Bearer"
  }
}
```

### 2. Access Control Examples

#### ✅ USER accessing own profile

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <USER_token>"
```

**Response:** 200 OK

#### ❌ USER trying to list all users

```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <USER_token>"
```

**Response:** 403 Forbidden

```json
{
  "success": false,
  "message": "Insufficient permissions",
  "error": {
    "required": ["ADMIN", "SUPER_USER"],
    "actual": "USER"
  }
}
```

#### ✅ ADMIN creating a new user

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <ADMIN_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New User",
    "email": "new@example.com",
    "password": "Password123!",
    "role": "USER"
  }'
```

**Response:** 201 Created

#### ❌ ADMIN trying to delete a user

```bash
curl -X DELETE http://localhost:8080/api/v1/users/123 \
  -H "Authorization: Bearer <ADMIN_token>"
```

**Response:** 403 Forbidden (only SUPER_USER can delete)

#### ✅ SUPER_USER deleting a user

```bash
curl -X DELETE http://localhost:8080/api/v1/users/123 \
  -H "Authorization: Bearer <SUPER_USER_token>"
```

**Response:** 200 OK

## Middleware Usage Patterns

### Basic Protection

```go
// Single role required
app.Get("/admin-only", 
    middleware.JWTAuth(),
    middleware.HasRole(constants.RoleAdmin),
    handler,
)
```

### Multiple Roles (OR Logic)

```go
// Any of these roles can access
app.Get("/staff", 
    middleware.JWTAuth(),
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    handler,
)
```

### All Authenticated Users

```go
// All roles including USER
app.Get("/profile", 
    middleware.JWTAuth(),
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    handler,
)
```

### Super User Only

```go
// Critical operations
app.Delete("/system/reset", 
    middleware.JWTAuth(),
    middleware.HasRole(constants.RoleSuperUser),
    handler,
)
```

## Adding RBAC to New Modules

### Step 1: Define Module Rules

Create an access matrix for your module:

| Action | Allowed Roles |
|--------|---------------|
| Create | ADMIN, SUPER_USER |
| Read | USER, ADMIN, SUPER_USER |
| Update | ADMIN, SUPER_USER |
| Delete | SUPER_USER |

### Step 2: Apply Middleware

```go
package mymodule

import (
    "github.com/gofiber/fiber/v2"
    "github.com/itsahyarr/go-fiber-boilerplate/middleware"
    "github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
)

func RegisterRoutes(app *fiber.App, deps ...) {
    handler := handler.NewMyModuleHandler(...)
    
    // Module routes
    route := app.Group("/api/v1/mymodule", middleware.JWTAuth())
    
    // Read - all authenticated users
    route.Get("/",
        middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
        handler.List,
    )
    
    // Create - admins only
    route.Post("/",
        middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
        handler.Create,
    )
    
    // Update - admins only
    route.Patch("/:id",
        middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
        handler.Update,
    )
    
    // Delete - super users only
    route.Delete("/:id",
        middleware.HasRole(constants.RoleSuperUser),
        handler.Delete,
    )
}
```

### Step 3: Update Validation

Update DTOs to validate role values:

```go
type CreateRequest struct {
    Role string `json:"role" validate:"required,oneof=SUPER_USER ADMIN USER"`
}
```

## Security Best Practices

### ✅ Do's

1. **Always use JWTAuth before HasRole**

   ```go
   app.Get("/protected",
       middleware.JWTAuth(),      // ← First: authenticate
       middleware.HasRole(...),   // ← Then: authorize
       handler,
   )
   ```

2. **Use SUPER_USER for destructive operations**
   - Database resets
   - User deletion
   - System configuration changes

3. **Be explicit with role requirements**

   ```go
   // Good - clear intent
   middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser)
   
   // Bad - magic strings
   middleware.HasRole("ADMIN", "SUPER_USER")
   ```

4. **Test each role's access**
   - Write integration tests for each role
   - Verify 403 responses for unauthorized access

### ❌ Don'ts

1. **Don't hardcode roles in handlers**

   ```go
   // Bad
   if userRole != "ADMIN" {
       return errors.New("forbidden")
   }
   
   // Good - let middleware handle it
   ```

2. **Don't skip JWT validation**

   ```go
   // Bad - no authentication
   app.Get("/users", middleware.HasRole(...), handler)
   
   // Good - authenticate first
   app.Get("/users", middleware.JWTAuth(), middleware.HasRole(...), handler)
   ```

3. **Don't mix authentication and authorization**
   - Keep JWTAuth and HasRole as separate middlewares
   - Better composability and testability

## Error Responses

### 401 Unauthorized (Missing/Invalid Token)

```json
{
  "success": false,
  "message": "Invalid or expired token",
  "error": "token expired"
}
```

### 403 Forbidden (Insufficient Permissions)

```json
{
  "success": false,
  "message": "Insufficient permissions",
  "error": {
    "required": ["ADMIN", "SUPER_USER"],
    "actual": "USER"
  }
}
```

## Testing RBAC

### Manual Testing Checklist

For each protected endpoint:

- [ ] Test with no token → 401
- [ ] Test with expired token → 401
- [ ] Test with USER role → verify access
- [ ] Test with ADMIN role → verify access
- [ ] Test with SUPER_USER role → verify access (should always work)
- [ ] Test with wrong role → 403

### Unit Test Example

```go
func TestHasRoleMiddleware(t *testing.T) {
    tests := []struct {
        name           string
        userRole       string
        requiredRoles  []string
        expectedStatus int
    }{
        {
            name:           "USER accessing USER endpoint",
            userRole:       constants.RoleUser,
            requiredRoles:  []string{constants.RoleUser},
            expectedStatus: 200,
        },
        {
            name:           "USER accessing ADMIN endpoint",
            userRole:       constants.RoleUser,
            requiredRoles:  []string{constants.RoleAdmin},
            expectedStatus: 403,
        },
        {
            name:           "SUPER_USER accessing any endpoint",
            userRole:       constants.RoleSuperUser,
            requiredRoles:  []string{constants.RoleAdmin},
            expectedStatus: 200, // SUPER_USER bypasses all checks
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Production Considerations

### 1. Role Assignment

**On Registration:**

- New users get `USER` role by default
- Never allow users to self-assign ADMIN or SUPER_USER

**Promotion:**

- Create dedicated endpoint for role changes
- Require SUPER_USER authentication
- Log all role changes for audit

### 2. Audit Logging

```go
// Log all privilege escalations
logger.Info("Role changed",
    "adminId", adminID,
    "targetUserId", userID,
    "oldRole", oldRole,
    "newRole", newRole,
)
```

### 3. Database Indexes

```sql
-- MongoDB
db.users.createIndex({ "role": 1 })

-- PostgreSQL
CREATE INDEX idx_users_role ON users(role);
```

### 4. Monitoring

Monitor for:

- Failed authorization attempts (potential attack)
- SUPER_USER access patterns
- Role change frequency

## Migration from Old Roles

If upgrading from `admin`/`user` to `SUPER_USER`/`ADMIN`/`USER`:

```bash
# MongoDB
db.users.updateMany(
  { role: "admin" },
  { $set: { role: "ADMIN" } }
)

db.users.updateMany(
  { role: "user" },
  { $set: { role: "USER" } }
)

# Manually assign SUPER_USER to system admins
db.users.updateOne(
  { email: "sysadmin@example.com" },
  { $set: { role: "SUPER_USER" } }
)
```

## Resources

- [OWASP Access Control Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Access_Control_Cheat_Sheet.html)
- [NIST RBAC Standard](https://csrc.nist.gov/projects/role-based-access-control)
- [OAuth 2.0 Scopes Best Practices](https://www.oauth.com/oauth2-servers/scope/)
