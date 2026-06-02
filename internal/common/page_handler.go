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
	MockCreateData      interface{}
	MockLoginPromptData interface{}
	MockFeedbackData    interface{}
	MockEditData        interface{}
	MockDetailData      interface{}
	MockProfileData     interface{}
	MockSavedIDs        map[string]bool
	Config              *config.Config
	MockUser            *domain.User
	module.BaseViewData
	MockListingFeatured domain.Listing
	MockListingJob      domain.Listing
	MockListingNormal   domain.Listing
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

	mockSavedIDs := map[string]bool{
		"mock-listing-featured": true,
		"mock-listing-normal":   false,
	}

	baseViewData := h.PopulateBase(c)
	baseViewData.User = mockUser

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

	vm := SandboxViewModel{
		BaseViewData:        baseViewData,
		Config:              h.App.Cfg,
		MockListingNormal:   mockListingNormal,
		MockListingFeatured: mockListingFeatured,
		MockListingJob:      mockListingJob,
		MockUser:            mockUser,
		MockSavedIDs:        mockSavedIDs,
		MockDetailData:      mockDetailData,
		MockProfileData:     mockProfileData,
		MockCreateData:      baseViewData,
		MockEditData:        mockEditData,
		MockFeedbackData:    baseViewData,
		MockLoginPromptData: baseViewData,
	}

	return h.RenderTyped(c, "sandbox.html", vm)
}
