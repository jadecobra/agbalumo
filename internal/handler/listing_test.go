package handler_test

import (
	"bytes"
	"context"
	"errors"
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
)

// Simple Template Renderer for testing
type TestRenderer struct {
	templates *template.Template
}

func (t *TestRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewMainTemplate() *template.Template {
	// We need to parse a minimal set of templates for the handler to work.
	// In a real scenario we might load from disk, but for unit tests strings are safer/faster.
	t := template.New("base")
	// index.html is a full page in our app, not a partial define
	t.New("index.html").Parse(`Index: {{len .Listings}} Listings`)
	t.New("listing_list.html").Parse(`{{range .Listings}}{{.Title}}{{end}}`)
	t.New("modal_detail.html").Parse(`{{.Listing.Title}} - {{.Listing.Description}}`)
	t.New("listing_card.html").Parse(`{{.Title}}`)
	t.New("admin_login.html").Parse(`Login Form: {{if .Error}}{{.Error}}{{end}}`)
	t.New("admin_dashboard.html").Parse(`Dashboard: {{len .Listings}} items`)
	t.New("modal_edit_listing.html").Parse(`Edit: {{.Title}}`)
	t.New("error.html").Parse(`Error Page`)
	return t
}

func TestHandleHome(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Repo
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			return []domain.Listing{
				{Title: "Test Listing 1"},
				{Title: "Test Listing 2"},
			}, nil
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	// Verify
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Index: 2 Listings") {
		t.Errorf("Expected body to contain listings count, got: %s", rec.Body.String())
	}
}

func TestHandleHome_Counts(t *testing.T) {
	// Setup
	e := echo.New()
	t_temp := template.New("base")
	t_temp.New("index.html").Parse(`Total: {{.TotalCount}}, Food: {{index .Counts "Food"}}, Business: {{index .Counts "Business"}}`)
	e.Renderer = &TestRenderer{templates: t_temp}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Repo
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			return []domain.Listing{}, nil
		},
		GetCountsFn: func(ctx context.Context) (map[domain.Category]int, error) {
			return map[domain.Category]int{
				domain.Food:     5,
				domain.Business: 3,
			}, nil
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	if err := h.HandleHome(c); err != nil {
		t.Fatalf("HandleHome failed: %v", err)
	}

	// Verify
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	expectedBody := "Total: 8, Food: 5, Business: 3"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}


func TestHandleFragment(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/listings/fragment?q=jollof&type=Business", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Repo
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			if filterType != "Business" {
				t.Errorf("Expected filterType Business, got %s", filterType)
			}
			if query != "jollof" {
				t.Errorf("Expected query jollof, got %s", query)
			}
			return []domain.Listing{{Title: "Jollof Place"}}, nil
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	if err := h.HandleFragment(c); err != nil {
		t.Fatalf("HandleFragment failed: %v", err)
	}

	// Verify
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Jollof Place") {
		t.Errorf("Expected body to contain listing title, got: %s", rec.Body.String())
	}
}

