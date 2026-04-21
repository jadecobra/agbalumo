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
	TemplateError          = "error.html"
	TemplateBase           = "base.html"
	TemplateIndex          = "index.html"
	TemplateAdminDashboard = "admin_dashboard.html"
	TemplateAdminListings  = "admin_listings.html"

	// Paths/Routes
	PathAdmin           = "/admin"
	PathAdminListings   = "/admin/listings"
	PathListings        = "/listings"
	PathProfile         = "/profile"
	PathLogin           = "/login"
	PathListingID       = "/listings/:id"
	PathAdminCategories = "/admin/categories"

	// File extensions
	ExtJPG      = ".jpg"
	ExtJPEG     = ".jpeg"
	ExtCSS      = ".css"
	ExtHTML     = ".html"
	ExtGo       = ".go"
	ExtMarkdown = ".md"

	// Field names (reused in CLI/Forms)
	FieldTitle             = "title"
	FieldDescription       = "description"
	FieldEmail             = "email"
	FieldPhone             = "phone"
	FieldAddress           = "address"
	FieldWhatsApp          = "whatsapp"
	FieldWebsite           = "website"
	FieldImageURL          = "image-url"
	FieldDeadline          = "deadline"
	FieldEventStart        = "event-start"
	FieldEventEnd          = "event-end"
	FieldSkills            = "skills"
	FieldJobStart          = "job-start"
	FieldApplyURL          = "apply-url"
	FieldCompany           = "company"
	FieldPayRange          = "pay-range"
	FieldType              = "type"
	FieldCity              = "city"
	FieldOwnerOrigin       = "owner_origin"
	FieldWebsiteURL        = "website_url"
	FieldContactEmail      = "contact_email"
	FieldContactPhone      = "contact_phone"
	FieldHoursOfOperation  = "hours_of_operation"
	FieldEventStartDate    = "event_start"
	FieldDeadlineDate      = "deadline_date"
	FieldContactWhatsApp   = "contact_whatsapp"
	FieldJobStartDate      = "job_start_date"
	FieldJobApplyURL       = "job_apply_url"
	FieldTopDish           = "top_dish"
	FieldRegionalSpecialty = "regional_specialty"
	FieldHeatLevel         = "heat_level"
	FieldRemoveImage       = "remove_image"

	// Context Keys
	CtxKeyUser = "User"

	// Sessions
	SessionName          = "auth_session"
	SessionKeyOAuthState = "oauth_state"

	// Fields (Additional)
	FieldStatus      = "status"
	FieldFeatured    = "featured"
	FieldCreatedAt   = "created_at"
	FieldAction      = "action"
	FieldNewCategory = "new_category"
	FieldClaimable   = "claimable"
	FieldAdminCode   = "admin_code"
	FieldCode        = "code"
	FieldName        = "name"
	FieldCSVFile     = "csv_file"
	FieldContent     = "content"

	// Headers
	HeaderHXTrigger = "HX-Trigger"

	// HTMX
	TriggerListingUpdatedPrefix = "listing-updated-"

	// HTMX Targets & Indicators
	TargetListingsContainer  = "#listings-container"
	IndicatorListingsLoading = "#listings-loading"
	TargetBody               = "body"

	// HTMX Swaps
	SwapBeforeEnd = "beforeend"
	SwapOuterHTML = "outerHTML"

	// Params
	ParamSource      = "source"
	ParamSourceAdmin = "admin"
	ParamCategory    = "category"
	ParamSort        = "sort"
	ParamOrder       = "order"
	ParamQuery       = "q"
	ParamID          = "id"
	ParamPage        = "page"
	ParamTarget      = "target"
	ParamState       = "state"
	ParamCode        = "code"
	ParamCSVFile     = "csv_file"
	ParamListingIDs  = "selectedListings"

	SessionKeyUserID = "user_id"
	FlashMessageKey  = "message"

	// Messages
	MsgFailedToOpenDB        = "Failed to open DB"
	MsgFailedToUpdateListing = "Failed to update listing"
	MsgFailedToCreateListing = "Failed to create listing"
	MsgFailedToLogin         = "Failed to login"

	// Protocols
	ProtoHTTPS = "https://"

	// Modals
	ModalModeration    = "moderationModal"
	ModalCategory      = "categoryModal"
	ModalUsers         = "usersModal"
	ModalCreateListing = "create-listing-modal"

	// Actions
	ActionOpenModal = "open-modal"

	// CSS Classes (Core)
	ClassFooter     = "footer-fruit"
	ClassEarthDark  = "bg-earth-dark"
	ClassEarthSand  = "bg-earth-sand"
	ClassEarthOchre = "bg-earth-ochre"

	// Fragment Paths
	PathListingsFragment = "/listings/fragment"
	CountryUSA           = "USA"
)
