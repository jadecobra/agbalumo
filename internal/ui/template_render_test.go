// template_render_test.go tests that templates render correct HTML output
// when given specific data. These tests replace expensive browser verification
// for component-level changes. Run: go test ./internal/ui/... -run TestRender
package ui_test

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRender_SaveButton_Saved(t *testing.T) {
	out := renderPartial(t, "save_button", map[string]interface{}{
		"ListingID": "abc-123",
		"IsSaved":   true,
	})
	assert.Contains(t, out, `id="save-btn-abc-123"`)
	assert.Contains(t, out, "favorite")           // filled heart icon
	assert.Contains(t, out, "text-red-500")       // saved state color
	assert.NotContains(t, out, "favorite_border") // NOT outline
}

func TestRender_SaveButton_Unsaved(t *testing.T) {
	out := renderPartial(t, "save_button", map[string]interface{}{
		"ListingID": "abc-123",
		"IsSaved":   false,
	})
	assert.Contains(t, out, "favorite_border") // outline heart
	assert.Contains(t, out, "text-stone-400")  // unsaved state color
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