func TestHandleHome_Error(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock Repo Error
	mockRepo := &mock.MockListingRepository{
		FindAllFn: func(ctx context.Context, filterType, query string, includeInactive bool) ([]domain.Listing, error) {
			return nil, errors.New("db connection failed")
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	_ = h.HandleHome(c)

	// The handler writes 500 to response.
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Error Page") {
		t.Errorf("Expected friendly error page, got: %s", rec.Body.String())
	}
}

func TestHandleDetail(t *testing.T) {
	// Setup
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodGet, "/listings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/listings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Mock Repo
	mockRepo := &mock.MockListingRepository{
		FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
			if id == "1" {
				return domain.Listing{Title: "Found It", Description: "Details here"}, nil
			}
			return domain.Listing{}, errors.New("not found")
		},
	}

	h := handler.NewListingHandler(mockRepo)

	// Execute
	if err := h.HandleDetail(c); err != nil {
		t.Fatalf("HandleDetail failed: %v", err)
	}

	// Verify
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
		mockSetup      func() *mock.MockListingRepository
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			body: "title=Test+Title&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&hours_of_operation=Mon-Fri+9-5&address=123+Street",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					SaveFn: func(ctx context.Context, l domain.Listing) error {
						if l.Title != "Test Title" {
							return errors.New("unexpected title")
						}
						// TDD: Check if HoursOfOperation is extracted
						// Note: This needs the input body to include it, adjusting body below
						if l.HoursOfOperation != "Mon-Fri 9-5" {
							return errors.New("expected HoursOfOperation to be 'Mon-Fri 9-5'")
						}
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "ValidationError",
			body: "title=Test+Title&type=Business", // Missing required fields
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					SaveFn: func(ctx context.Context, l domain.Listing) error {
						return errors.New("shoud not be called")
					},
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Validation Error",
		},
		{
			name: "RepoError",
			body: "title=Test+Title&type=Business&owner_origin=Nigeria&description=Cool&contact_email=test@example.com&address=123+St",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					SaveFn: func(ctx context.Context, l domain.Listing) error {
						return errors.New("save failed")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestContext(http.MethodPost, "/listings", strings.NewReader(tt.body))

			h := handler.NewListingHandler(tt.mockSetup())

			// Inject User for Auth
			c.Set("User", domain.User{ID: "test-user-id", Email: "test@example.com"})

			err := h.HandleCreate(c)
			if err != nil {
				// Some errors are handled by helper, but checking response code covers it
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
			if tt.expectedBody != "" && !strings.Contains(rec.Body.String(), tt.expectedBody) {
				t.Errorf("Expected body to contain %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestHandleEdit(t *testing.T) {
	tests := []struct {
		name           string
		user           domain.User
		mockSetup      func() *mock.MockListingRepository
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "owner-1"},
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
						return domain.Listing{ID: "1", OwnerID: "owner-1", Title: "Title"}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "other-user"},
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
						return domain.Listing{ID: "1", OwnerID: "owner-1"}, nil
					},
				}
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "NotFound",
			user: domain.User{ID: "owner-1"},
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
						return domain.Listing{}, errors.New("not found")
					},
				}
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

			h := handler.NewListingHandler(tt.mockSetup())

			if err := h.HandleEdit(c); err != nil {
				// Echo handler returns error for some statuses, mainly we check response code
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandleUpdate(t *testing.T) {
	tests := []struct {
		name           string
		user           domain.User
		body           string
		mockSetup      func() *mock.MockListingRepository
		expectedStatus int
	}{
		{
			name: "Success",
			user: domain.User{ID: "user1", Email: "owner@example.com"},
			body: "title=Updated+Title&type=Business&owner_origin=Ghana&description=Updated&contact_email=new@example.com&address=123+St",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
						return domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title"}, nil
					},
					SaveFn: func(ctx context.Context, l domain.Listing) error {
						if l.Title != "Updated Title" {
							return errors.New("unexpected title")
						}
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Forbidden",
			user: domain.User{ID: "user2", Email: "hacker@example.com"},
			body: "",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
						return domain.Listing{ID: "1", OwnerID: "user1", Title: "Old Title"}, nil
					},
				}
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "RepoError",
			user: domain.User{ID: "user1", Email: "owner@example.com"},
			body: "title=Updated+Title&type=Business&owner_origin=Ghana&description=Updated&contact_email=new@example.com&address=123+St",
			mockSetup: func() *mock.MockListingRepository {
				return &mock.MockListingRepository{
					FindByIDFn: func(ctx context.Context, id string) (domain.Listing, error) {
						return domain.Listing{ID: "1", OwnerID: "user1"}, nil
					},
					SaveFn: func(ctx context.Context, l domain.Listing) error {
						return errors.New("update failed")
					},
				}
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

			h := handler.NewListingHandler(tt.mockSetup())

			_ = h.HandleUpdate(c)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandleCreate_WithImage(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	// Create Multipart Form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Fields
	writer.WriteField("title", "Image Listing")
	writer.WriteField("type", "Business")
	writer.WriteField("owner_origin", "Ghana")
	writer.WriteField("description", "Desc")
	writer.WriteField("contact_email", "img@example.com")
	writer.WriteField("address", "123 Image St")

	// File
	part, err := writer.CreateFormFile("image", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake image content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/listings", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Clean up uploads dir after test
	defer os.RemoveAll("ui")

	mockRepo := &mock.MockListingRepository{
		SaveFn: func(ctx context.Context, l domain.Listing) error {
			if l.ImageURL == "" {
				t.Error("Expected ImageURL to be set")
			}
			if !strings.HasPrefix(l.ImageURL, "/static/uploads/") {
				t.Errorf("Expected ImageURL path /static/uploads/, got %s", l.ImageURL)
			}
			return nil
		},
	}
	h := handler.NewListingHandler(mockRepo)

	// Inject User for Auth
	c.Set("User", domain.User{ID: "test-user-id", Email: "test@example.com"})

	if err := h.HandleCreate(c); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Code)
	}
}
