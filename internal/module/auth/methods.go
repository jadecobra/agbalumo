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
	email := c.QueryParam("email")
	if email == "" {
		email = h.App.Cfg.DevAuthEmail
	}

	if h.App.Cfg.Env != "development" {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusForbidden, "Dev login disabled in production"))
	}

	googleID := "dev-" + email
	name := "Dev User"
	avatar := "https://ui-avatars.com/api/?name=Dev+User&background=random"

	user, err := h.findOrCreateUser(c.Request().Context(), googleID, email, name, avatar)
	if err != nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, domain.MsgFailedToLogin))
	}

	return h.setSessionAndRedirect(c, user.ID)
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	if !h.App.Cfg.HasGoogleAuth {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusServiceUnavailable, "Google OAuth is not configured"))
	}

	state := uuid.New().String()
	baseURL := os.Getenv("BASE_URL")
	isSecure := h.App.Cfg.Env == "production" || strings.HasPrefix(baseURL, "https://")

	cookie := new(http.Cookie)
	cookie.Name = "oauth_state"
	cookie.Value = state
	cookie.Path = "/"
	cookie.MaxAge = 10 * 60
	cookie.HttpOnly = true
	cookie.Secure = c.Scheme() == "https" || isSecure
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
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Session Store Missing"))
	}

	baseURL := os.Getenv("BASE_URL")
	isSecure := h.App.Cfg.Env == "production" || strings.HasPrefix(baseURL, "https://")

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   c.Scheme() == "https" || isSecure,
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["user_id"] = userID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session"))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")

	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "States don't match or expired"))
	}

	deleteCookie := new(http.Cookie)
	deleteCookie.Name = "oauth_state"
	deleteCookie.Value = ""
	deleteCookie.Path = "/"
	deleteCookie.MaxAge = -1
	c.SetCookie(deleteCookie)

	code := c.QueryParam("code")
	token, err := h.GoogleProvider.Exchange(c.Request().Context(), code, c.Scheme(), c.Request().Host)
	if err != nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Code exchange failed"))
	}

	gUser, err := h.GoogleProvider.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "User data fetch failed"))
	}

	user, err := h.findOrCreateUser(c.Request().Context(), gUser.ID, gUser.Email, gUser.Name, gUser.Picture)
	if err != nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, domain.MsgFailedToLogin))
	}

	return h.setSessionAndRedirect(c, user.ID)
}
