package listing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSaveToggle_SaveAndUnsave(t *testing.T) {
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	h := NewListingHandler(app)
	e := echo.New()

	// 1. Create a user and listing
	u := &domain.User{ID: "user-1", Email: "test@example.com"}
	err := app.DB.SaveUser(context.Background(), *u)
	require.NoError(t, err)

	l := domain.Listing{ID: "listing-1", Title: "Test Listing", Type: domain.Food}
	err = app.DB.Save(context.Background(), l)
	require.NoError(t, err)

	// 2. First POST: Save listing
	req := httptest.NewRequest(http.MethodPost, "/listings/listing-1/save", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("listing-1")
	c.Set(domain.CtxKeyUser, u)

	err = h.HandleSaveToggle(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ListingID")
	assert.Contains(t, rec.Body.String(), "listing-1")

	// Verify in DB
	saved, err := app.DB.IsListingSaved(context.Background(), u.ID, l.ID)
	require.NoError(t, err)
	assert.True(t, saved)

	// 3. Second POST: Unsave listing
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("listing-1")
	c.Set(domain.CtxKeyUser, u)

	err = h.HandleSaveToggle(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify in DB
	saved, err = app.DB.IsListingSaved(context.Background(), u.ID, l.ID)
	require.NoError(t, err)
	assert.False(t, saved)
}
