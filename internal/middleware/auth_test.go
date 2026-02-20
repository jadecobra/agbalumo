package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAuthMiddleware_RequireAuth_Redirect(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	authMw := middleware.NewAuthMiddleware(mockRepo)

	handlerFunc := authMw.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "/auth/google/login", rec.Header().Get("Location"))
}

func TestAuthMiddleware_RequireAuth_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["user_id"] = "user-123"
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	authMw := middleware.NewAuthMiddleware(mockRepo)

	handlerFunc := authMw.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Success", rec.Body.String())
}

func TestAuthMiddleware_OptionalAuth_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	sess.Values["user_id"] = "user-123"
	c.Set("session", sess)

	mockRepo := &mock.MockListingRepository{}
	user := domain.User{ID: "user-123"}
	mockRepo.On("FindUserByID", testifyMock.Anything, "user-123").Return(user, nil)

	authMw := middleware.NewAuthMiddleware(mockRepo)

	handlerFunc := authMw.OptionalAuth(func(c echo.Context) error {
		u := c.Get("User")
		if u == nil {
			return c.String(http.StatusOK, "No User")
		}
		return c.String(http.StatusOK, "Has User")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, "Has User", rec.Body.String())
}

func TestAuthMiddleware_OptionalAuth_NoSession(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := &mock.MockListingRepository{}
	authMw := middleware.NewAuthMiddleware(mockRepo)

	handlerFunc := authMw.OptionalAuth(func(c echo.Context) error {
		u := c.Get("User")
		if u == nil {
			return c.String(http.StatusOK, "No User")
		}
		return c.String(http.StatusOK, "Has User")
	})

	err := handlerFunc(c)

	assert.NoError(t, err)
	assert.Equal(t, "No User", rec.Body.String())
}
