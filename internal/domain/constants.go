package domain

const (
	// Env
	EnvProduction  = "production"
	EnvDevelopment = "development"

	// Env var keys
	EnvKeyDatabaseURL = "DATABASE_URL"
	EnvKeyAppEnv      = "AGBALUMO_ENV"
	EnvKeyAppURL      = "APP_URL"
	EnvKeyBaseURL     = "BASE_URL"

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
	PathAdmin    = "/admin"
	PathListings = "/listings"
	PathProfile  = "/profile"

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
)
