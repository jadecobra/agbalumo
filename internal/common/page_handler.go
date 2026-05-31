package common

import (
	"net/http"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/labstack/echo/v4"
)

// PageHandler manages rendering for generic pages like About.
type PageHandler struct {
	module.BaseHandler
}

// NewPageHandler creates a new PageHandler.
func NewPageHandler(app *env.AppEnv) *PageHandler {
	return &PageHandler{
		BaseHandler: module.BaseHandler{App: app},
	}
}

// AboutViewModel represents the strongly-typed data passed to the about page template.
type AboutViewModel struct {
	Config *config.Config
	module.BaseViewData
}

// HandleAbout renders the generic about page.
func (h *PageHandler) HandleAbout(c echo.Context) error {
	vm := AboutViewModel{
		BaseViewData: h.PopulateBase(c),
		Config:       h.App.Cfg,
	}

	return h.RenderTyped(c, "about.html", vm)
}

// SandboxViewModel represents the strongly-typed data passed to the sandbox page template.
type SandboxViewModel struct {
	Config       *config.Config
	MockUser     *domain.User
	MockSavedIDs map[string]bool
	module.BaseViewData
	MockListingNormal     domain.Listing
	MockListingFeatured   domain.Listing
	MockListingJob        domain.Listing
	MockListingEvent      domain.Listing
	MockListingBusiness   domain.Listing
	MockListingZeroRating domain.Listing
	MockListingNotOwned   domain.Listing
	MockPagination1       domain.Pagination
	MockPagination2       domain.Pagination
	MockDetailData        interface{}
	MockProfileData       interface{}
	MockCreateData        interface{}
	MockEditData          interface{}
	MockFeedbackData      interface{}
	MockLoginPromptData   interface{}
	MockCategories        []domain.CategoryData
	MockCounts            map[string]int
	MockAdminListings     []domain.Listing
	MockAdminFeedback     domain.Feedback
	MockAdminMetrics      interface{}
	MockAdminUsers        []domain.User
	MockNilUser           *domain.User
	MockToastID           string
	MockToastTitle        string
	MockToastMessage      string
}

