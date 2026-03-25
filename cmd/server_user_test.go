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
	form := strings.NewReader("title=" + title + "&type=Service&owner_origin=Nigeria&description=Testing&contact_email=dev@test.com&city=Lagos&_csrf=" + csrfToken)
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
			form = strings.NewReader("title=" + title + "+Updated&type=Service&owner_origin=Nigeria&description=Testing&contact_email=dev@test.com&city=Lagos&_csrf=" + csrfToken)
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
