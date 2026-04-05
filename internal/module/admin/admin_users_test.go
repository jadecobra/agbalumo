package admin_test

import (
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module/admin"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_HandleUsers_Success(t *testing.T) {
	c, rec := setupAdminTestContext(http.MethodGet, "/admin/users", nil)
	c.Set("User", domain.User{Role: domain.UserRoleAdmin})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	// Seed a user
	user := domain.User{ID: "u1", Name: "Test User", Email: "test@test.com", Role: domain.UserRoleUser}
	err := app.DB.SaveUser(c.Request().Context(), user)
	require.NoError(t, err)

	h := admin.NewAdminHandler(app)
	_ = h.HandleUsers(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}
