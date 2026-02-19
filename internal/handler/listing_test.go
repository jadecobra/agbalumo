package handler_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

// Simple Template Renderer for testing
type TestRenderer struct {
	templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewMainTemplate() *template.Template {
	t := template.New("base")
	t.New("index.html").Parse(`Index: {{len .Listings}} Listings`)
	t.New("listing_list").Parse(`{{range .Listings}}{{.Title}}{{end}}`)
	t.New("modal_detail").Parse(`{{.Listing.Title}} - {{.Listing.Description}}`)
	t.New("listing_card").Parse(`{{.Title}}`)
	t.New("admin_login.html").Parse(`Login Form: {{if .Error}}{{.Error}}{{end}}`)
	t.New("admin_dashboard.html").Parse(`Dashboard: {{len .PendingListings}} items`)
	t.New("modal_edit_listing.html").Parse(`Edit: {{.Title}}`)
	t.New("modal_feedback.html").Parse(`Feedback Modal`)
	t.New("modal_profile").Parse(`Profile: {{.User.Name}}, Listings: {{len .Listings}}`)
	t.New("error.html").Parse(`Error Page`)
	return t
}

func TestHandleHome(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", context.Background(), "", "", false).Return([]domain.Listing{
		{Title: "Test Listing 1"},
		{Title: "Test Listing 2"},
	}, nil)
	mockRepo.On("GetCounts", context.Background()).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetFeaturedListings", context.Background()).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Index: 2 Listings") {
		t.Errorf("Expected body to contain listings count, got: %s", rec.Body.String())
	}
}

func TestHandleHome_Counts(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	t_temp.New("index.html").Parse(`Total: {{.TotalCount}}, Food: {{index .Counts "Food"}}, Business: {{index .Counts "Business"}}`)
	e.Renderer = &TestRenderer{templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", context.Background(), "", "", false).Return([]domain.Listing{}, nil)
	mockRepo.On("GetCounts", context.Background()).Return(map[domain.Category]int{
		domain.Food:     5,
		domain.Business: 3,
	}, nil)
	mockRepo.On("GetFeaturedListings", context.Background()).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	expectedBody := "Total: 8, Food: 5, Business: 3"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestHandleFragment(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?q=jollof&type=Business", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", context.Background(), "Business", "jollof", false).Return([]domain.Listing{{Title: "Jollof Place"}}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Jollof Place") {
		t.Errorf("Expected body to contain listing title, got: %s", rec.Body.String())
	}
}

func TestHandleHome_Error(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", context.Background(), "", "", false).Return([]domain.Listing{}, errors.New("db connection failed"))

	h := handler.NewListingHandler(mockRepo, nil)

	_ = h.HandleHome(c)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Error Page") {
		t.Errorf("Expected friendly error page, got: %s", rec.Body.String())
	}
}

func TestHandleDetail(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/listings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", context.Background(), "1").Return(domain.Listing{Title: "Found It", Description: "Details here"}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleDetail(c); err != nil {
		t.Fatalf("HandleDetail failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Found It - Details here") {
		t.Errorf("Expected details in body, got: %s", rec.Body.String())
	}
}

func setupTestContext(method, target string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(method, target, body)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

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
					return l.Title == "Test Title" && l.HoursOfOperation == "Mon-Fri 9-5"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "ValidationError",
			body: "title=Test+Title&type=Business", // Missing required fields
			setupMock: func(m *mock.MockListingRepository) {
				// Save should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation Error",
		},
		{
			name: "RepoError",
			body: "title=Test+Title&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&address=123+St",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("Save", testifyMock.Anything, testifyMock.Anything).Return(errors.New("save failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))
			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)

			h := handler.NewListingHandler(mockRepo, nil)
			c.Set("User", domain.User{ID: "test-user-id", Email: "test@example.com"})

			err := h.HandleCreate(c)
			if err != nil {
				// handled
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.expectedBody != "" && !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, rec.Body.String())
			}
			mockRepo.AssertExpectations(t)
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
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "other-user"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "NotFound",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
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
			h := handler.NewListingHandler(mockRepo, nil)

			if err := h.HandleEdit(c); err != nil {
				// handled
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			mockRepo.AssertExpectations(t)
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
				m.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
					return l.Title == "Updated Title"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "user2", Email: "hacker@example.com"},
			body: "",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title"}, nil)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "RepoError",
			user: domain.User{ID: "user1", Email: "owner@example.com"},
			body: "title=Updated+Title&type=Business&owner_origin=Ghana&description=Updated&contact_email=new@example.com&address=123+St",
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "user1"}, nil)
				m.On("Save", testifyMock.Anything, testifyMock.Anything).Return(errors.New("update failed"))
			},
			expectedStatus: http.StatusInternalServerError,
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

			h := handler.NewListingHandler(mockRepo, nil)
			_ = h.HandleUpdate(c)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandleCreate_WithImage(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Image Listing")
	writer.WriteField("type", "Business")
	writer.WriteField("owner_origin", "Ghana")
	writer.WriteField("description", "Desc")
	writer.WriteField("contact_email", "img@example.com")
	writer.WriteField("address", "123 Image St")

	part, err := writer.CreateFormFile("image", "test.png")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	defer os.RemoveAll("ui")

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
		return l.ImageURL != "" && strings.HasPrefix(l.ImageURL, "/static/uploads/")
	})).Return(nil)

	h := handler.NewListingHandler(mockRepo, nil)
	c.Set("User", domain.User{ID: "test-user-id", Email: "test@example.com"})

	if err := h.HandleCreate(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Code)
	}
	mockRepo.AssertExpectations(t)
}

