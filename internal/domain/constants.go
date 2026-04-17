package domain

const (
	// Env
	EnvProduction  = "production"
	EnvDevelopment = "development"

	// Env var keys
	EnvKeyDatabaseURL        = "DATABASE_URL"
	EnvKeyAppEnv             = "AGBALUMO_ENV"
	EnvKeyAppURL             = "APP_URL"
	EnvKeyBaseURL            = "BASE_URL"
	EnvKeyGoogleClientID     = "GOOGLE_CLIENT_ID"
	EnvKeyGoogleClientSecret = "GOOGLE_CLIENT_SECRET"
	EnvKeyGoogleMapsAPIKey   = "GOOGLE_MAPS_API_KEY" // #nosec G101 - This is an env var name, not a credential
	EnvKeyAdminCode          = "ADMIN_CODE"
	EnvKeyMockAuth           = "MOCK_AUTH"
	EnvKeySessionSecret      = "SESSION_SECRET"
	EnvKeyDevAuthEmail       = "DEV_AUTH_EMAIL"
	EnvKeyUploadDir          = "UPLOAD_DIR"
	EnvKeyRateLimitRate      = "RATE_LIMIT_RATE"
	EnvKeyRateLimitBurst     = "RATE_LIMIT_BURST"
	EnvKeySlowQueryThreshold = "SLOW_QUERY_THRESHOLD_MS"

	// Audit
	SeparatorLine = "--------------------------------"

	// File paths
	DefaultUploadDir = "ui/static/uploads"

	// Database
	DefaultDatabaseURL = ".tester/data/agbalumo.db"
	SQLiteDriver       = "sqlite"
	SQLiteMemory       = ":memory:"

	// Date formats
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02T15:04"

	// Templates
	TemplateError = "error.html"
	TemplateBase  = "base.html"
	TemplateIndex = "index.html"

	// Paths/Routes
	PathAdmin         = "/admin"
	PathAdminListings = "/admin/listings"
	PathListings      = "/listings"
	PathProfile       = "/profile"
	PathLogin         = "/login"
	PathListingID     = "/listings/:id"

	// File extensions
	ExtJPG      = ".jpg"
	ExtJPEG     = ".jpeg"
	ExtCSS      = ".css"
	ExtHTML     = ".html"
	ExtGo       = ".go"
	ExtMarkdown = ".md"

	// Field names (reused in CLI/Forms)
	FieldTitle       = "title"
	FieldDescription = "description"
	FieldEmail       = "email"
	FieldPhone       = "phone"
	FieldAddress     = "address"
	FieldWhatsApp    = "whatsapp"
	FieldWebsite     = "website"
	FieldImageURL    = "image-url"
	FieldDeadline    = "deadline"
	FieldEventStart  = "event-start"
	FieldEventEnd    = "event-end"
	FieldSkills      = "skills"
	FieldJobStart    = "job-start"
	FieldApplyURL    = "apply-url"
	FieldCompany     = "company"
	FieldPayRange    = "pay-range"
	FieldType        = "type"
	FieldCity        = "city"

	// Context Keys
	CtxKeyUser = "User"

	// Sessions
	SessionName          = "auth_session"
	SessionKeyOAuthState = "oauth_state"

	// Fields (Additional)
	FieldStatus    = "status"
	FieldFeatured  = "featured"
	FieldCreatedAt = "created_at"

	// Headers
	HeaderHXTrigger = "HX-Trigger"

	// HTMX
	TriggerListingUpdatedPrefix = "listing-updated-"

	// Params
	ParamSource      = "source"
	ParamSourceAdmin = "admin"

	SessionKeyUserID = "user_id"
	FlashMessageKey  = "message"

	// Messages
	MsgFailedToOpenDB        = "Failed to open DB"
	MsgFailedToUpdateListing = "Failed to update listing"
	MsgFailedToCreateListing = "Failed to create listing"
	MsgFailedToLogin         = "Failed to login"

	// Protocols
	ProtoHTTPS = "https://"
)



