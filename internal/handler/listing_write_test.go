package handler_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleCreate(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setupMock      func(*mock.MockListingRepository)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			body: "title=Test+Title&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&hours_of_operation=Mon-Fri+9-5&address=123+Street",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
					return l.Title == "Test Title"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "ValidationError",
			body: "title=Test+Title&type=Business",
			setupMock: func(m *mock.MockListingRepository) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Error Page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)
			mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
			mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

			h := handler.NewListingHandler(mockRepo, nil, "")
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
		setupMock      func(*mock.MockListingRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title"}, nil)
				m.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "other-user", Role: domain.UserRoleUser},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
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

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)
			h := handler.NewListingHandler(mockRepo, nil, "")

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
		setupMock      func(*mock.MockListingRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "user1", Email: "owner@example.com"},
			body: "title=Updated+Title&type=Business&owner_origin=Ghana&description=Updated&contact_email=new@example.com&address=123+St",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title"}, nil)
				m.On("Save", testifyMock.Anything, testifyMock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "user2", Email: "hacker@example.com", Role: domain.UserRoleUser},
			body: "",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title"}, nil)
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

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)
			mockRepo.On("FindByTitle", testifyMock.Anything, testifyMock.Anything).Return([]domain.Listing{}, nil).Maybe()
			mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

			h := handler.NewListingHandler(mockRepo, nil, "")
			_ = h.HandleUpdate(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestHandleDelete(t *testing.T) {
	tests := []struct {
		name       string
		user       interface{}
		setupMock  func(*mock.MockListingRepository)
		expectCode int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
				m.On("Delete", testifyMock.Anything, "1").Return(nil)
			},
			expectCode: http.StatusSeeOther,
		},
		{
			name:       "NoUser_Unauthorized",
			user:       nil,
			setupMock:  func(m *mock.MockListingRepository) {},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "NotFound",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("not found"))
			},
			expectCode: http.StatusNotFound,
		},
		{
			name: "Forbidden_NotOwner",
			user: domain.User{ID: "other-user"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
			},
			expectCode: http.StatusForbidden,
		},
		{
			name: "DeleteError",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
				m.On("Delete", testifyMock.Anything, "1").Return(errors.New("db error"))
			},
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodDelete, "/listings/1", nil)
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			if tt.user != nil {
				c.Set("User", tt.user)
			}

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)
			h := handler.NewListingHandler(mockRepo, nil, "")
			_ = h.HandleDelete(c)

			assert.Equal(t, tt.expectCode, rec.Code)
		})
	}
}

func TestHandleClaim(t *testing.T) {
	c, rec := setupTestContext(http.MethodPost, "/listings/1/claim", nil)
	c.SetPath("/listings/:id/claim")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "claimer", Name: "Claimer", Email: "c@e.com"})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", Title: "Biz", Type: domain.Business}, nil)
	mockRepo.On("GetCategory", testifyMock.Anything, string(domain.Business)).Return(domain.CategoryData{Claimable: true}, nil)
	mockRepo.On("GetClaimRequestByUserAndListing", testifyMock.Anything, "claimer", "1").Return(domain.ClaimRequest{}, errors.New("not found"))
	mockRepo.On("SaveClaimRequest", testifyMock.Anything, testifyMock.Anything).Return(nil)

	h := handler.NewListingHandler(mockRepo, nil, "")
	_ = h.HandleClaim(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandleUpdate_NotFound(t *testing.T) {
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1"})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("not found"))

	h := handler.NewListingHandler(mockRepo, nil, "")
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleUpdate_NoUser(t *testing.T) {
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewListingHandler(mockRepo, nil, "")
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleUpdate_DuplicateTitle(t *testing.T) {
	body := "title=Taken+Title&type=Business&owner_origin=Ghana&description=Desc&contact_email=t@e.com&address=123+St"
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(body))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "user1"})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "user1", Title: "Old"}, nil)
	mockRepo.On("FindByTitle", testifyMock.Anything, "Taken Title").Return([]domain.Listing{{ID: "2", Title: "Taken Title"}}, nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	_ = h.HandleUpdate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleCreate_NoUser(t *testing.T) {
	body := "title=Test&type=Business&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St"
	c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(body))

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewListingHandler(mockRepo, nil, "")
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandleCreate_DuplicateTitle(t *testing.T) {
	body := "title=Existing&type=Business&owner_origin=Nigeria&description=Cool&contact_email=t@e.com&address=123+St"
	c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(body))
	c.Set("User", domain.User{ID: "user1"})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByTitle", testifyMock.Anything, "Existing").Return([]domain.Listing{{ID: "x", Title: "Existing"}}, nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{}, nil).Maybe()

	h := handler.NewListingHandler(mockRepo, nil, "")
	_ = h.HandleCreate(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
