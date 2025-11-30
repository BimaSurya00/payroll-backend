# camelCase Struct Tags Migration

## Overview

All struct tags (`json`, `bson`, `gorm`) have been updated to use camelCase naming convention. This improves consistency with modern API standards and JavaScript/TypeScript frontend integration.

## Changes Summary

### General Pattern

```go
// Before (snake_case)
Field string `json:"field_name" bson:"field_name"`

// After (camelCase)
Field string `json:"fieldName" bson:"fieldName"`
```

**Note:** Field names remain unchanged. Only tag values were updated.

## Files Modified

### 1. shared/entity/base.go

```go
// Before
type BaseEntity struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// After
type BaseEntity struct {
    ID        string    `json:"id" bson:"_id,omitempty"`
    CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
```

### 2. internal/user/entity/user.go

```go
// Before
type User struct {
    entity.BaseEntity `bson:",inline"`
    Name              string `json:"name" bson:"name"`
    Email             string `json:"email" bson:"email"`
    Password          string `json:"-" bson:"password"`
    Role              string `json:"role" bson:"role"`
    IsActive          bool   `json:"is_active" bson:"is_active"`
}

// After
type User struct {
    entity.BaseEntity `bson:",inline"`
    Name              string `json:"name" bson:"name"`
    Email             string `json:"email" bson:"email"`
    Password          string `json:"-" bson:"password"`
    Role              string `json:"role" bson:"role"`
    IsActive          bool   `json:"isActive" bson:"isActive"`
}
```

### 3. internal/user/dto/user_response.go

```go
// Before
type UserResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// After
type UserResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    IsActive  bool      `json:"isActive"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

### 4. internal/user/dto/update_user.go

```go
// Before
type UpdateUserRequest struct {
    Name     *string `json:"name,omitempty" validate:"..."`
    Email    *string `json:"email,omitempty" validate:"..."`
    Password *string `json:"password,omitempty" validate:"..."`
    Role     *string `json:"role,omitempty" validate:"..."`
    IsActive *bool   `json:"is_active,omitempty"`
}

// After
type UpdateUserRequest struct {
    Name     *string `json:"name,omitempty" validate:"..."`
    Email    *string `json:"email,omitempty" validate:"..."`
    Password *string `json:"password,omitempty" validate:"..."`
    Role     *string `json:"role,omitempty" validate:"..."`
    IsActive *bool   `json:"isActive,omitempty"`
}
```

### 5. internal/auth/entity/token.go

```go
// Before
type TokenPair struct {
    AccessToken       string `json:"access_token"`
    RefreshToken      string `json:"refresh_token"`
    AccessTokenExpiry int64  `json:"access_token_expiry"`
}

type RefreshTokenData struct {
    TokenID   string    `json:"token_id"`
    UserID    string    `json:"user_id"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}

// After
type TokenPair struct {
    AccessToken       string `json:"accessToken"`
    RefreshToken      string `json:"refreshToken"`
    AccessTokenExpiry int64  `json:"accessTokenExpiry"`
}

type RefreshTokenData struct {
    TokenID   string    `json:"tokenId"`
    UserID    string    `json:"userId"`
    ExpiresAt time.Time `json:"expiresAt"`
    CreatedAt time.Time `json:"createdAt"`
}
```

### 6. internal/auth/dto/refresh.go

```go
// Before
type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
    UserID       string `json:"user_id" validate:"required"`
}

// After
type RefreshTokenRequest struct {
    RefreshToken string `json:"refreshToken" validate:"required"`
    UserID       string `json:"userId" validate:"required"`
}
```

### 7. internal/auth/dto/response.go

```go
// Before
type AuthResponse struct {
    User         *dto.UserResponse `json:"user"`
    AccessToken  string            `json:"access_token"`
    RefreshToken string            `json:"refresh_token"`
    ExpiresAt    int64             `json:"expires_at"`
    TokenType    string            `json:"token_type"`
}

// After
type AuthResponse struct {
    User         *dto.UserResponse `json:"user"`
    AccessToken  string            `json:"accessToken"`
    RefreshToken string            `json:"refreshToken"`
    ExpiresAt    int64             `json:"expiresAt"`
    TokenType    string            `json:"tokenType"`
}
```

### 8. internal/user/helper/pagination.go

```go
// Before
type Pagination[T any] struct {
    CurrentPage  int     `json:"current_page"`
    Data         []T     `json:"data"`
    FirstPageURL string  `json:"first_page_url"`
    From         int     `json:"from"`
    LastPage     int     `json:"last_page"`
    LastPageURL  string  `json:"last_page_url"`
    Links        []Link  `json:"links"`
    NextPageURL  *string `json:"next_page_url"`
    Path         string  `json:"path"`
    PerPage      int     `json:"per_page"`
    PrevPageURL  *string `json:"prev_page_url"`
    To           int     `json:"to"`
    Total        int64   `json:"total"`
}

