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

func TestHandleUpdate_JobSuccess(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	// Existing Job Listing
	existingListing := domain.Listing{
		ID:           "job-1",
		OwnerID:      "owner-1",
		Type:         domain.Job,
		Title:        "Senior Go Dev",
		Description:  "Write Go code",
		Company:      "Tech Corp",
		Skills:       "Go, SQL",
		PayRange:     "100k-150k",
		JobStartDate: time.Now().Add(24 * time.Hour),
		JobApplyURL:  "https://example.com",
		City:         "Lagos",
		OwnerOrigin:  "Nigeria",
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	// This simulates the form data that the updated UI will send
	// Note: We use the format expected by time.Parse inside the handler
	jobStart := time.Now().Add(48 * time.Hour).Format("2006-01-02T15:04")
	
	formData := "title=Senior+Go+Dev+Updated&type=Job&owner_origin=Nigeria&description=Updated+Desc&contact_email=job@example.com&city=Lagos" +
		"&company=Updated+Corp&skills=Go,+Rust&pay_range=200k&job_apply_url=https://updated.com&job_start_date=" + jobStart

	req := httptest.NewRequest(http.MethodPost, "/listings/job-1", strings.NewReader(formData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("job-1")
	c.Set("User", domain.User{ID: "owner-1"})

	// Mock Repo
	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			return existingListing, nil
		},
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			// Verify Update fields
			if l.Title != "Senior Go Dev Updated" {
				t.Errorf("Expected updated title, got %s", l.Title)
			}
			if l.Company != "Updated Corp" {
				t.Errorf("Expected Updated Corp, got %s", l.Company)
			}
			if l.Skills != "Go, Rust" {
				t.Errorf("Expected Go, Rust, got %s", l.Skills)
			}
			if l.PayRange != "200k" {
				t.Errorf("Expected 200k, got %s", l.PayRange)
			}
			return nil
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	err := h.HandleUpdate(c)

	if err != nil {
		t.Fatalf("HandleUpdate failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
		t.Logf("Response Body: %s", rec.Body.String())
	}
}
