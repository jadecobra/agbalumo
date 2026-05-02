package listing

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildListingViewData(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	h := NewListingHandler(app)
	ctx := context.Background()

	listings := []domain.Listing{
		{ID: "l1", Title: "Listing 1"},
		{ID: "l2", Title: "Listing 2"},
	}

	t.Run("Authenticated User", func(t *testing.T) {
		u := &domain.User{ID: "user-1", Email: "test@example.com"}
		err := app.DB.SaveUser(ctx, *u)
		require.NoError(t, err)

		// Save one listing
		err = app.DB.SaveListing(ctx, u.ID, "l1")
		require.NoError(t, err)

		c, _ := testutil.SetupModuleContext(http.MethodGet, "/", nil)
		c.Set(domain.CtxKeyUser, u)

		data := h.buildListingViewData(c, listings)

		assert.Equal(t, listings, data["Listings"])
		assert.Equal(t, u, data["User"])
		assert.NotNil(t, data["SavedIDs"])
		assert.Contains(t, data["SavedIDs"], "l1")
		assert.NotContains(t, data["SavedIDs"], "l2")
	})

	t.Run("Unauthenticated User", func(t *testing.T) {
		c, _ := testutil.SetupModuleContext(http.MethodGet, "/", nil)
		// No user set in context

		data := h.buildListingViewData(c, listings)

		assert.Equal(t, listings, data["Listings"])
		assert.Nil(t, data["User"])
		assert.Nil(t, data["SavedIDs"])
	})
}
