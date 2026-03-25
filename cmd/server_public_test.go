package cmd_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicRoutes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   []string
	}{
		{
			name:           "Home loads",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   []string{"listings-container"}, // "Verify homepage loads and shows listings"
		},
		{
			name:           "About loads",
			method:         http.MethodGet,
			path:           "/about",
			expectedStatus: http.StatusOK,
			expectedBody:   []string{"about"},
		},
		{
			name:           "Search/Filter Fragment loads",
			method:         http.MethodGet,
			path:           "/listings/fragment",
			expectedStatus: http.StatusOK,
			expectedBody:   []string{"card-juicy", "hx-get=", "/listings/"},
		},
		{
			name:           "Google OAuth initiates",
			method:         http.MethodGet,
			path:           "/auth/google/login",
			expectedStatus: http.StatusTemporaryRedirect, // OAuth will redirect
			expectedBody:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			body := rec.Body.String()
			for _, exp := range tt.expectedBody {
				assert.Contains(t, body, exp)
			}
		})
	}
}

func TestListingDetail(t *testing.T) {
	// Let's rely on the seeder having created "Lagos Import Export"
	// We can fetch the fragment and extract an ID or just search via API?
	// The seeder sets Title = "Lagos Import Export" and generates a UUID.
	// Easiest is to hit the DB directly via an injected or available repo, but we don't have it here.
	// Instead, let's parse the fragment response for `/listings/`
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body := rec.Body.String()
	idx := strings.Index(body, "hx-get=\"/listings/")
	if idx == -1 {
		t.Fatal("Could not find any listing in the fragment response to test Detail view")
	}

	// Extract the UUID format: "hx-get=\"/listings/UUID\""
	sub := body[idx+len("hx-get=\"/listings/"):]
	endIdx := strings.Index(sub, "\"")
	if endIdx == -1 {
		t.Fatal("Could not parse listing ID from fragment")
	}

	listingID := sub[:endIdx]

	// Now fetch the detail modal directly
	reqDetail := httptest.NewRequest(http.MethodGet, "/listings/"+listingID, nil)
	recDetail := httptest.NewRecorder()
	e.ServeHTTP(recDetail, reqDetail)

	assert.Equal(t, http.StatusOK, recDetail.Code)
	// It should render modal_detail which has the close button with data-modal-action
	assert.Contains(t, recDetail.Body.String(), "data-modal-action=\"close\"")
}