func TestHandleDelete(t *testing.T) {
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
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
				m.On("Delete", testifyMock.Anything, "1").Return(nil)
			},
			expectedStatus: http.StatusSeeOther,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "other-user"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "RepoError_Find",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("db error"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "RepoError_Delete",
			user: domain.User{ID: "owner-1"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "owner-1"}, nil)
				m.On("Delete", testifyMock.Anything, "1").Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodDelete, "/listings/1", nil)
			c.SetPath("/listings/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("User", tt.user)

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)

			h := handler.NewListingHandler(mockRepo, nil)
			_ = h.HandleDelete(c)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandleProfile(t *testing.T) {
	e := echo.New()
	t_temp := template.New("base")
	t_temp.New("modal_profile").Parse(`Profile: {{.User.Name}}, Listings: {{len .Listings}}`)
	e.Renderer = &TestRenderer{templates: t_temp}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := domain.User{ID: "u1", Name: "Test User"}
	c.Set("User", user)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAllByOwner", testifyMock.Anything, "u1").Return([]domain.Listing{
		{Title: "L1"}, {Title: "L2"},
	}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleProfile(c); err != nil {
		t.Fatalf("HandleProfile failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	expectedBody := "Profile: Test User, Listings: 2"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestHandleAbout(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewListingHandler(mockRepo, nil)

	t_temp := template.New("base")
	t_temp.New("about.html").Parse(`About Page: {{.User}}`)
	e.Renderer = &TestRenderer{templates: t_temp}

	err := h.HandleAbout(c)
	if err != nil {
		t.Fatalf("HandleAbout failed: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("About Page")) {
		t.Errorf("Expected body to contain 'About Page', got %q", rec.Body.String())
	}
}

func TestHandleClaim(t *testing.T) {
	tests := []struct {
		name           string
		user           interface{}
		setupMock      func(*mock.MockListingRepository)
		expectedStatus int
	}{
		{
			name: "Unauthenticated",
			user: nil,
			setupMock: func(m *mock.MockListingRepository) {
			},
			expectedStatus: http.StatusFound, // Redirect to login
		},
		{
			name: "ListingNotFound",
			user: domain.User{ID: "claimer"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "AlreadyOwned",
			user: domain.User{ID: "claimer"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "existing-owner"}, nil)
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Success",
			user: domain.User{ID: "claimer"},
			setupMock: func(m *mock.MockListingRepository) {
				m.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{ID: "1", OwnerID: "", Type: domain.Business}, nil)
				m.On("Save", testifyMock.Anything, testifyMock.MatchedBy(func(l domain.Listing) bool {
					return l.OwnerID == "claimer"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Renderer = &TestRenderer{templates: NewMainTemplate()}
			req := httptest.NewRequest(http.MethodPost, "/listings/1/claim", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/listings/:id/claim")
			c.SetParamNames("id")
			c.SetParamValues("1")

			if tt.user != nil {
				c.Set("User", tt.user)
			}

			mockRepo := &mock.MockListingRepository{}
			tt.setupMock(mockRepo)

			h := handler.NewListingHandler(mockRepo, nil)
			_ = h.HandleClaim(c)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandleCreate_InvalidDates(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid Deadline",
			body:           "title=T&type=Request&deadline_date=invalid-date",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Date Format",
		},
		{
			name:           "Invalid Event Start",
			body:           "title=T&type=Event&event_start=invalid-time",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Start Date Format",
		},
		{
			name:           "Invalid Event End",
			body:           "title=T&type=Event&event_end=invalid-time",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid End Date Format",
		},
		{
			name:           "Invalid Job Start",
			body:           "title=T&type=Job&job_start_date=invalid-time",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Job Start Date Format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))

			mockRepo := &mock.MockListingRepository{}
			h := handler.NewListingHandler(mockRepo, nil)
			c.Set("User", domain.User{ID: "u1"})

			err := h.HandleCreate(c)

			// Handle implementation details: some validations return echo.HTTPError, others write to c.String/c.Render
			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					if he.Code != tt.expectedStatus {
						t.Errorf("Expected status %d, got %d", tt.expectedStatus, he.Code)
					}
					if !strings.Contains(fmt.Sprintf("%v", he.Message), tt.expectedBody) {
						t.Errorf("Expected message to contain %q, got %q", tt.expectedBody, he.Message)
					}
					return
				} else {
					t.Fatalf("Unexpected non-HTTP error: %v", err)
				}
			}

			// If no error, check recorder
			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestHandleCreate_ImageUploadError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "Image Listing")
	writer.WriteField("type", "Business")

	part, _ := writer.CreateFormFile("image", "test.png")
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}

	mockImageService := &MockImageService{}
	mockImageService.On("UploadImage", testifyMock.Anything, testifyMock.Anything, testifyMock.Anything).Return("", errors.New("upload failed"))

	h := handler.NewListingHandler(mockRepo, mockImageService)
	c.Set("User", domain.User{ID: "u1"})

	err := h.HandleCreate(c)

	if err != nil {
		he, ok := err.(*echo.HTTPError)
		// handler.RespondError wraps errors, and might return HTTPError with 500
		if ok {
			if he.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 500, got %d", he.Code)
			}
			return
		}
	}

	// If it wrote to response directly
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}
}

// MockImageService for testing
type MockImageService struct {
	testifyMock.Mock
}

func (m *MockImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	args := m.Called(ctx, file, listingID)
	return args.String(0), args.Error(1)
}

// --- Profile Edge Case Tests ---

func TestHandleProfile_NoUser(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// No user set

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleProfile(c); err != nil {
		t.Fatalf("HandleProfile failed: %v", err)
	}

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect 307, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/auth/google/login" {
		t.Errorf("Expected redirect to login, got: %s", rec.Header().Get("Location"))
	}
}

