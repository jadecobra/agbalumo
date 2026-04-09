package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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
