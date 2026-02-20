package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

// Data Validation Integration Tests
func TestIntegration_DataValidation(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	// 1. Positive Test Case: Known Good Data (Subset of seed data)
	goodData := []string{
		"title=Good+Biz&type=Business&owner_origin=Nigeria&description=Valid&contact_email=good@test.com&address=123+Main+St",
		"title=Good+Req&type=Request&owner_origin=Ghana&description=Valid&contact_whatsapp=+123456&deadline_date=" + time.Now().Add(24*time.Hour).Format("2006-01-02"),
	}

	for _, bodyStr := range goodData {
		mockRepo := &mock.MockListingRepository{}
		mockRepo.On("Save", testifyMock.Anything, testifyMock.Anything).Return(nil)

		h := handler.NewListingHandler(mockRepo, nil)

		req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(bodyStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("User", domain.User{ID: "test-user-id", Email: "good@test.com"})

		if err := h.HandleCreate(c); err != nil {
			t.Errorf("Unexpected error for good data: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for good data, got %d. Body: %s", rec.Code, rec.Body.String())
		}
		mockRepo.AssertExpectations(t)
	}

	// 2. Negative Test Cases: Known Bad Data
	badData := []struct {
		name      string
		body      string
		wantError string
	}{
		{
			name:      "Missing Origin",
			body:      "title=Bad+Biz&type=Business&description=No+Origin&contact_email=a@b.com",
			wantError: "owner origin is required",
		},
		{
			name:      "Invalid Origin",
			body:      "title=Bad+Biz&type=Business&owner_origin=Mars&description=Alien&contact_email=a@b.com",
			wantError: "owner origin must be a West African country",
		},
		{
			name:      "Missing Contact",
			body:      "title=Bad+Service&type=Service&owner_origin=Nigeria&description=Ghost",
			wantError: "at least one contact method is required",
		},
		{
			name:      "Request Deadline in Past",
			body:      "title=Bad+Req&type=Request&owner_origin=Ghana&description=Late&contact_email=a@b.com&deadline_date=" + time.Now().Add(-48*time.Hour).Format("2006-01-02"),
			wantError: "deadline cannot be in the past",
		},
		{
			name:      "Request Deadline Too Far (>90 days)",
			body:      "title=Bad+Req&type=Request&owner_origin=Ghana&description=Future&contact_email=a@b.com&deadline_date=" + time.Now().Add(100*24*time.Hour).Format("2006-01-02"),
			wantError: "request deadline cannot exceed 90 days",
		},
	}

	for _, tc := range badData {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mock.MockListingRepository{}
			// Save shouldn't be called, so no expectation set, or Expect failure?
			// Validation failures usually don't reach repository Save.
			// So default mock (AssertExpectations) works fine if we don't set any expectations, it means "no calls allowed".
			// Except we usually need to specify "no calls" explicitly or just not set "On".
			// If we don't set "On" and it calls, it will panic or fail.
			// Let's assume validation happens before Save.

			h := handler.NewListingHandler(mockRepo, nil)

			req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("User", domain.User{ID: "test-user-id", Email: "bad@test.com"})

			h.HandleCreate(c)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 Bad Request, got %d", rec.Code)
			}
			// The instruction seems to indicate that the error message should be rendered within an HTML template.
			// We'll check if the body contains the expected error message, potentially wrapped in the template structure.
			// The provided snippet `t.New("error.html").Parse(`Error Page: {{if .Message}}{{.Message}}{{end}}`)`
			// looks like an attempt to define or parse a template, which is not valid in this context.
			// Assuming the intent is to check if the error message is present in the rendered HTML.
			if !strings.Contains(rec.Body.String(), tc.wantError) {
				t.Errorf("Expected error '%s', got '%s'", tc.wantError, rec.Body.String())
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
