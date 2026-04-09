package admin_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleToggleFeatured_Error(t *testing.T) {
	t.Parallel()
	_, h, mockRepo := setupAdminMockTest(t)
	mockRepo.ErrorOn["SetFeatured"] = fmt.Errorf("db error")

	formData := url.Values{}
	formData.Set("featured", "true")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/123/featured", strings.NewReader(formData.Encode()))
	setupAdminAuth(t, c)
	c.SetParamNames("id")
	c.SetParamValues("123")

	err := h.HandleToggleFeatured(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleToggleFeatured_BadRequest_MissingID(t *testing.T) {
	t.Parallel()
	_, h, _ := setupAdminMockTest(t)

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings//featured", nil)
	// Missing ID param
	c.SetParamNames("id")
	c.SetParamValues("")
	setupAdminAuth(t, c)

	err := h.HandleToggleFeatured(c)
	assert.NoError(t, err) // Echo handlers return nil and specify code in c.JSON
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ui.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v. Body tracking: %s", err, rec.Body.String())
	}

	assert.Equal(t, "Listing ID is required", errResp.Error)
	assert.Equal(t, http.StatusBadRequest, errResp.Code) // Fails when current implementation just returns a map without code
}
