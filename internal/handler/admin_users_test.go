package handler_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminHandler_HandleUsers_Success(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/users", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetAllUsers", testifyMock.Anything, 50, 0).Return([]domain.User{{ID: "u1"}}, nil)

	h := handler.NewAdminHandler(mockRepo, nil, config.LoadConfig())
	_ = h.HandleUsers(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}