// After
type Pagination[T any] struct {
    CurrentPage  int     `json:"currentPage"`
    Data         []T     `json:"data"`
    FirstPageURL string  `json:"firstPageUrl"`
    From         int     `json:"from"`
    LastPage     int     `json:"lastPage"`
    LastPageURL  string  `json:"lastPageUrl"`
    Links        []Link  `json:"links"`
    NextPageURL  *string `json:"nextPageUrl"`
    Path         string  `json:"path"`
    PerPage      int     `json:"perPage"`
    PrevPageURL  *string `json:"prevPageUrl"`
    To           int     `json:"to"`
    Total        int64   `json:"total"`
}
```

## Complete Tag Mapping

| Old Tag (snake_case) | New Tag (camelCase) |
|---------------------|---------------------|
| `created_at` | `createdAt` |
| `updated_at` | `updatedAt` |
| `is_active` | `isActive` |
| `access_token` | `accessToken` |
| `refresh_token` | `refreshToken` |
| `access_token_expiry` | `accessTokenExpiry` |
| `token_id` | `tokenId` |
| `user_id` | `userId` |
| `expires_at` | `expiresAt` |
| `token_type` | `tokenType` |
| `current_page` | `currentPage` |
| `first_page_url` | `firstPageUrl` |
| `last_page` | `lastPage` |
| `last_page_url` | `lastPageUrl` |
| `next_page_url` | `nextPageUrl` |
| `per_page` | `perPage` |
| `prev_page_url` | `prevPageUrl` |

## Impact on API Responses

### Before (snake_case)

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "123",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "USER",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGci...",
    "refresh_token": "uuid-here",
    "expires_at": 1701345678,
    "token_type": "Bearer"
  }
}
```

### After (camelCase)

```json
{
  "success": true,
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

## Frontend Integration

### TypeScript Interfaces

Now the API responses match TypeScript conventions:

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

interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
  tokenType: string;
}

interface Pagination<T> {
  currentPage: number;
  data: T[];
  firstPageUrl: string;
  from: number;
  lastPage: number;
  lastPageUrl: string;
  links: Link[];
  nextPageUrl: string | null;
  path: string;
  perPage: number;
  prevPageUrl: string | null;
  to: number;
  total: number;
}
```

### React/Next.js Example

```typescript
// Before - needed transformation
const fetchUser = async (id: string) => {
  const response = await fetch(`/api/users/${id}`);
  const data = await response.json();
  
  // Had to transform snake_case to camelCase
  return {
    ...data.user,
    isActive: data.user.is_active,
    createdAt: data.user.created_at,
    updatedAt: data.user.updated_at,
  };
};

// After - direct usage
const fetchUser = async (id: string) => {
  const response = await fetch(`/api/users/${id}`);
  const { user } = await response.json();
  return user; // Already in camelCase!
};
```

## Database Compatibility

### MongoDB

The `bson` tags are updated to camelCase, which affects how data is stored:

```javascript
// MongoDB Document (after migration)
{
  "_id": "123",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "USER",
  "isActive": true,
  "createdAt": ISODate("2024-01-01T00:00:00Z"),
  "updatedAt": ISODate("2024-01-01T00:00:00Z")
}
```

### Migration Script (if needed)

If you have existing data in snake_case format:

```javascript
// MongoDB migration
db.users.find().forEach(function(doc) {
  db.users.updateOne(
    { _id: doc._id },
    {
      $rename: {
        "is_active": "isActive",
        "created_at": "createdAt",
        "updated_at": "updatedAt"
      }
    }
  );
});
```

## Validation Tags (Unchanged)

Note that `validate` tags remain unchanged:

```go
// Validation tags keep their original format
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=3,max=100,trimmed_string"`
    Email string `json:"email" validate:"required,email"`
    Role  string `json:"role" validate:"required,oneof=SUPER_USER ADMIN USER"`
}
```

## Testing Updates

Update your tests to use camelCase:

```go
// Before
assert.Equal(t, "John Doe", user["name"])
assert.Equal(t, true, user["is_active"])

// After
assert.Equal(t, "John Doe", user["name"])
assert.Equal(t, true, user["isActive"])
```

## Benefits

1. **Frontend Consistency**: Matches JavaScript/TypeScript naming conventions
2. **No Transformation**: Frontend doesn't need to convert snake_case to camelCase
3. **Type Safety**: TypeScript interfaces align perfectly with API responses
4. **Modern Standards**: Follows contemporary REST API naming conventions
5. **Better DX**: Improved developer experience for frontend teams

## Breaking Changes

⚠️ **This is a breaking change for existing API consumers**

If you have existing clients consuming your API:

1. **Version your API**: Consider releasing this as v2
2. **Deprecation Period**: Keep v1 with snake_case for transition
3. **Update Documentation**: Ensure all examples use new format
4. **Notify Consumers**: Alert frontend teams of the change

## Checklist

- [x] Updated all `json` tags to camelCase
- [x] Updated all `bson` tags to camelCase
- [x] Kept `validate` tags unchanged
- [x] Updated documentation
- [x] Created TypeScript interfaces
- [x] Tested API responses
- [x] Updated integration tests
