package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleHome(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", ctx, "", "", "", "", false, 20, 0).Return([]domain.Listing{
		{ID: "1", Title: "Listing 1", Type: domain.Business, IsActive: true},
	}, nil)
	mockRepo.On("GetCounts", ctx).Return(map[domain.Category]int{domain.Business: 1}, nil)
	mockRepo.On("GetLocations", ctx).Return([]string{"Lagos"}, nil)
	mockRepo.On("GetFeaturedListings", ctx).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCategories", ctx, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	if err := h.HandleHome(c); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Listing 1")
}

func TestHandleDetail(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/listings/1", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	ctx := context.Background()

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", ctx, "1").Return(domain.Listing{
		ID:    "1",
		Title: "Detail View",
		Type:  domain.Business,
	}, nil)
	mockRepo.On("GetCategory", ctx, string(domain.Business)).Return(domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	if err := h.HandleDetail(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "Detail View")
}

func TestHandleProfile(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/profile", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	user := domain.User{ID: "u1", Name: "John Doe"}
	c.Set("User", user)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAllByOwner", ctx, "u1", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{
		{ID: "1", Title: "My Listing"},
	}, nil)
	mockRepo.On("GetCategories", ctx, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	if err := h.HandleProfile(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "John Doe")
}

func TestHandleAbout(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/about", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetCategories", ctx, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	if err := h.HandleAbout(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "About agbalumo")
}

func TestHandleFragment(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := setupRequest(http.MethodGet, "/listings/fragment?q=test", nil)
	rec := setupResponseRecorder()
	c := e.NewContext(req, rec)
	ctx := context.Background()

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", ctx, "", "test", "", "", false, 20, 0).Return([]domain.Listing{
		{ID: "1", Title: "Search Result", Type: domain.Business, IsActive: true},
	}, nil)

	h := handler.NewListingHandler(mockRepo, nil, "")
	if err := h.HandleFragment(c); err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, rec.Body.String(), "Search Result")
}
