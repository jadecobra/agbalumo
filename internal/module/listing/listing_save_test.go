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

func TestHandleSaveToggle_SaveAndUnsave(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	h := NewListingHandler(app)

	// 1. Create a user and listing
	u := &domain.User{ID: "user-1", Email: "test@example.com"}
	err := app.DB.SaveUser(context.Background(), *u)
	require.NoError(t, err)

	l := domain.Listing{ID: "listing-1", Title: "Test Listing", Type: domain.Food}
	err = app.DB.Save(context.Background(), l)
	require.NoError(t, err)

	// 2. First POST: Save listing
	c, rec := testutil.SetupModuleContext(http.MethodPost, "/listings/listing-1/save", nil)
	c.SetParamNames("id")
	c.SetParamValues("listing-1")
	c.Set(domain.CtxKeyUser, u)

	err = h.HandleSaveToggle(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "save-button-listing-1")
	assert.Contains(t, rec.Body.String(), "Saved")

	// Verify in DB
	saved, err := app.DB.IsListingSaved(context.Background(), u.ID, l.ID)
	require.NoError(t, err)
	assert.True(t, saved)

	// 3. Second POST: Unsave listing
	c2, rec2 := testutil.SetupModuleContext(http.MethodPost, "/listings/listing-1/save", nil)
	c2.SetParamNames("id")
	c2.SetParamValues("listing-1")
	c2.Set(domain.CtxKeyUser, u)

	err = h.HandleSaveToggle(c2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec2.Code)
	assert.Contains(t, rec2.Body.String(), "Save") // Not "Saved"

	// Verify in DB
	saved, err = app.DB.IsListingSaved(context.Background(), u.ID, l.ID)
	require.NoError(t, err)
	assert.False(t, saved)
}
