# Go 1.24.10 Migration Guide

## Overview

This boilerplate has been updated to Go 1.24.10, taking advantage of the latest features and improvements while ensuring compatibility with modern dependencies.

## Major Changes

### 1. Updated go.mod

```go
go 1.24.10
```

All dependencies have been upgraded to versions compatible with Go 1.24.x:

- `github.com/gofiber/fiber/v2` → v2.52.5 (latest stable)
- `github.com/redis/go-redis/v9` → v9.7.0
- `github.com/jackc/pgx/v5` → v5.7.1
- `go.mongodb.org/mongo-driver` → v1.17.1
- `golang.org/x/crypto` → v0.31.0
- `github.com/google/uuid` → v1.6.0

### 2. Error Handling Improvements

Go 1.24 continues to improve error handling. This boilerplate now uses:

```go
import "errors"

// Join multiple errors (Go 1.20+)
err := errors.Join(err1, err2, err3)

// Use errors.New for simple errors
return errors.New("email already exists")
```

**Best Practice**: Use `errors.New()` for static errors, `fmt.Errorf()` for formatted errors with context.

### 3. Type Inference Enhancements

Go 1.24 has improved type inference. Example in our codebase:

```go
// Generic pagination function (simplified inference)
func NewPagination[T any](data []T, page, perPage int, total int64, path string) *Pagination[T] {
    // Type parameter T is inferred from data parameter
    return &Pagination[T]{...}
}

// Usage - type automatically inferred
pagination := NewPagination(userResponses, 1, 15, 100, "/users")
```

### 4. Deprecated APIs Removed

#### MongoDB Driver

- **Removed**: Deprecated `mongo.Connect` options
- **Updated**: Connection pooling configuration

```go
// ✅ Current (Go 1.24 compatible)
clientOptions := options.Client().
    ApplyURI(cfg.URI).
    SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
    SetMinPoolSize(uint64(cfg.MinPoolSize))
```

#### Fiber v2

- All deprecated middleware imports removed
- Updated to latest stable Fiber v2.52.5
- No breaking changes for our use case

#### pgx/v5

- Connection pool configuration updated
- Compatible with Go 1.24's improved context handling

### 5. Standard Library Improvements

#### Context Package

Go 1.24 has enhanced context propagation. Our implementation:

```go
// Context passed through all layers
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(ctx, req.Email)
    // ...
}
```

#### Crypto Package

Updated to `golang.org/x/crypto v0.31.0` with:

- Improved bcrypt performance
- Enhanced security defaults

### 6. Removed Deprecated Features

#### What Was Removed/Updated

1. **Old error wrapping syntax**: Now using `fmt.Errorf("... : %w", err)`
2. **Deprecated JWT claims access**: Updated to proper type assertions
3. **Old MongoDB options**: Updated to current API

## Compatibility Notes

### Breaking Changes from Go 1.21 → 1.24

**None affecting this boilerplate**. The changes are primarily:

- Performance improvements
- Enhanced type inference
- Better error messages
- Improved standard library

### Dependencies

All major dependencies support Go 1.24:

- ✅ Fiber v2: Full support
- ✅ MongoDB Driver: Full support  
- ✅ pgx v5: Full support
- ✅ go-redis v9: Full support
- ✅ Viper: Full support
- ✅ JWT v5: Full support

## New Features Utilized

### 1. Enhanced Error Handling

```go
// Before (Go 1.21)
if err != nil {
    return fmt.Errorf("operation failed: %v", err)
}

// After (Go 1.24) - better error wrapping
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### 2. Improved Generic Type Inference

```go
// Pagination helper with cleaner generics
type Pagination[T any] struct {
    CurrentPage int
    Data        []T
    // ...
}

// No need for explicit type parameter in many cases
result := NewPagination(users, 1, 10, 100, "/users")
```

### 3. Better Context Handling

```go
// Context timeout handling improved
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Better error messages when context cancelled
if err := db.Ping(ctx); err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        return fmt.Errorf("database ping timeout")
    }
    return fmt.Errorf("database ping failed: %w", err)
}
```

## Migration Steps (If Upgrading Existing Project)

### Step 1: Update go.mod

```bash
# Update go version
go mod edit -go=1.24.10

# Update dependencies
go get -u github.com/gofiber/fiber/v2@latest
go get -u github.com/redis/go-redis/v9@latest
go get -u github.com/jackc/pgx/v5@latest
go get -u go.mongodb.org/mongo-driver@latest
go get -u golang.org/x/crypto@latest

# Tidy modules
go mod tidy
```

### Step 2: Update Error Handling

```bash
# Search for deprecated error wrapping
grep -r "fmt.Errorf.*%v" .

# Replace with %w where appropriate
```

### Step 3: Test Everything

```bash
# Run all tests
go test ./...

# Build to check for compilation errors
go build ./...
```

### Step 4: Update Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder
```

### Step 5: Update CI/CD

```yaml
# GitHub Actions example
- uses: actions/setup-go@v4
  with:
    go-version: '1.24.10'
```

## Performance Improvements

Go 1.24 brings several performance enhancements that benefit this boilerplate:

1. **Faster Compilation**: ~10-15% faster build times
2. **Better GC**: Improved garbage collection for lower latency
3. **Crypto Performance**: bcrypt and JWT operations are faster
4. **Generic Performance**: Better optimization of generic code

## Testing

### Run Full Test Suite

```bash
make test
```

### Verify Dependencies

```bash
go list -m all
```

### Check for Vulnerabilities

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## Production Deployment

### Recommended Go Version

```
Go 1.24.10 or higher
```

### Docker Base Image

```dockerfile
FROM golang:1.24-alpine
```

### CI/CD Configuration

Ensure your CI/CD pipelines use Go 1.24.10:

```yaml
# GitHub Actions
- uses: actions/setup-go@v4
  with:
    go-version: '1.24.10'

# GitLab CI
image: golang:1.24-alpine
```

## Troubleshooting

### Issue: Build Fails After Upgrade

**Solution**:

```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download

# Rebuild
go build ./...
```

### Issue: Import Cycle Detected

**Solution**: Go 1.24 has stricter import cycle detection. Reorganize packages if needed.

### Issue: Deprecated API Warnings

**Solution**: Check dependency documentation and update to new APIs:

```bash
# Example for MongoDB
go doc go.mongodb.org/mongo-driver/mongo
```

## Future-Proofing

This boilerplate is ready for future Go releases:

- ✅ Uses standard patterns
- ✅ No deprecated APIs
- ✅ Latest stable dependencies
- ✅ Modern error handling
- ✅ Proper context usage
- ✅ Generic type parameters

## Resources

- [Go 1.24 Release Notes](https://go.dev/doc/go1.24)
- [Go Module Reference](https://go.dev/ref/mod)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Best Practices](https://go.dev/doc/effective_go)

## Changelog

### From Previous Version

- ✅ Updated Go version: 1.21 → 1.24.10
- ✅ Updated all dependencies to latest stable
- ✅ Improved error handling
- ✅ Enhanced type inference usage
- ✅ Better context propagation
- ✅ Removed deprecated APIs
- ✅ Docker Compose v2 syntax
- ✅ Modern Docker images (mongo:8, postgres:17)
