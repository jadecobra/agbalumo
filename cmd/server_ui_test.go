package cmd_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/cmd"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var e *echo.Echo

func TestMain(m *testing.M) {
	os.Setenv("AGBALUMO_ENV", "development")
	// Keep ENV=test for test compatibility but set high rate limits to avoid 429 in tests
	os.Setenv("RATE_LIMIT_RATE", "10000")
	os.Setenv("RATE_LIMIT_BURST", "20000")
	os.Setenv("DB_URL", "file:test_ui.db?mode=memory&cache=shared")
	// SetupServer handles seeding as long as ENV != "production"
	var err error

	// We need to change to the project root directory so template paths work
	os.Chdir("..")

	e, err = cmd.SetupServer()
	if err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

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

func getSessionCookie(rec *httptest.ResponseRecorder) string {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "auth_session" {
			return cookie.String()
		}
	}
	return ""
}

func TestUserRoutes(t *testing.T) {
	// 1. Dev Login
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=dev@agbalumo.com", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code) // Should redirect to /
	cookie := getSessionCookie(rec)
	if cookie == "" {
		for _, c := range rec.Result().Cookies() {
			t.Logf("Found cookie: %s=%s", c.Name, c.Value)
			// Gorilla sessions default name is "session" unless specified
			if c.Name == "session" {
				cookie = c.String()
			}
		}
	}
	assert.NotEmpty(t, cookie, "Expected session cookie after login")

	// 2. User Profile
	req = httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, "dev@agbalumo.com") // Profile shows email

	// 3. Create Listing
	// Extract CSRF token
	idx := strings.Index(body, "name=\"_csrf\" value=\"")
	if idx == -1 {
		t.Fatal("Could not find CSRF token in profile page")
	}
	sub := body[idx+len("name=\"_csrf\" value=\""):]
	endIdx := strings.Index(sub, "\"")
	csrfToken := sub[:endIdx]

	var csrfCookie string
	for _, c := range rec.Result().Cookies() {
		if c.Name == "_csrf" {
			csrfCookie = c.String()
		}
	}

	title := fmt.Sprintf("My Test Listing %d", time.Now().UnixNano())
	form := strings.NewReader("title=" + title + "&type=Service&owner_origin=Nigeria&description=Testing&contact_email=dev@test.com&_csrf=" + csrfToken)
	req = httptest.NewRequest(http.MethodPost, "/listings", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie+"; "+csrfCookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code) // Returns 200 OK with the new listing card
	assert.Contains(t, rec.Body.String(), title)

	// Extract listing ID from the card
	createdBody := rec.Body.String()
	idx = strings.Index(createdBody, "id=\"listing-")
	if idx != -1 {
		subStr := createdBody[idx+len("id=\"listing-"):]
		endId := strings.Index(subStr, "\"")
		if endId != -1 {
			listingID := subStr[:endId]

			// Edit
			req = httptest.NewRequest(http.MethodGet, "/listings/"+listingID+"/edit", nil)
			req.Header.Set("Cookie", cookie)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)

			// Update
			form = strings.NewReader("title=" + title + "+Updated&type=Service&owner_origin=Nigeria&description=Testing&contact_email=dev@test.com&_csrf=" + csrfToken)
			req = httptest.NewRequest(http.MethodPost, "/listings/"+listingID, form)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Cookie", cookie+"; "+csrfCookie)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Updated")

			// Claim (returns Forbidden/already owned as we are the owner, but it tests the route)
			req = httptest.NewRequest(http.MethodPost, "/listings/"+listingID+"/claim", nil)
			req.Header.Set("Cookie", cookie+"; "+csrfCookie)
			req.Header.Set("X-CSRF-Token", csrfToken)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.NotEqual(t, http.StatusNotFound, rec.Code)

			// Delete
			req = httptest.NewRequest(http.MethodDelete, "/listings/"+listingID, nil)
			req.Header.Set("Cookie", cookie+"; "+csrfCookie)
			req.Header.Set("X-CSRF-Token", csrfToken)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusSeeOther, rec.Code) // redirects
		}
	}

	// Feedback Modal & Submit
	req = httptest.NewRequest(http.MethodGet, "/feedback/modal", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	form = strings.NewReader("content=Great+App&type=bug&_csrf=" + csrfToken)
	req = httptest.NewRequest(http.MethodPost, "/feedback", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie+"; "+csrfCookie)
	req.Header.Set("X-CSRF-Token", csrfToken)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Logout
	req = httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code) // Should redirect (307)
}

