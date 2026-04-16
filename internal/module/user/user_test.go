package user

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// testRenderer is a minimal echo.Renderer for unit tests that
// need ui.RespondErrorMsg to write a response without a real template engine.
type testRenderer struct{}

func (r *testRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	_, err := fmt.Fprintf(w, "<%s>", name)
	return err
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Case 1: No user
	user, ok := GetUser(c)
	assert.False(t, ok)
	assert.Nil(t, user)

	// Case 2: User pointer present
	mockUser := &domain.User{ID: "u1"}
	c.Set("User", mockUser)
	user, ok = GetUser(c)
	assert.True(t, ok)
	assert.Equal(t, mockUser, user)

	// Case 3: User value present
	valUser := domain.User{ID: "u2"}
	c.Set("User", valUser)
	user, ok = GetUser(c)
	assert.True(t, ok)
	assert.Equal(t, "u2", user.ID)

	// Case 4: Invalid type
	c.Set("User", "not a user")
	_, ok = GetUser(c)
	assert.False(t, ok)
}

func TestMustUser(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Case 1: Panic on no user
	assert.Panics(t, func() {
		MustUser(c)
	})

	// Case 2: Return user when present
	mockUser := &domain.User{ID: "u1"}
	c.Set("User", mockUser)
	user := MustUser(c)
	assert.Equal(t, mockUser, user)
}

func TestRequireUserAPI(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testRenderer{}

	t.Run("no user returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		u, _ := RequireUserAPI(c)
		assert.Nil(t, u)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("user present returns user and no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mockUser := &domain.User{ID: "u1"}
		c.Set("User", mockUser)

		u, err := RequireUserAPI(c)
		assert.NoError(t, err)
		assert.Equal(t, mockUser, u)
	})
}
