package handler_test

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleHome_FeaturedPrioritization(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("index.html").Parse(`Listings: {{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &TestRenderer{templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}

	featured := []domain.Listing{
		{ID: "f1", Title: "Featured 1", Featured: true},
		{ID: "f2", Title: "Featured 2", Featured: true},
	}

	regular := []domain.Listing{
		{ID: "r1", Title: "Regular 1", Featured: false},
		{ID: "f1", Title: "Featured 1", Featured: true}, // Also in regular results
		{ID: "r2", Title: "Regular 2", Featured: false},
	}

	mockRepo.On("GetFeaturedListings", context.Background()).Return(featured, nil)
	mockRepo.On("FindAll", context.Background(), "", "", "", "", false, 20, 0).Return(regular, nil)
	mockRepo.On("GetCounts", context.Background()).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetLocations", context.Background()).Return([]string{}, nil)

	h := handler.NewListingHandler(mockRepo, nil, "")

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	// EXPECTED: Featured 1, Featured 2, Regular 1, Regular 2 (no duplicate f1)
	expectedBody := "Listings: Featured 1,Featured 2,Regular 1,Regular 2,"
	assert.Equal(t, expectedBody, rec.Body.String())
}

func TestHandleFragment_FeaturedPrioritization(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	template.Must(t_temp.New("listing_list").Parse(`{{range .Listings}}{{.Title}},{{end}}`))
	e.Renderer = &TestRenderer{templates: t_temp}

	// Page 1, no filters
	req := httptest.NewRequest(http.MethodGet, "/listings?page=1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}

	featured := []domain.Listing{{ID: "f1", Title: "Featured 1"}}
	regular := []domain.Listing{{ID: "r1", Title: "Regular 1"}}

	mockRepo.On("GetFeaturedListings", context.Background()).Return(featured, nil)
	mockRepo.On("FindAll", context.Background(), "", "", "", "", false, 20, 0).Return(regular, nil)

	h := handler.NewListingHandler(mockRepo, nil, "")

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	expectedBody := "Featured 1,Regular 1,"
	assert.Equal(t, expectedBody, rec.Body.String())
}
