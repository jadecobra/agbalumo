package handler_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_HandleUsers_Success(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/users", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	repo := handler.SetupTestRepository(t)
	// Seed a user
	user := domain.User{ID: "u1", Name: "Test User", Email: "test@test.com", Role: domain.UserRoleUser}
	err := repo.SaveUser(c.Request().Context(), user)
	require.NoError(t, err)

	h := handler.NewAdminHandler(repo, nil, config.LoadConfig())
	_ = h.HandleUsers(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}
