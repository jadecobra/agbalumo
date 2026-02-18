package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

type MockRenderer struct{}

func (m *MockRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return nil
}

func TestHandleCreate_EventParsing(t *testing.T) {
	e := echo.New()
	e.Renderer = &MockRenderer{}

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		// Asset that the parsed listing has the correct dates
		if l.Type != domain.Event {
			return false
		}

		// Expected parsed time
		expectedStart, _ := time.Parse("2006-01-02T15:04", "2026-12-25T10:00")
		expectedEnd, _ := time.Parse("2006-01-02T15:04", "2026-12-25T14:00")

		// Check within duration manually or just trust assert inside match if meaningful?
		// MatchedBy expects a bool. We can't use t.Errorf here easily.
		// So we return true/false.
		// To be precise, let's just use strict equality or approximate if needed.
		// Given it's time parsing, exact match of Unix or Equal is good.
		if !l.EventStart.Equal(expectedStart) {
			return false
		}
		if !l.EventEnd.Equal(expectedEnd) {
			return false
		}
		return true
	})).Return(nil)

	h := handler.NewListingHandler(mockRepo, nil)

	// Create form data
	form := url.Values{}
	form.Set("title", "Test Event")
	form.Set("type", "Event")
	form.Set("owner_origin", "Nigeria")
	form.Set("contact_email", "test@example.com")
	form.Set("event_start", "2026-12-25T10:00") // standard datetime-local format
	form.Set("event_end", "2026-12-25T14:00")

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.Set("User", domain.User{ID: "event-user", Email: "event@example.com"})

	// Execute
	err := h.HandleCreate(c)
	assert.NoError(t, err)
	if rec.Code != http.StatusOK {
		t.Logf("Response Body: %s", rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}
