package handler_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleAddCategory_Success(t *testing.T) {
	formData := url.Values{}
	formData.Set("name", "Music")
	c, rec := setupAdminTestContext(http.MethodPost, "/admin/categories/add", strings.NewReader(formData.Encode()))

	adminUser := domain.User{ID: "admin1", Role: domain.UserRoleAdmin}
	c.Set("User", adminUser)
	store := middleware.NewTestSessionStore()
	session, _ := store.Get(c.Request(), "auth_session")
	c.Set("session", session)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("SaveCategory", testifyMock.Anything, testifyMock.Anything).Return(nil)

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	_ = h.HandleAddCategory(c)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "/admin", rec.Header().Get("Location"))
}
