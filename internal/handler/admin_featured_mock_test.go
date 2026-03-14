package handler_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
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

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	err := h.HandleToggleFeatured(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
