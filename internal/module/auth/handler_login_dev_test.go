package auth_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_DevLogin_Production(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/dev", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	app.Cfg.Env = "production"
	h := auth.NewAuthHandler(app)

	err := h.DevLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthHandler_DevLogin_Success(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=test@dev.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := auth.NewAuthHandler(app)

	_ = os.Setenv("AGBALUMO_ENV", "development")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	err := h.DevLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.NotEmpty(t, sess.Values[domain.SessionKeyUserID])
}
