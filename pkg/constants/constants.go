package constants

const (
	// Error codes
	SuccessCode         = 0
	ErrorCode           = 1000
	ParamErrorCode      = 1001
	UnauthorizedCode    = 1002
	PermissionDeniedCode = 1003
	NotFoundCode        = 1004
	AlreadyExistsCode   = 1005
	OperationFailedCode = 1006
	RateLimitCode       = 1007

	UserNotFoundCode    = 2001
	PasswordErrorCode   = 2002
	UserLockedCode      = 2003

	CourseNotFoundCode  = 3001
	CourseClosedCode    = 3002

	TokenInvalidCode    = 4001
	TokenExpiredCode    = 4002

	// User roles
	RoleAdmin = "admin"
	RoleUser  = "user"

	// Resource types
	ResourceTypeVideo = "video"
	ResourceTypeImage = "image"

	// Course hour types
	HourTypeVideo = "video"

	// JWT context keys
	ContextKeyUserID = "user_id"
	ContextKeyEmail  = "email"
	ContextKeyRole   = "role"
	ContextKeyJTI    = "jti"

	// LDAP sync actions
	LdapSyncActionFull = "full"
	LdapSyncActionIncremental = "incremental"

	// LDAP sync status
	LdapSyncStatusRunning = "running"
	LdapSyncStatusSuccess = "success"
	LdapSyncStatusFailed  = "failed"

	// Cache keys
	CacheKeyUser          = "user:%d"
	CacheKeyCourse        = "course:%d"
	CacheKeyCategory      = "category:tree"
	CacheKeyResourceCategory = "resource_category:tree"
	CacheKeyUserDepartments  = "user:departments:%d"
	CacheKeyTokenBlacklist   = "token:blacklist:%s"

	// Cache TTL
	CacheTTLUser     = 300  // 5 minutes
	CacheTTLCourse   = 600  // 10 minutes
	CacheTTLCategory = 1800 // 30 minutes
	CacheTTLToken    = 1296000 // 15 days

	// Default values
	DefaultPageSize = 10
	DefaultPage     = 1
	MaxPageSize     = 100

	// File upload
	MaxVideoSize = 1024 * 1024 * 1024 * 2 // 2GB
	MaxImageSize = 1024 * 1024 * 10        // 10MB

	// Learning progress
	FinishThreshold = 80 // 80% to consider as finished
)
