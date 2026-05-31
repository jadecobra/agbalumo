package ui_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/module/feedback"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/stretchr/testify/require"
)

var mockUser = domain.User{
	ID:        "user-1",
	Name:      "Test User",
	Email:     "test@example.com",
	AvatarURL: "https://example.com/avatar.png",
	Role:      domain.UserRoleUser,
}

var baseData = module.BaseViewData{
	Env:  "test",
	User: &mockUser,
	CSRF: "test-csrf",
}

var mockListing = domain.Listing{
	ID:        "listing-1",
	Title:     "Test Listing",
	CreatedAt: time.Now(),
	Status:    domain.ListingStatusApproved,
}

// 1. The Global Mock Registry
// Every template defined in the system MUST have an entry here with strictly typed mock data.
var TemplateMockRegistry = map[string]interface{}{
	// Pages
	"index.html": listing.HomeViewModel{
		BaseViewData: baseData,
		Listings:     []domain.Listing{mockListing},
		Pagination:   domain.Pagination{Page: 1, TotalPages: 1},
	},
	"profile.html": listing.ProfileViewModel{
		BaseViewData: baseData,
		Listings:     []domain.Listing{mockListing},
		SavedIDs:     map[string]bool{},
	},
	"modal_profile": listing.ProfileViewModel{
		BaseViewData: baseData,
		Listings:     []domain.Listing{mockListing},
		SavedIDs:     map[string]bool{},
	},
	"admin_dashboard.html": admin.AdminDashboardViewModel{
		BaseViewData: baseData,
		UserCount:    10,
	},
	"admin_listings.html": admin.AdminListingsViewModel{
		BaseViewData: baseData,
		Listings:     []domain.Listing{mockListing},
		Pagination:   domain.Pagination{Page: 1, TotalPages: 1},
	},
	"admin_users.html": admin.AdminUsersViewModel{
		BaseViewData: baseData,
		Users:        []domain.User{mockUser},
		Pagination:   domain.Pagination{Page: 1, TotalPages: 1},
	},
	"admin_login.html": map[string]interface{}{"Error": "", "CSRF": "test-csrf"},
	"about.html":       baseData,
	"error.html":       map[string]interface{}{"Message": "Test Error"},
	"admin_delete_confirm.html": admin.AdminDeleteViewModel{
		BaseViewData: baseData,
		IDs:          []string{"1", "2"},
	},
	"sandbox.html": common.SandboxViewModel{
		BaseViewData: baseData,
	},

	// Components/Partials (Keys match defined name OR filename without .html)
	"head_meta": baseData,
	"admin_modal_charts.html": admin.AdminDashboardViewModel{
		ListingGrowth: []domain.DailyMetric{{Date: "2024-01-01", Count: 10}},
		UserGrowth:    []domain.DailyMetric{{Date: "2024-01-01", Count: 5}},
	},
	"home_listings_section": listing.HomeViewModel{
		Listings: []domain.Listing{mockListing},
	},
	"admin_modal_moderation.html": map[string]interface{}{
		"ClaimRequests": []domain.ClaimRequest{},
	},
	"admin_listing_bulk_actions": admin.AdminListingsViewModel{
		BaseViewData: baseData,
	},
	"admin_listing_filters": admin.AdminListingsViewModel{
		BaseViewData: baseData,
	},
	"admin_modal_category.html": map[string]interface{}{
		"CSRF":       "test-csrf",
		"Categories": []domain.CategoryData{},
	},
	"listing_form_job_fields": map[string]interface{}{
		"IDPrefix": "test-",
		"Visible":  true,
		"Listing":  mockListing,
	},
	"admin_modal_bulk.html": map[string]interface{}{
		"CSRF": "test-csrf",
	},
	"admin_listing_table_header": admin.AdminListingsViewModel{
		BaseViewData: baseData,
		Pagination:   domain.Pagination{Page: 1},
		SortField:    "title",
		SortOrder:    "ASC",
		Category:     "Food",
	},
	"admin_listing_table_row": mockListing,
	"listing_form_type_origin": map[string]interface{}{
		"SelectedType":   "Business",
		"SelectedOrigin": "Nigeria",
		"Categories":     []domain.CategoryData{{Name: "Business"}},
	},
	"home_hero_search": listing.HomeViewModel{
		BaseViewData: module.BaseViewData{Categories: []domain.CategoryData{{Name: "Food"}}},
	},
	"admin_metrics_banner": admin.AdminDashboardViewModel{},
	"custom_country_options": map[string]interface{}{
		"Selected": "Nigeria",
	},
	"admin_modal_users.html": map[string]interface{}{"Users": []domain.User{mockUser}},
	"admin_feedback_item":    domain.Feedback{CreatedAt: time.Now()},
	"listing_form_event_fields": map[string]interface{}{
		"IDPrefix": "test-",
		"Visible":  true,
		"Listing":  mockListing,
	},
	"listing_form_image_fields": map[string]interface{}{
		"IDPrefix":  "test-",
		"ListingID": "1",
		"ImageURL":  "http://example.com/img.png",
	},
	"footer":     baseData,
	"mobile_nav": baseData,
	"listing_form_contact_fields": map[string]interface{}{
		"Listing":       mockListing,
		"User":          mockUser,
		"SplitWhatsApp": false,
	},
	"listing_form_common_fields": map[string]interface{}{
		"Value":            "Test",
		"IDPrefix":         "test-",
		"ListingID":        "1",
		"Address":          "123 Main",
		"GoogleMapsApiKey": "key",
		"City":             "Lagos",
		"Hours":            "9-5",
	},
	"listing_form_title":       map[string]interface{}{"Value": "Test"},
	"listing_form_description": map[string]interface{}{"Value": "Test", "Type": ""},
	"listing_form_ada_signals": mockListing,
	"listing_form_location": map[string]interface{}{
		"IDPrefix":         "test-",
		"ListingID":        "1",
		"Address":          "123 Main",
		"City":             "Lagos",
		"Hours":            "9-5",
		"GoogleMapsApiKey": "key",
	},
	"listing_form_website": map[string]interface{}{"Value": "http://example.com"},
	"admin_tools_grid":     admin.AdminDashboardViewModel{},
	"header_search":        baseData,
	"navigation":           baseData,
	"admin_header_content": admin.AdminDashboardViewModel{},
	"pagination": map[string]interface{}{
		"Pagination": domain.Pagination{Page: 1, TotalPages: 5},
		"FilterType": "Food",
		"City":       "Lagos",
		"Query":      "Suya",
	},
	"admin_pagination.html": map[string]interface{}{
		"Pagination": domain.Pagination{Page: 1, TotalPages: 5},
		"Category":   "",
		"QueryText":  "",
		"SortField":  "",
		"SortOrder":  "",
	},
	"save_button": map[string]interface{}{"ListingID": "1", "IsSaved": true, "Classes": "", "TextColorClass": "", "IDPrefix": "", "OOB": false},
	"modal_feedback": feedback.FeedbackViewModel{
		BaseViewData: baseData,
		Path:         "/test",
	},
	"modal_feedback_content": feedback.FeedbackViewModel{
		BaseViewData: baseData,
		Path:         "/test",
	},
	"featured_card": map[string]interface{}{"Listing": mockListing},
	"ui_components": map[string]interface{}{},
	"modal_create_listing": map[string]interface{}{
		"CSRF":             "test-csrf",
		"Categories":       []domain.CategoryData{{Name: "Food"}},
		"GoogleMapsApiKey": "key",
		"User":             mockUser,
	},
	"modal_edit_listing": listing.EditViewModel{
		BaseViewData: baseData,
		Listing:      mockListing,
	},
	"modal_login_prompt.html": map[string]interface{}{},
	"modal_detail": listing.DetailViewModel{
		Listing:  mockListing,
		CanClaim: true,
		SavedIDs: map[string]bool{},
	},
	"modal_detail_content": listing.DetailViewModel{
		Listing:  mockListing,
		CanClaim: true,
		SavedIDs: map[string]bool{},
	},
	"toast_error.html": map[string]interface{}{
		"Message": "Test Error",
		"ID":      "toast-1",
		"Title":   "Error",
	},
	"listing_list": listing.ListingFragmentViewModel{
		Listings:   []domain.Listing{mockListing},
		Pagination: domain.Pagination{Page: 1, TotalPages: 1},
	},
	"listing_card": map[string]interface{}{
		"Listing":   mockListing,
		"SavedIDs":  map[string]bool{},
		"User":      mockUser,
		"IDPrefix":  "",
		"GridClass": "",
	},
	"rating_stars": map[string]interface{}{
		"Listing":        mockListing,
		"TextColorClass": "text-earth-ochre dark:text-yellow-600",
	},

	"button_sharp": map[string]interface{}{
		"Label": "Test",
	},
	"admin_tool_link_sharp": map[string]interface{}{
		"Link":           "/",
		"IconBgClass":    "bg-red-500",
		"IconColorClass": "text-white",
		"Icon":           "icon",
		"Label":          "Test",
	},
	"metric_stat_sharp": map[string]interface{}{
		"Value": "10",
		"Label": "Test",
	},
	"status_badge_sharp": map[string]interface{}{
		"Label": "Approved",
	},
	"admin_tool_btn_sharp": map[string]interface{}{
		"IconBgClass":    "bg-red-500",
		"IconColorClass": "text-white",
		"Icon":           "icon",
		"Label":          "Test",
	},
	"btn_close": map[string]interface{}{
		"Classes": "",
		"Action":  "close",
		"TestID":  "test-id",
	},
	"modal_base": map[string]interface{}{
		"ID":    "test",
		"Title": "Test Modal",
	},
	"modal_create_listing_fields": baseData,
	"modal_profile_content": listing.ProfileViewModel{
		BaseViewData: baseData,
		Listings:     []domain.Listing{mockListing},
		SavedIDs:     map[string]bool{},
	},
	"modal_edit_listing_fields": listing.EditViewModel{
		Listing: mockListing,
	},
	"modal_login_prompt_content": baseData,
	"location_permission_prompt": map[string]interface{}{},
}

