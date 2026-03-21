package admin_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_HandleToggleFeatured_Error(t *testing.T) {
	mockRepo := NewMockRepository()
	mockRepo.ErrorOn["SetFeatured"] = fmt.Errorf("db error")

	formData := url.Values{}
	formData.Set("featured", "true")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings/123/featured", strings.NewReader(formData.Encode()))
	c.SetParamNames("id")
	c.SetParamValues("123")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: mockRepo, FeedbackStore: mockRepo, AnalyticsStore: mockRepo, CategoryStore: mockRepo, UserStore: mockRepo, ListingStore: mockRepo, ClaimRequestStore: mockRepo, CSVService: nil, Cfg: config.LoadConfig()})
	err := h.HandleToggleFeatured(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_HandleToggleFeatured_BadRequest_MissingID(t *testing.T) {
	mockRepo := NewMockRepository()

	c, rec := setupAdminTestContext(http.MethodPost, "/admin/listings//featured", nil)
	// Missing ID param
	c.SetParamNames("id")
	c.SetParamValues("")
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	h := admin.NewAdminHandler(admin.AdminDependencies{AdminStore: mockRepo, FeedbackStore: mockRepo, AnalyticsStore: mockRepo, CategoryStore: mockRepo, UserStore: mockRepo, ListingStore: mockRepo, ClaimRequestStore: mockRepo, CSVService: nil, Cfg: config.LoadConfig()})
	
	err := h.HandleToggleFeatured(c)
	assert.NoError(t, err) // Echo handlers return nil and specify code in c.JSON
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp handler.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v. Body tracking: %s", err, rec.Body.String())
	}

	assert.Equal(t, "Listing ID is required", errResp.Error)
	assert.Equal(t, http.StatusBadRequest, errResp.Code) // Fails when current implementation just returns a map without code
}
