package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) DevLogin(c echo.Context) error {
	email := c.QueryParam(domain.FieldEmail)
	if email == "" {
		email = h.App.Cfg.DevAuthEmail
	}

	if h.App.Cfg.Env != domain.EnvDevelopment {
		return ui.RespondErrorMsg(c, http.StatusForbidden, "Dev login disabled in production")
	}

	googleID := "dev-" + email
	name := "Dev User"
	avatar := "https://ui-avatars.com/api/?name=Dev+User&background=random"

	return h.loginWith(c, googleID, email, name, avatar)
}

// loginWith looks up or creates the user then sets the session and redirects.
func (h *AuthHandler) loginWith(c echo.Context, googleID, email, name, avatar string) error {
	u, err := h.findOrCreateUser(c.Request().Context(), googleID, email, name, avatar)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, domain.MsgFailedToLogin)
	}
	return h.setSessionAndRedirect(c, u.ID)
}

func isSecureCookie(c echo.Context, env string) bool {
	baseURL := os.Getenv("BASE_URL")
	return env == domain.EnvProduction || strings.HasPrefix(baseURL, "https://") || c.Scheme() == "https"
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	if !h.App.Cfg.HasGoogleAuth {
		return ui.RespondErrorMsg(c, http.StatusServiceUnavailable, "Google OAuth is not configured")
	}

	state := uuid.New().String()

	cookie := new(http.Cookie)
	cookie.Name = domain.SessionKeyOAuthState
	cookie.Value = state
	cookie.Path = "/"
	cookie.MaxAge = 10 * 60
	cookie.HttpOnly = true
	cookie.Secure = isSecureCookie(c, h.App.Cfg.Env)
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	url := h.GoogleProvider.GetAuthCodeURL(state, c.Scheme(), c.Request().Host)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.Options.MaxAge = -1
		_ = sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *AuthHandler) findOrCreateUser(ctx context.Context, googleID, email, name, avatar string) (*domain.User, error) {
	user, err := h.App.DB.FindUserByGoogleID(ctx, googleID)
	if err != nil {
		user = domain.User{
			ID:        uuid.New().String(),
			GoogleID:  googleID,
			Email:     email,
			Name:      name,
			AvatarURL: avatar,
			CreatedAt: time.Now(),
		}
		if err := h.App.DB.SaveUser(ctx, user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	if user.AvatarURL != avatar || user.Name != name {
		user.AvatarURL = avatar
		user.Name = name
		_ = h.App.DB.SaveUser(ctx, user)
	}
	return &user, nil
}

func (h *AuthHandler) setSessionAndRedirect(c echo.Context, userID string) error {
	sess := customMiddleware.GetSession(c)
	if sess == nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, "Session Store Missing")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   isSecureCookie(c, h.App.Cfg.Env),
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values[domain.SessionKeyUserID] = userID

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, "Failed to save session")
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	state := c.QueryParam(domain.ParamState)

	stateCookie, err := c.Cookie(domain.SessionKeyOAuthState)
	if err != nil || stateCookie.Value != state {
		return ui.RespondErrorMsg(c, http.StatusBadRequest, "States don't match or expired")
	}

	deleteCookie := new(http.Cookie)
	deleteCookie.Name = domain.SessionKeyOAuthState
	deleteCookie.Value = ""
	deleteCookie.Path = "/"
	deleteCookie.MaxAge = -1
	c.SetCookie(deleteCookie)

	code := c.QueryParam(domain.ParamCode)
	token, err := h.GoogleProvider.Exchange(c.Request().Context(), code, c.Scheme(), c.Request().Host)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, "Code exchange failed")
	}

	gUser, err := h.GoogleProvider.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, "User data fetch failed")
	}

	return h.loginWith(c, gUser.ID, gUser.Email, gUser.Name, gUser.Picture)
}