func TestAdminRoutes(t *testing.T) {
	adminEmail := fmt.Sprintf("admin-dev-%d@agbalumo.com", time.Now().UnixNano())
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email="+adminEmail, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	cookie := getSessionCookie(rec)
	assert.NotEmpty(t, cookie, "Expected session cookie")

	// 1. Admin Login View (accessible without admin auth)
	req = httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusTemporaryRedirect {
		t.Logf("Redirected to: %s", rec.Header().Get("Location"))
	}
	assert.Equal(t, http.StatusOK, rec.Code)

	// Extract CSRF
	body := rec.Body.String()
	idx := strings.Index(body, "name=\"_csrf\" value=\"")
	if idx == -1 {
		t.Fatal("Could not find CSRF token in admin login page")
	}
	sub := body[idx+len("name=\"_csrf\" value=\""):]
	endIdx := strings.Index(sub, "\"")
	csrfToken := sub[:endIdx]

	var csrfCookie string
	for _, c := range rec.Result().Cookies() {
		if c.Name == "_csrf" {
			csrfCookie = c.String()
		}
	}

	// 2. Admin Login Action (POST /admin/login)
	// Default code in test is "agbalumo2024" as defined in the test config or code
	form := strings.NewReader("code=agbalumo2024&_csrf=" + csrfToken)
	req = httptest.NewRequest(http.MethodPost, "/admin/login", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie+"; "+csrfCookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusFound, rec.Code) // Redirects to /admin

	// Now we have an admin session
	// 3. Admin Dashboard
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Ensure we have a valid CSRF token in dashboard to use for POST requests
	body = rec.Body.String()
	idx = strings.Index(body, "_csrf")
	if idx != -1 {
		// just re-parse it loosely or we can just keep the old one, since the session hasn't regenerated CSRF.
		// Wait, the CSRF middleware usually keeps the same token for the cookie session.
	}

	// 4. User Management
	req = httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 5. Listing Management
	req = httptest.NewRequest(http.MethodGet, "/admin/listings", nil)
	req.Header.Set("Cookie", cookie)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Create a dummy listing directly to test approve/reject/hard delete/featured
	// Extracting an ID from /admin/listings might work, but it might be empty if not seeded yet.
	// Oh wait, seeder runs! We can just grab the first listing from /admin/listings.
	body = rec.Body.String()
	idx = strings.Index(body, "<th>")
	// If it has listings, let's find `hx-post="/admin/listings/`
	idx = strings.Index(body, "hx-post=\"/admin/listings/")
	if idx != -1 {
		subStr := body[idx+len("hx-post=\"/admin/listings/"):]
		endId := strings.Index(subStr, "/")
		if endId != -1 {
			listingID := subStr[:endId]

			// 6. Toggle Featured
			req = httptest.NewRequest(http.MethodPost, "/admin/listings/"+listingID+"/featured", nil)
			req.Header.Set("Cookie", cookie+"; "+csrfCookie)
			req.Header.Set("X-CSRF-Token", csrfToken)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)

			// 7. Hard Delete Confirmation View
			req = httptest.NewRequest(http.MethodGet, "/admin/listings/delete-confirm?id="+listingID, nil)
			req.Header.Set("Cookie", cookie)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)

			// 8. Hard Delete Action
			form = strings.NewReader("id=" + listingID + "&_csrf=" + csrfToken)
			req = httptest.NewRequest(http.MethodPost, "/admin/listings/delete", form)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Cookie", cookie+"; "+csrfCookie)
			req.Header.Set("X-CSRF-Token", csrfToken)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code) // usually returns HX-Redirect or 200 snippet
		}
	}

	// 11. Bulk Actions (Approve multiple)
	form = strings.NewReader("action=approve&listing_ids=first-id&listing_ids=second-id&_csrf=" + csrfToken)
	req = httptest.NewRequest(http.MethodPost, "/admin/listings/bulk", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie+"; "+csrfCookie)
	req.Header.Set("X-CSRF-Token", csrfToken)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusFound, rec.Code)

	// 12. Bulk Upload (Testing access, sending empty/invalid CSV for 200 return)
	// Using multipart form requires a bit more setup, but an empty POST form returns 400 or 200 with error
	req = httptest.NewRequest(http.MethodPost, "/admin/upload", nil)
	req.Header.Set("Cookie", cookie+"; "+csrfCookie)
	req.Header.Set("X-CSRF-Token", csrfToken)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusFound, rec.Code)
}