func TestGlobalTemplateCoverage(t *testing.T) {
	// Initialize the true engine (use relative path to templates)
	renderer, err := ui.NewTemplateRenderer("../../ui/templates/*.html", "../../ui/templates/*/*.html")
	require.NoError(t, err)

	// 2. Discover all templates natively
	names := renderer.GetAllTemplateNames()

	// 3. Assert 100% Coverage & Contract Validity
	for _, name := range names {
		// Skip common block names that are redefined across pages.
		// These are tested when the page templates themselves are rendered.
		if name == "content" || name == "title" || name == "header_content" ||
			name == "header_classes" || name == "extra_scripts" || name == "filters" ||
			name == "inner_content" || name == "base.html" || name == "base" || name == "" {
			continue
		}

		// Check if it exists in registry
		mockData, exists := TemplateMockRegistry[name]
		if !exists {
			// Try without .html extension
			if strings.HasSuffix(name, ".html") {
				mockData, exists = TemplateMockRegistry[strings.TrimSuffix(name, ".html")]
			}
		}

		require.True(t, exists, "Template '%s' is missing from the mock registry! Add it to TemplateMockRegistry to verify its data contract.", name)

		// Render to verify contract
		err := renderer.RenderDefinition(io.Discard, name, mockData)
		require.NoError(t, err, "Template '%s' failed to render. Missing key or contract violation.", name)
	}
}
