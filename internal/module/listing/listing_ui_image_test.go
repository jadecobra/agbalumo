package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEditListingModal_ImageRemovalUI(t *testing.T) {
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplate(t)}

	listing := domain.Listing{
		ID:       "test-ui-listing",
		Title:    "Test UI",
		ImageURL: "/static/uploads/test.webp",
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	data := map[string]interface{}{
		"Listing": listing,
		"User":    domain.User{ID: "u1"},
	}

	// Render the partial
	err := e.Renderer.Render(rec, "modal_edit_listing", data, c)
	assert.NoError(t, err)

	body := rec.Body.String()

	// 1. Verify hidden remove_image input is present with correct ID
	assert.Contains(t, body, `input type="hidden" name="remove_image" id="edit-remove-image-test-ui-listing" value="false"`)

	// 2. Verify "Remove" button (close icon) is present and calls clearEditImage
	assert.Contains(t, body, `hx-on:click="clearEditImage('test-ui-listing')"`)
	assert.Contains(t, body, `class="absolute -top-2 -right-2 bg-earth-ochre text-earth-dark p-1 hover:bg-earth-ochre-light shadow-lg"`)

	// 3. Verify the existing image div has the correct ID
	assert.Contains(t, body, `id="edit-existing-image-test-ui-listing"`)
}