// HandleSandbox renders the component sandbox page for non-production environments.
func (h *PageHandler) HandleSandbox(c echo.Context) error {
	if h.App.Cfg.Env == domain.EnvProduction {
		return echo.NewHTTPError(http.StatusNotFound, "Sandbox not available in production")
	}

	mockUser := &domain.User{
		ID:        "mock-user-123",
		Email:     "owner@agbalumo.com",
		Name:      "Chidi Egwu",
		AvatarURL: "",
		Role:      domain.UserRoleUser,
	}

	mockListingNormal := domain.Listing{
		ID:                "mock-listing-normal",
		Title:             "Brutalist Jollof Joint",
		Description:       "Authentic West African flavors served in a premium, high-contrast environment.",
		Type:              domain.Food,
		City:              "Lagos",
		Country:           "Nigeria",
		Rating:            3.5,
		ReviewCount:       12,
		OwnerOrigin:       "Nigeria",
		RegionalSpecialty: "Jollof Rice",
		WebsiteURL:        "https://jollof.example.com",
		OwnerID:           "mock-user-123",
		CreatedAt:         time.Now().Add(-2 * 24 * time.Hour), // Not new (older than 24h)
	}

	mockListingFeatured := domain.Listing{
		ID:                "mock-listing-featured",
		Title:             "Premium Suya & Grill",
		Description:       "Spicy and smokey suya, grilled to perfection using authentic family recipes.",
		Type:              domain.Food,
		City:              "London",
		Country:           "United Kingdom",
		Rating:            5.0,
		ReviewCount:       42,
		OwnerOrigin:       "Nigeria",
		RegionalSpecialty: "Suya",
		WebsiteURL:        "https://suya.example.com",
		Featured:          true,
		OwnerID:           "mock-user-123",
		CreatedAt:         time.Now(), // New (created just now)
	}

	mockListingJob := domain.Listing{
		ID:           "mock-listing-job",
		Title:        "Senior Suya Chef",
		Description:  "We are seeking an experienced chef specializing in authentic Nigerian Suya and grilling techniques.",
		Type:         domain.Job,
		City:         "London",
		Country:      "United Kingdom",
		Company:      "Suya Grill Palace Ltd",
		PayRange:     "£35,000 - £45,000 / year",
		Skills:       "Grilling, Meat Prep, Spice Blending, Food Safety",
		JobApplyURL:  "https://suyapalace.example.com/apply",
		JobStartDate: time.Now().Add(14 * 24 * time.Hour),
		OwnerID:      "mock-user-123",
		CreatedAt:    time.Now().Add(-5 * time.Hour),
	}

	mockListingEvent := domain.Listing{
		ID:          "mock-listing-event",
		Title:       "West African Street Food Fest",
		Description: "Join us for a vibrant weekend of music, culture, and tasting authentic street food from West Africa.",
		Type:        domain.Event,
		City:        "Lagos",
		Country:     "Nigeria",
		EventStart:  time.Now().Add(48 * time.Hour),
		EventEnd:    time.Now().Add(54 * time.Hour),
		OwnerID:     "mock-user-123",
		CreatedAt:   time.Now().Add(-12 * time.Hour),
	}

	mockListingBusiness := domain.Listing{
		ID:               "mock-listing-business",
		Title:            "Onyx Shea Butter Distributors",
		Description:      "Wholesale and retail suppliers of 100% organic, raw unrefined shea butter directly from Ghana.",
		Type:             domain.Business,
		City:             "Accra",
		Country:          "Ghana",
		HoursOfOperation: "9:00 AM - 5:00 PM (GMT)",
		OwnerID:          "mock-user-123",
		CreatedAt:        time.Now().Add(-10 * 24 * time.Hour),
	}

	mockListingZeroRating := domain.Listing{
		ID:                "mock-listing-zero",
		Title:             "New Age Chin Chin Bakery",
		Description:       "Handmade gourmet chin chin in various flavors: vanilla, cinnamon, and spicy ginger.",
		Type:              domain.Food,
		City:              "Toronto",
		Country:           "Canada",
		Rating:            0.0,
		ReviewCount:       0,
		OwnerOrigin:       "Nigeria",
		RegionalSpecialty: "Chin Chin",
		OwnerID:           "mock-user-123",
		CreatedAt:         time.Now().Add(-30 * time.Minute),
	}

	mockListingNotOwned := domain.Listing{
		ID:                "mock-listing-not-owned",
		Title:             "Mama Put Restaurant",
		Description:       "Traditional home-style African dishes. Fufu, Egusi, Okra soup and much more.",
		Type:              domain.Food,
		City:              "Houston",
		Country:           "United States",
		Rating:            4.8,
		ReviewCount:       124,
		OwnerOrigin:       "Nigeria",
		RegionalSpecialty: "Egusi Soup",
		OwnerID:           "other-owner-789", // Not mock-user-123
		CreatedAt:         time.Now().Add(-90 * 24 * time.Hour),
	}

	mockSavedIDs := map[string]bool{
		"mock-listing-featured": true,
		"mock-listing-normal":   false,
	}

	mockPagination1 := domain.Pagination{
		Page:        1,
		Limit:       10,
		Offset:      0,
		TotalCount:  25,
		TotalPages:  3,
		HasNextPage: true,
	}

	mockPagination2 := domain.Pagination{
		Page:        2,
		Limit:       10,
		Offset:      10,
		TotalCount:  25,
		TotalPages:  3,
		HasNextPage: true,
	}

	baseViewData := h.PopulateBase(c)

	mockDetailData := struct {
		GoogleMapsApiKey string
		Category         domain.CategoryData
		SavedIDs         map[string]bool
		module.BaseViewData
		Listing  domain.Listing
		CanClaim bool
	}{
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
		Category: domain.CategoryData{
			Name:      "Food",
			Claimable: true,
		},
		SavedIDs:     mockSavedIDs,
		BaseViewData: baseViewData,
		Listing:      mockListingFeatured,
		CanClaim:     true,
	}

	mockProfileData := struct {
		SavedIDs         map[string]bool
		GoogleMapsApiKey string
		Listings         []domain.Listing
		SavedListings    []domain.Listing
		module.BaseViewData
	}{
		SavedIDs:         mockSavedIDs,
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
		Listings:         []domain.Listing{mockListingNormal, mockListingFeatured, mockListingJob},
		SavedListings:    []domain.Listing{mockListingFeatured},
		BaseViewData:     baseViewData,
	}

	mockEditData := struct {
		TargetID         string
		Source           string
		GoogleMapsApiKey string
		module.BaseViewData
		Listing domain.Listing
	}{
		TargetID:         "edit-listing-featured",
		Source:           "sandbox",
		GoogleMapsApiKey: h.App.Cfg.GoogleMapsAPIKey,
		Listing:          mockListingFeatured,
		BaseViewData:     baseViewData,
	}

	mockFeedbackData := struct {
		Path string
		module.BaseViewData
	}{
		Path:         "/sandbox",
		BaseViewData: baseViewData,
	}

	mockCategories := []domain.CategoryData{
		{Name: "Food"},
		{Name: "Job"},
		{Name: "Event"},
		{Name: "Business"},
	}

	mockCounts := map[string]int{
		"Food":     24,
		"Job":      5,
		"Event":    3,
		"Business": 12,
	}

	mockAdminListings := []domain.Listing{mockListingNormal, mockListingFeatured, mockListingJob, mockListingEvent, mockListingBusiness}
	mockAdminFeedback := domain.Feedback{
		ID:        "fb-1",
		UserID:    mockUser.ID,
		Type:      domain.FeedbackTypeFeature,
		Content:   "I love the brand color schema and the premium brutalist typography. Latency is extremely fast!",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}

	mockAdminMetrics := struct {
		ListingGrowth []domain.DailyMetric
		UserGrowth    []domain.DailyMetric
		module.BaseViewData
		AdaDiscoveryAvg float64
		UserCount       int
		ListingCount    int
	}{
		UserCount:       128,
		ListingCount:    49,
		AdaDiscoveryAvg: 4.2,
		ListingGrowth: []domain.DailyMetric{
			{Date: "05-25", Count: 3},
			{Date: "05-26", Count: 5},
			{Date: "05-27", Count: 8},
			{Date: "05-28", Count: 12},
			{Date: "05-29", Count: 19},
			{Date: "05-30", Count: 31},
			{Date: "05-31", Count: 49},
		},
		UserGrowth: []domain.DailyMetric{
			{Date: "05-25", Count: 10},
			{Date: "05-26", Count: 25},
			{Date: "05-27", Count: 40},
			{Date: "05-28", Count: 62},
			{Date: "05-29", Count: 85},
			{Date: "05-30", Count: 108},
			{Date: "05-31", Count: 128},
		},
		BaseViewData: baseViewData,
	}

	mockAdminUsers := []domain.User{
		*mockUser,
		{ID: "user-2", Email: "african.gourmet@gmail.com", Name: "Amina Bello", Role: domain.UserRoleUser},
		{ID: "user-3", Email: "kwame.delivery@yahoo.com", Name: "Kwame Mensah", Role: domain.UserRoleUser},
	}

	vm := SandboxViewModel{
		BaseViewData:          baseViewData,
		Config:                h.App.Cfg,
		MockListingNormal:     mockListingNormal,
		MockListingFeatured:   mockListingFeatured,
		MockListingJob:        mockListingJob,
		MockListingEvent:      mockListingEvent,
		MockListingBusiness:   mockListingBusiness,
		MockListingZeroRating: mockListingZeroRating,
		MockListingNotOwned:   mockListingNotOwned,
		MockUser:              mockUser,
		MockSavedIDs:          mockSavedIDs,
		MockPagination1:       mockPagination1,
		MockPagination2:       mockPagination2,
		MockDetailData:        mockDetailData,
		MockProfileData:       mockProfileData,
		MockCreateData:        baseViewData,
		MockEditData:          mockEditData,
		MockFeedbackData:      mockFeedbackData,
		MockLoginPromptData:   baseViewData,
		MockCategories:        mockCategories,
		MockCounts:            mockCounts,
		MockAdminListings:     mockAdminListings,
		MockAdminFeedback:     mockAdminFeedback,
		MockAdminMetrics:      mockAdminMetrics,
		MockAdminUsers:        mockAdminUsers,
		MockNilUser:           nil,
		MockToastID:           "toast-sandbox-error",
		MockToastTitle:        "Error Occurred",
		MockToastMessage:      "Failed to load listing details. Please try again.",
	}

	return h.RenderTyped(c, "sandbox.html", vm)
}
