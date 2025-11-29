# Refresh Token Implementation Guide

## Overview

This document explains the secure refresh token implementation in this boilerplate, following modern security best practices.

## Architecture Decision

### Why Separate Refresh Tokens from Access Tokens?

**Access Token (JWT)**:

- Short-lived (15 minutes)
- Contains user claims (user_id, role)
- Stateless - no server lookup required
- Used for API authentication

**Refresh Token (UUID)**:

- Long-lived (7 days)
- Random UUID v4 - no data exposure
- Stateful - stored in KeyDB
- Used only for obtaining new access tokens

### Why UUID Instead of JWT for Refresh Tokens?

1. **Security**: UUIDs contain no data - if leaked, attacker gets no information
2. **Revocation**: Easier to revoke specific tokens stored in database
3. **Size**: UUIDs are smaller than JWTs
4. **Best Practice**: Refresh tokens should be opaque tokens, not self-contained

## KeyDB Storage Structure

### Key Naming Convention

```
refresh_token:{user_id}:{token_id}

Example:
refresh_token:550e8400-e29b-41d4-a716-446655440000:a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

### Value Structure (JSON)

```json
{
  "token_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2024-12-06T10:30:00Z",
  "created_at": "2024-11-29T10:30:00Z"
}
```

### TTL (Time To Live)

- Automatically set to 7 days (configurable via `JWT_REFRESH_EXPIRY`)
- KeyDB automatically deletes expired tokens
- No manual cleanup required

## Token Rotation Flow

Token rotation is a **critical security feature** that prevents token replay attacks.

### How It Works

```
Step 1: User sends refresh token
   ↓
Step 2: Validate token exists in KeyDB
   ↓
Step 3: Delete old refresh token (IMPORTANT!)
   ↓
Step 4: Generate NEW access + refresh tokens
   ↓
Step 5: Store NEW refresh token in KeyDB
   ↓
Step 6: Return new tokens to client
```

### Why Rotation?

1. **Prevents Replay Attacks**: Old tokens become useless immediately
2. **Limits Exposure Window**: If token is stolen, it only works once
3. **Audit Trail**: Each refresh creates a new token with timestamp
4. **Industry Standard**: Recommended by OAuth 2.0 best practices

## API Examples

### 1. Login

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
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "user@example.com",
      "role": "user"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "expires_at": 1701345678,
    "token_type": "Bearer"
  }
}
```

### 2. Refresh Access Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "user": {...},
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "new-uuid-here",
    "expires_at": 1701346578,
    "token_type": "Bearer"
  }
}
```

**Important**: The `refresh_token` in the response is a NEW token. The old one is invalidated.

### 3. Logout (Single Device)

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  }'
```

### 4. Logout All Devices

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

This revokes ALL refresh tokens for the user.

## Security Considerations

### ✅ Implemented Best Practices

1. **Separate Token Types**: Access tokens (JWT) vs Refresh tokens (UUID)
2. **Token Rotation**: Old refresh token deleted on use
3. **Secure Storage**: KeyDB only, never in memory or logs
4. **Automatic Expiration**: TTL-based cleanup in KeyDB
5. **User ID Validation**: Refresh requires both token and user ID
6. **No Sensitive Data**: Refresh tokens are opaque UUIDs
7. **Revocation Support**: Single token or all tokens

### 🔒 Additional Recommendations

1. **HTTPS Only**: Always use HTTPS in production
2. **Secure Cookie**: Consider storing refresh token in HTTP-only cookie
3. **Rate Limiting**: Limit refresh token endpoint (prevents brute force)
4. **Monitoring**: Log all refresh token operations
5. **Token Fingerprinting**: Optional: tie tokens to device/IP
6. **Refresh Token Rotation Window**: Optional: allow grace period for network issues

## Code Structure

### Repository Layer (`internal/auth/repository/token_repository.go`)

```go
type TokenRepository interface {
    SetRefreshToken(ctx, userID, tokenID, expiresAt)
    GetRefreshToken(ctx, userID, tokenID)
    DeleteRefreshToken(ctx, userID, tokenID)
    DeleteAllUserRefreshTokens(ctx, userID)
    RefreshTokenExists(ctx, userID, tokenID)
}
```

### Service Layer (`internal/auth/service/auth_service.go`)

```go
func (s *authService) RefreshToken(ctx, req) {
    // 1. Validate refresh token
    // 2. Delete old token (rotation)
    // 3. Generate new tokens
    // 4. Store new refresh token
    // 5. Return new tokens
}
```

### Handler Layer (`internal/auth/handler/auth_handler.go`)

```go
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
    // 1. Parse request
    // 2. Validate input
    // 3. Call service
    // 4. Return response
}
```

## Testing Refresh Token Flow

### Manual Testing Steps

1. **Login** and save both tokens
2. **Use access token** to call protected endpoints
3. **Wait for access token expiry** (or manually expire it)
4. **Call refresh endpoint** with refresh token + user ID
5. **Verify** new tokens are returned
6. **Try using old refresh token** - should fail (rotation)
7. **Call logout-all**
8. **Try refreshing** - should fail (revoked)

### Unit Test Example

```go
func TestRefreshTokenRotation(t *testing.T) {
    // Setup
    service := setupAuthService()
    userID := "user-123"
    
    // Login and get tokens
    auth, _ := service.Login(ctx, loginReq)
    oldRefreshToken := auth.RefreshToken
    
    // Refresh tokens
    refreshReq := &dto.RefreshTokenRequest{
        RefreshToken: oldRefreshToken,
        UserID: userID,
    }
    newAuth, _ := service.RefreshToken(ctx, refreshReq)
    
    // Verify rotation
    assert.NotEqual(t, oldRefreshToken, newAuth.RefreshToken)
    
    // Try using old token (should fail)
    _, err := service.RefreshToken(ctx, refreshReq)
    assert.Error(t, err)
}
```

## Troubleshooting

### Common Issues

**Issue**: Refresh token not found

- **Cause**: Token expired or never created
- **Solution**: Check KeyDB TTL, verify SetRefreshToken was called

**Issue**: Token rotation not working

- **Cause**: DeleteRefreshToken not called before generating new token
- **Solution**: Ensure deletion happens before new token creation

**Issue**: Can't logout specific device

- **Cause**: Not passing refresh token to logout endpoint
- **Solution**: Store refresh token on client, send with logout request

## Production Checklist

- [ ] HTTPS enabled for all endpoints
- [ ] JWT_SECRET is strong and secret (min 32 characters)
- [ ] JWT_REFRESH_EXPIRY set appropriately (7 days recommended)
- [ ] KeyDB password enabled in production
- [ ] Rate limiting on /auth/refresh endpoint
- [ ] Monitoring and alerting for suspicious refresh patterns
- [ ] Regular security audits of token flow
- [ ] Consider refresh token max lifetime (e.g., 30 days absolute max)

## References

- [OAuth 2.0 Best Practices](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
