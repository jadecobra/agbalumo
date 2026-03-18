package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreate_EventParsing(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	listingSvc := listmod.NewListingService(repo, repo, repo)
	h := listmod.NewListingHandler(repo, repo, listingSvc, nil, &MockGeocodingService{}, &config.Config{})

	// Create form data
	form := url.Values{}
	form.Set("title", "Test Event")
	form.Set("type", "Event")
	form.Set("owner_origin", "Nigeria")
	form.Set("contact_email", "test@example.com")
	form.Set("city", "Lagos")
	form.Set("event_start", "2026-12-25T10:00") // standard datetime-local format
	form.Set("event_end", "2026-12-25T14:00")

	c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(form.Encode()))
	c.Set("User", domain.User{ID: "event-user", Email: "event@example.com"})

	// Execute
	err := h.HandleCreate(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify DB state
	listings, _ := repo.FindByTitle(context.Background(), "Test Event")
	assert.Len(t, listings, 1)
	l := listings[0]
	assert.Equal(t, domain.Event, l.Type)

	expectedStart, _ := time.Parse("2006-01-02T15:04", "2026-12-25T10:00")
	expectedEnd, _ := time.Parse("2006-01-02T15:04", "2026-12-25T14:00")

	assert.True(t, l.EventStart.Equal(expectedStart))
	assert.True(t, l.EventEnd.Equal(expectedEnd))
}
