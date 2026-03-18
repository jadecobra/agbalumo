package handler_test

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestListingHandler_FormParsing(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setup          func(t *testing.T, repo domain.ListingRepository)
		expectedStatus int
		verify         func(t *testing.T, repo domain.ListingRepository)
	}{
		{
			name:           "Success_EventWithDates",
			body:           "title=Event+Test&type=Event&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St&city=Lagos&event_start=2027-12-01T10:00&event_end=2027-12-01T12:00",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, repo domain.ListingRepository) {
				listings, _, err := repo.FindAll(context.Background(), "", "Event Test", "", "", false, 1, 0)
				assert.NoError(t, err)
				if assert.Len(t, listings, 1) {
					assert.Equal(t, 2027, listings[0].EventStart.Year())
					assert.Equal(t, 12, int(listings[0].EventStart.Month()))
					assert.Equal(t, 1, listings[0].EventStart.Day())
				}
			},
		},
		{
			name:           "Success_JobWithStartDateAndURL",
			body:           "title=Job+Test&type=Job&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St&job_start_date=2027-12-01T09:00&job_apply_url=example.com&company=Acme&skills=Golang&pay_range=100k-200k&city=Lagos",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, repo domain.ListingRepository) {
				listings, _, err := repo.FindAll(context.Background(), "", "Job Test", "", "", false, 1, 0)
				assert.NoError(t, err)
				if assert.Len(t, listings, 1) {
					assert.Equal(t, 2027, listings[0].JobStartDate.Year())
					assert.Equal(t, "https://example.com", listings[0].JobApplyURL)
				}
			},
		},
		{
			name:           "Success_RequestWithDeadline",
			body:           "title=Request+Test&type=Request&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St&city=Lagos&deadline_date=2026-04-30",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
			verify: func(t *testing.T, repo domain.ListingRepository) {
				listings, _, err := repo.FindAll(context.Background(), "", "Request Test", "", "", false, 1, 0)
				assert.NoError(t, err)
				if assert.Len(t, listings, 1) {
					assert.Equal(t, 2026, listings[0].Deadline.Year())
					assert.Equal(t, 4, int(listings[0].Deadline.Month()))
				}
			},
		},
		{
			name:           "Failure_InvalidEventDate",
			body:           "title=Bad+Event&type=Event&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St&city=Lagos&event_start=invalid",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := handler.SetupTestRepository(t)
			tt.setup(t, repo)

			listingSvc := service.NewListingService(repo, repo, repo)

			h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			c.Set("User", domain.User{ID: "user-1"})

			err := h.HandleCreate(c)
			if err != nil {
				t.Logf("HandleCreate error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Fatalf("Test %s failed: expected %d, got %d. Body: %s", tt.name, tt.expectedStatus, rec.Code, rec.Body.String())
			}
			if tt.verify != nil {
				tt.verify(t, repo)
			}
		})
	}
}

func TestListingHandler_URLNormalization(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	listingSvc := service.NewListingService(repo, repo, repo)
	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"NoPrefix", "jadecobra.dev", "https://jadecobra.dev"},
		{"HttpPrefix", "http://example.com", "http://example.com"},
		{"HttpsPrefix", "https://test.com", "https://test.com"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := "title=URL+Test+" + tt.name + "&type=Business&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St&city=Lagos&website_url=" + tt.url
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(body))
			c.Set("User", domain.User{ID: "user-1"})

			_ = h.HandleCreate(c)

			assert.Equal(t, http.StatusOK, rec.Code)
			listings, _, err := repo.FindAll(context.Background(), "", "URL Test "+tt.name, "", "", false, 1, 0)
			assert.NoError(t, err)
			if assert.Len(t, listings, 1) {
				assert.Equal(t, tt.expected, listings[0].WebsiteURL)
			}
		})
	}
}
