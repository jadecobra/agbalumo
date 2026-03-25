package cmd_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
