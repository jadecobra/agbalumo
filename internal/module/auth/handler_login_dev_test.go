package auth_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_DevLogin_Production(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	req := httptest.NewRequest(http.MethodGet, "/auth/dev", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	cfg := config.LoadConfig()
	cfg.Env = "production"
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: cfg})

	err := h.DevLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthHandler_DevLogin_Success(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	store := sessions.NewCookieStore([]byte("secret"))
	req := httptest.NewRequest(http.MethodGet, "/auth/dev?email=test@dev.com", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	repo := handler.SetupTestRepository(t)
	h := auth.NewAuthHandler(auth.AuthDependencies{UserStore: repo, Config: config.LoadConfig()})

	_ = os.Setenv("AGBALUMO_ENV", "development")
	defer func() { _ = os.Unsetenv("AGBALUMO_ENV") }()

	err := h.DevLogin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.NotEmpty(t, sess.Values["user_id"])
}
