package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
)

// Data Validation Integration Tests
func TestIntegration_DataValidation(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	// 1. Positive Test Case: Known Good Data (Subset of seed data)
	goodData := []string{
		"title=Good+Biz&type=Business&owner_origin=Nigeria&description=Valid&contact_email=good@test.com",
		"title=Good+Req&type=Request&owner_origin=Ghana&description=Valid&contact_whatsapp=+123456&deadline_date=" + time.Now().Add(24*time.Hour).Format("2006-01-02"),
	}

	for _, bodyStr := range goodData {
		req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(bodyStr))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock Repo - Success
		mockRepo := &mock.MockListingRepository{
			SaveFn: func(ctx context.Context, l domain.Listing) error {
				return nil
			},
		}
		h := handler.NewListingHandler(mockRepo)

		if err := h.HandleCreate(c); err != nil {
			t.Errorf("Unexpected error for good data: %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for good data, got %d. Body: %s", rec.Code, rec.Body.String())
		}
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
			body:      "title=Bad+Biz&type=Business&owner_origin=Nigeria&description=Ghost",
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
			req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockRepo := &mock.MockListingRepository{SaveFn: func(ctx context.Context, l domain.Listing) error { return nil }}
			h := handler.NewListingHandler(mockRepo)

			h.HandleCreate(c) // Should return error handled by echo or return error

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 Bad Request, got %d", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tc.wantError) {
				t.Errorf("Expected error '%s', got '%s'", tc.wantError, rec.Body.String())
			}
		})
	}
}
