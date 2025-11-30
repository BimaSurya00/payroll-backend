package constants

const (
	// User Roles (RBAC)
	RoleSuperUser = "SUPER_USER"
	RoleAdmin     = "ADMIN"
	RoleUser      = "USER"

	// Token Types
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"

	// Context Keys
	ContextKeyUserID   = "user_id"
	ContextKeyUserRole = "user_role"

	// Cache Keys
	CacheKeyRefreshToken       = "refresh_token:"
	CacheKeyUserSession        = "user_session:"
	CacheKeyRefreshTokenPrefix = "refresh_token:%s:%s" // user_id:token_id

	// Pagination
	DefaultPage    = 1
	DefaultPerPage = 15
	MaxPerPage     = 100
)
