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

// SandboxPagination replicates the pagination structure locally to avoid import cycles.
type SandboxPagination struct {
	Page        int
	Limit       int
	Offset      int
	TotalCount  int
	TotalPages  int
	HasNextPage bool
}

// GetPageRange returns a slice of page numbers for template display.
func (p SandboxPagination) GetPageRange() []int {
	if p.TotalPages <= 1 {
		return nil
	}
	start := p.Page - 2
	if start < 1 {
		start = 1
	}
	end := start + 4
	if end > p.TotalPages {
		end = p.TotalPages
		start = end - 4
		if start < 1 {
			start = 1
		}
	}
	var pages []int
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}

// SandboxViewModel represents the strongly-typed data passed to the sandbox page template.
type SandboxViewModel struct {
	Config       *config.Config
	MockUser     *domain.User
	MockSavedIDs map[string]bool
	module.BaseViewData
	MockListingNormal   domain.Listing
	MockListingFeatured domain.Listing
	MockPagination1     SandboxPagination
	MockPagination2     SandboxPagination
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

	mockSavedIDs := map[string]bool{
		"mock-listing-featured": true,
		"mock-listing-normal":   false,
	}

	mockPagination1 := SandboxPagination{
		Page:        1,
		Limit:       10,
		Offset:      0,
		TotalCount:  25,
		TotalPages:  3,
		HasNextPage: true,
	}

	mockPagination2 := SandboxPagination{
		Page:        2,
		Limit:       10,
		Offset:      10,
		TotalCount:  25,
		TotalPages:  3,
		HasNextPage: true,
	}

	vm := SandboxViewModel{
		BaseViewData:        h.PopulateBase(c),
		Config:              h.App.Cfg,
		MockListingNormal:   mockListingNormal,
		MockListingFeatured: mockListingFeatured,
		MockUser:            mockUser,
		MockSavedIDs:        mockSavedIDs,
		MockPagination1:     mockPagination1,
		MockPagination2:     mockPagination2,
	}

	return h.RenderTyped(c, "sandbox.html", vm)
}
