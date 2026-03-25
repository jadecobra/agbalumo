package listing_test

import (
	listmod "github.com/jadecobra/agbalumo/internal/module/listing"

	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/stretchr/testify/assert"
)
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
			listingSvc := listmod.NewListingService(repo, repo, repo)
			h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})

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

			listingSvc := listmod.NewListingService(repo, repo, repo)

			h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
			_ = h.HandleUpdate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
func TestHandleUpdate_NotFound(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1"})

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestHandleUpdate_NoUser(t *testing.T) {
	repo := handler.SetupTestRepository(t)
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
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

	listingSvc := listmod.NewListingService(repo, repo, repo)

	h := listmod.NewListingHandler(listmod.ListingDependencies{
		ListingStore:     repo,
		CategoryStore:    repo,
		ListingSvc:       listingSvc,
		ImageService:     nil,
		GeocodingSvc:     &MockGeocodingService{},
		Config:           &config.Config{},
	})
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
