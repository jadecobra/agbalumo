package handler_test

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/service"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreate(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setup          func(t *testing.T, repo domain.ListingRepository)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			body: "title=Test+Title&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&hours_of_operation=Mon-Fri+9-5&city=Lagos&address=123+Street",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				// No extra setup needed
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ValidationError",
			body:           "title=Test+Title&type=Business",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Page",
		},
		{
			name:           "RequestWithoutDeadline",
			body:           "title=Req&type=Request&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&city=Lagos&address=123+St",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "BusinessWithNoPrefixURL",
			body:           "title=Biz&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&city=Lagos&address=123+Street&website_url=example.com",
			setup:          func(t *testing.T, repo domain.ListingRepository) {},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			repo := handler.SetupTestRepository(t)
			tt.setup(t, repo)

			listingSvc := service.NewListingService(repo, repo, repo)

			h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
			c.Set("User", domain.User{ID: "test-user-id"})

			_ = h.HandleCreate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandleEdit(t *testing.T) {
	tests := []struct {
		name           string
		user           domain.User
		setup          func(t *testing.T, repo domain.ListingRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "other-user", Role: domain.UserRoleUser},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodGet, "/listings/1/edit", nil)
			c.SetPath("/listings/:id/edit")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			repo := handler.SetupTestRepository(t)
			tt.setup(t, repo)
			listingSvc := service.NewListingService(repo, repo, repo)
			h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})

			_ = h.HandleEdit(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestHandleUpdate(t *testing.T) {
	tests := []struct {
		name           string
		user           domain.User
		body           string
		setup          func(t *testing.T, repo domain.ListingRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "user1", Email: "owner@example.com"},
			body: "title=Updated+Title&type=Business&owner_origin=Ghana&description=Updated&contact_email=new@example.com&city=Accra&address=123+St",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "user2", Email: "hacker@example.com", Role: domain.UserRoleUser},
			body: "",
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(tt.body))
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			repo := handler.SetupTestRepository(t)
			tt.setup(t, repo)

			listingSvc := service.NewListingService(repo, repo, repo)

			h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
			_ = h.HandleUpdate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestHandleDelete(t *testing.T) {
	tests := []struct {
		name       string
		user       interface{}
		setup      func(t *testing.T, repo domain.ListingRepository)
		expectCode int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectCode: http.StatusSeeOther,
		},
		{
			name:       "NoUser_Unauthorized",
			user:       nil,
			setup:      func(t *testing.T, repo domain.ListingRepository) {},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "NotFound",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
			},
			expectCode: http.StatusNotFound,
		},
		{
			name: "Forbidden_NotOwner",
			user: domain.User{ID: "other-user"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
			},
			expectCode: http.StatusForbidden,
		},
		{
			name: "DeleteError",
			user: domain.User{ID: "owner-1"},
			setup: func(t *testing.T, repo domain.ListingRepository) {
				// We can't trigger a DB error easily with real SQLite without some trickery
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "DeleteError" {
				t.Skip("Hard to trigger DB error with real SQLite")
			}
			c, rec := setupTestContext(http.MethodDelete, "/listings/1", nil)
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			if tt.user != nil {
				c.Set("User", tt.user)
			}

			repo := handler.SetupTestRepository(t)
			tt.setup(t, repo)

			listingSvc := service.NewListingService(repo, repo, repo)

			h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
			_ = h.HandleDelete(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

func TestHandleClaim(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	_ = repo.Save(context.Background(), domain.Listing{ID: "1", Title: "Biz", Type: domain.Business, Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", City: "Lagos", Address: "123 St"})
	_ = repo.SaveCategory(context.Background(), domain.CategoryData{ID: string(domain.Business), Name: string(domain.Business), Claimable: true, Active: true})

	c, rec := setupTestContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "claimer", Name: "Claimer", Email: "c@e.com"})

	listingSvc := service.NewListingService(repo, repo, repo)

	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleUpdate_NotFound(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1"})

	listingSvc := service.NewListingService(repo, repo, repo)

	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleUpdate_NoUser(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	listingSvc := service.NewListingService(repo, repo, repo)

	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleUpdate_DuplicateTitle(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	_ = repo.Save(context.Background(), domain.Listing{ID: "1", OwnerID: "user1", Title: "Old", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})
	_ = repo.Save(context.Background(), domain.Listing{ID: "2", OwnerID: "user2", Title: "Taken Title", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "456 St"})

	body := "title=Taken+Title&type=Business&owner_origin=Ghana&description=Desc&contact_email=t@e.com&city=Kumasi&address=123+St"
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(body))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1"})

	listingSvc := service.NewListingService(repo, repo, repo)

	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleCreate_NoUser(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	body := "title=Test&type=Business&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&city=Lagos&address=123+St"
	c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(body))

	listingSvc := service.NewListingService(repo, repo, repo)

	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleCreate_DuplicateTitle(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	_ = repo.Save(context.Background(), domain.Listing{ID: "x", Title: "Existing", Status: domain.ListingStatusApproved, IsActive: true, OwnerOrigin: "Nigeria", ContactEmail: "test@example.com", Type: domain.Business, City: "Lagos", Address: "123 St"})

	body := "title=Existing&type=Business&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&city=Lagos&address=123+St"
	c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(body))
	c.Set("User", domain.User{ID: "user1"})

	listingSvc := service.NewListingService(repo, repo, repo)

	h := handler.NewListingHandler(repo, repo, listingSvc, nil, &handler.MockGeocodingService{}, &config.Config{})
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
