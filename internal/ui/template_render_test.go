// template_render_test.go tests that templates render correct HTML output
// when given specific data. These tests replace expensive browser verification
// for component-level changes. Run: go test ./internal/ui/... -run TestRender
package ui_test

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/stretchr/testify/assert"
)

func TestRender_SaveButton_Saved(t *testing.T) {
	out := renderPartial(t, "save_button", map[string]interface{}{
		"ListingID":      "abc-123",
		"IsSaved":        true,
		"Classes":        "",
		"TextColorClass": "",
		"IDPrefix":       "",
		"OOB":            false,
	})
	assert.Contains(t, out, `id="save-btn-abc-123"`)
	assert.Contains(t, out, "favorite")          // filled heart icon
	assert.Contains(t, out, "text-earth-accent") // saved state color
	assert.Contains(t, out, "[font-variation-settings:'FILL'_1]")
	assert.NotContains(t, out, "hover:[font-variation-settings:'FILL'_1]")
}

func TestRender_SaveButton_Unsaved(t *testing.T) {
	out := renderPartial(t, "save_button", map[string]interface{}{
		"ListingID":      "abc-123",
		"IsSaved":        false,
		"Classes":        "",
		"TextColorClass": "",
		"IDPrefix":       "",
		"OOB":            false,
	})
	assert.Contains(t, out, "favorite")
	assert.Contains(t, out, "hover:[font-variation-settings:'FILL'_1]")
	assert.Contains(t, out, "text-earth-clay/50") // unsaved state color
}

func TestRender_ListingCard_WithSavedIDs(t *testing.T) {
	out := renderPartial(t, "listing_card", map[string]interface{}{
		"Listing":   domain.Listing{ID: "abc-123", Title: "Test Food", Type: domain.Food},
		"User":      &domain.User{ID: "u1"},
		"SavedIDs":  map[string]bool{"abc-123": true},
		"Index":     0,
		"IDPrefix":  "",
		"GridClass": "",
	})
	assert.Contains(t, out, "save-btn-abc-123") // heart button rendered
	assert.Contains(t, out, "Test Food")        // title rendered
}

func TestRender_ListingCard_NoUser_NoHeart(t *testing.T) {
	out := renderPartial(t, "listing_card", map[string]interface{}{
		"Listing":   domain.Listing{ID: "abc-123", Title: "Test"},
		"Index":     0,
		"SavedIDs":  map[string]bool{},
		"User":      nil,
		"IDPrefix":  "",
		"GridClass": "",
	})
	assert.NotContains(t, out, "save-btn") // no heart for anonymous
}

func TestRender_ListingList_RedesignTokens(t *testing.T) {
	mockListing := domain.Listing{ID: "l-1", Title: "Jollof Rice Spot", Type: domain.Food}

	// 1. Featured Section should have radiant amber header text
	outFeatured := renderPartial(t, "listing_list", listing.ListingFragmentViewModel{
		Featured: []domain.Listing{mockListing},
	})
	assert.Contains(t, outFeatured, `text-earth-accent">Featured</h2>`)
	assert.Contains(t, outFeatured, "bg-earth-ochre/10")

	// 2. Fallback City container should use dark espresso and amber button
	outFallback := renderPartial(t, "listing_list", listing.ListingFragmentViewModel{
		FallbackCity: "Dallas",
		Listings:     []domain.Listing{mockListing},
		SavedIDs:     map[string]bool{},
	})
	assert.Contains(t, outFallback, "bg-earth-espresso/60")
	assert.Contains(t, outFallback, "text-white")
	assert.Contains(t, outFallback, "from-earth-accent to-earth-ochre")

	// 3. Location status indicator should have text-earth-accent
	outLocation := renderPartial(t, "listing_list", listing.ListingFragmentViewModel{
		Radius: 25.0,
		City:   "Dallas",
	})
	assert.Contains(t, outLocation, "text-earth-accent")
}

func TestRender_ModalDetail_RedesignTokens(t *testing.T) {
	mockListing := domain.Listing{
		ID:      "l-2",
		Title:   "Suya Spot",
		Type:    domain.Food,
		MenuURL: "https://example.com/menu",
		Rating:  4.5,
	}
	out := renderPartial(t, "modal_detail", listing.DetailViewModel{
		Listing:  mockListing,
		SavedIDs: map[string]bool{},
	})
	// Rating star should use text-earth-accent (parity with listing_card)
	assert.Contains(t, out, "text-earth-accent dark:text-yellow-600")
	// Menu CTA button should use radiant amber gradient
	assert.Contains(t, out, "from-earth-accent to-earth-ochre")
	assert.Contains(t, out, "text-earth-dark font-extrabold")
}

func TestRender_HomeListingsSection_LoadingSkeleton(t *testing.T) {
	out := renderPartial(t, "home_listings_section", listing.HomeViewModel{
		Listings: []domain.Listing{},
		SavedIDs: map[string]bool{},
	})
	// Skeleton loader should have espresso background and amber border
	assert.Contains(t, out, "bg-earth-espresso/60")
	assert.Contains(t, out, "border-earth-accent/20")
}