func TestHandleProfile_RepoError(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := domain.User{ID: "u1", Name: "Test User"}
	c.Set("User", user)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAllByOwner", testifyMock.Anything, "u1").Return([]domain.Listing{}, errors.New("db error"))

	h := handler.NewListingHandler(mockRepo, nil)

	_ = h.HandleProfile(c)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}
}

// --- Fragment Edge Case Tests ---

func TestHandleFragment_WithHTMXHeader(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?type=Food", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "Food", "", false).Return([]domain.Listing{{Title: "Jollof Rice"}}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Jollof Rice") {
		t.Errorf("Expected body to contain listing title, got: %s", rec.Body.String())
	}
}

func TestHandleFragment_Error(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false).Return([]domain.Listing{}, errors.New("db error"))

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment returned error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}
}

// --- Home Graceful Fallback Tests ---

func TestHandleHome_CountsError_Fallback(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false).Return([]domain.Listing{{Title: "L1"}}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, errors.New("counts query failed"))
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, nil)

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome should not fail on counts error: %v", err)
	}

	// Should still render OK with empty counts
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestHandleHome_FeaturedError_Fallback(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindAll", testifyMock.Anything, "", "", false).Return([]domain.Listing{{Title: "L1"}}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{}, nil)
	mockRepo.On("GetFeaturedListings", testifyMock.Anything).Return([]domain.Listing{}, errors.New("featured query failed"))

	h := handler.NewListingHandler(mockRepo, nil)

	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome should not fail on featured error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

// --- Update Edge Case Tests ---

func TestHandleUpdate_NoUser(t *testing.T) {
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	// No user set

	mockRepo := &mock.MockListingRepository{}
	h := handler.NewListingHandler(mockRepo, nil)

	_ = h.HandleUpdate(c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rec.Code)
	}
}

func TestHandleUpdate_NotFound(t *testing.T) {
	c, rec := setupTestContext(http.MethodPost, "/listings/1", strings.NewReader(""))
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("User", domain.User{ID: "u1"})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("FindByID", testifyMock.Anything, "1").Return(domain.Listing{}, errors.New("not found"))

	h := handler.NewListingHandler(mockRepo, nil)

	_ = h.HandleUpdate(c)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}
