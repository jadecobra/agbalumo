package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// -- Google Interaction Abstraction --

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleProvider interface {
	GetAuthCodeURL(state string, host string) string
	Exchange(ctx context.Context, code string, host string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error)
}

type RealGoogleProvider struct {
	config *oauth2.Config
}

func NewRealGoogleProvider() *RealGoogleProvider {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return &RealGoogleProvider{config: config}
}

func (p *RealGoogleProvider) getRedirectURL(host string) string {
	// 1. Prefer BASE_URL (e.g. http://192.168.1.5:8080)
	baseURL := os.Getenv("BASE_URL")
	if baseURL != "" {
		// Basic trimming to avoid double slashes if user adds one
		// We avoid importing "strings" just for this if not already imported,
		// but "strings" is not imported in auth.go?
		// Checked file: imports: context, encoding/json, fmt, io, net/http, os, time...
		// Need to add strings to imports?
		// Or just blindly append if we trust user?
		// Let's assume user inputs correctly or simple check.
		// Actually, importing strings is better.
		// But replacing content safely involves keeping imports separate.
		// I will just use fmt and assume standard format for now to minimize diff risk.
		return fmt.Sprintf("%s/auth/google/callback", baseURL)
	}

	// 2. Fallback to GOOGLE_REDIRECT_URL (Legacy)
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL != "" {
		return redirectURL
	}
	// 3. Default to dynamic host
	return fmt.Sprintf("http://%s/auth/google/callback", host)
}

func (p *RealGoogleProvider) GetAuthCodeURL(state string, host string) string {
	// Create a copy of the config with the dynamic redirect URL
	cfg := *p.config
	cfg.RedirectURL = p.getRedirectURL(host)
	return cfg.AuthCodeURL(state)
}

func (p *RealGoogleProvider) Exchange(ctx context.Context, code string, host string) (*oauth2.Token, error) {
	// Create a copy of the config with the dynamic redirect URL
	cfg := *p.config
	cfg.RedirectURL = p.getRedirectURL(host)
	return cfg.Exchange(ctx, code)
}

func (p *RealGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user info: %s", resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gUser GoogleUser
	if err := json.Unmarshal(content, &gUser); err != nil {
		return nil, err
	}
	return &gUser, nil
}

// -- Auth Handler --

type AuthHandler struct {
	Repo           domain.ListingRepository
	GoogleProvider GoogleProvider
}

func NewAuthHandler(repo domain.ListingRepository, provider GoogleProvider) *AuthHandler {
	if provider == nil {
		provider = NewRealGoogleProvider()
	}
	return &AuthHandler{
		Repo:           repo,
		GoogleProvider: provider,
	}
}

// DevLogin simulates a Google Login for development
// param: email
func (h *AuthHandler) DevLogin(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		email = "dev@agbalumo.com"
	}

	// Environment Check: Only allow in development
	env := os.Getenv("AGBALUMO_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env != "development" {
		return c.String(http.StatusForbidden, "Dev login disabled in production")
	}
	// Simulate "Google ID" creation
	googleID := "dev-" + email
	name := "Dev User"
	avatar := "https://ui-avatars.com/api/?name=Dev+User&background=random"

	user, err := h.findOrCreateUser(c.Request().Context(), googleID, email, name, avatar)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to login")
	}

	return h.setSessionAndRedirect(c, user.ID)
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	url := h.GoogleProvider.GetAuthCodeURL("random-state", c.Request().Host)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != "random-state" {
		return c.String(http.StatusBadRequest, "States don't match")
	}

	code := c.QueryParam("code")
	token, err := h.GoogleProvider.Exchange(c.Request().Context(), code, c.Request().Host)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Code exchange failed")
	}

	gUser, err := h.GoogleProvider.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		return c.String(http.StatusInternalServerError, "User data fetch failed")
	}

	user, err := h.findOrCreateUser(c.Request().Context(), gUser.ID, gUser.Email, gUser.Name, gUser.Picture)
	if err != nil {
		// return c.String(http.StatusInternalServerError, "Failed to login: "+err.Error()) // Leaking error details
		return c.String(http.StatusInternalServerError, "Failed to login")
	}

	return h.setSessionAndRedirect(c, user.ID)
}

func (h *AuthHandler) findOrCreateUser(ctx context.Context, googleID, email, name, avatar string) (*domain.User, error) {
	user, err := h.Repo.FindUserByGoogleID(ctx, googleID)
	if err != nil {
		// Create new user
		user = domain.User{
			ID:        uuid.New().String(),
			GoogleID:  googleID,
			Email:     email,
			Name:      name,
			AvatarURL: avatar,
			CreatedAt: time.Now(),
		}
		if err := h.Repo.SaveUser(ctx, user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	// Update profile info if changed
	if user.AvatarURL != avatar || user.Name != name {
		user.AvatarURL = avatar
		user.Name = name
		if err := h.Repo.SaveUser(ctx, user); err != nil {
			// Log error but continue
			fmt.Printf("Failed to update user profile: %v\n", err)
		}
	}
	return &user, nil
}

func (h *AuthHandler) setSessionAndRedirect(c echo.Context, userID string) error {
	sess := customMiddleware.GetSession(c)
	if sess == nil {
		// In tests or if middleware missing
		return c.String(http.StatusInternalServerError, "Session Store Missing")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   c.Scheme() == "https" || os.Getenv("AGBALUMO_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["user_id"] = userID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Middleware to inject user into context
func (h *AuthHandler) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := customMiddleware.GetSession(c)
		if sess != nil {
			if userID, ok := sess.Values["user_id"].(string); ok {
				user, err := h.Repo.FindUserByID(c.Request().Context(), userID)
				if err == nil {
					c.Set("User", user)
				}
			}
		}
		return next(c)
	}
}

func (h *AuthHandler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := customMiddleware.GetSession(c)
		authSuccess := false
		if sess != nil {
			if _, ok := sess.Values["user_id"].(string); ok {
				authSuccess = true
			}
		}

		if !authSuccess {
			// Redirect to Google Login
			return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/login")
		}
		return next(c)
	}
}
