package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

// -- Google Interaction Abstraction --

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GoogleProvider interface {
	GetAuthCodeURL(state string, scheme string, host string) string
	Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error)
}

type RealGoogleProvider struct {
	config      *oauth2.Config
	UserInfoURL string
}

func NewRealGoogleProvider() *RealGoogleProvider {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Fprintf(os.Stderr, "WARNING: GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET not set. OAuth will fail.\n")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return &RealGoogleProvider{
		config:      config,
		UserInfoURL: googleUserInfoURL,
	}
}

func (p *RealGoogleProvider) getRedirectURL(scheme string, host string) string {
	// 1. Prefer BASE_URL if set (must include scheme and host)
	baseURL := os.Getenv("BASE_URL")
	if baseURL != "" {
		return fmt.Sprintf("%s/auth/google/callback", baseURL)
	}

	// 2. Fallback to GOOGLE_REDIRECT_URL (Explicit)
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL != "" {
		return redirectURL
	}

	// 3. Robust dynamic host/scheme detection
	if scheme == "" {
		scheme = "http"
		// Secure by default if port matches standard TLS or in production
		if strings.HasSuffix(host, ":8443") || os.Getenv("AGBALUMO_ENV") == "production" {
			scheme = "https"
		}
	}

	generated := fmt.Sprintf("%s://%s/auth/google/callback", scheme, host)
	fmt.Printf("[DEBUG] OAuth Redirect URI: %s\n", generated)
	return generated
}

func (p *RealGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	// Create a copy of the config with the dynamic redirect URL
	cfg := *p.config
	cfg.RedirectURL = p.getRedirectURL(scheme, host)
	return cfg.AuthCodeURL(state)
}

func (p *RealGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
	// Create a copy of the config with the dynamic redirect URL
	cfg := *p.config
	cfg.RedirectURL = p.getRedirectURL(scheme, host)
	return cfg.Exchange(ctx, code)
}

func (p *RealGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := p.config.Client(ctx, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	// #nosec G107 G704 - SSRF check: The URL is from the provider's configuration (defaulting to a constant) and AccessToken is handled by the authenticated client.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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

// -- Mock Google Provider for Browser Audits --

type MockGoogleProvider struct {
	Email string
	Name  string
}

func (p *MockGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	if scheme == "" {
		scheme = "http"
	}
	// Pure redirect to callback to simulate the flow
	return fmt.Sprintf("%s://%s/auth/google/callback?state=%s&code=mock-code", scheme, host, state)
}

func (p *MockGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "mock-token"}, nil // #nosec G101 - Rationale: Mock provider for testing/audits. Mock token is intentional and non-sensitive.
}

func (p *MockGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	return &GoogleUser{
		ID:      "mock-" + p.Email,
		Email:   p.Email,
		Name:    p.Name,
		Picture: "https://ui-avatars.com/api/?name=Mock+User",
	}, nil
}

// -- Auth Handler --

type AuthHandler struct {
	Repo           domain.UserStore
	GoogleProvider GoogleProvider
	Cfg            *config.Config
}

func NewAuthHandler(deps AuthDependencies) *AuthHandler {
	if deps.GoogleProvider == nil {
		deps.GoogleProvider = NewRealGoogleProvider()
	}
	return &AuthHandler{
		Repo:           deps.UserStore,
		GoogleProvider: deps.GoogleProvider,
		Cfg:            deps.Config,
	}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	e.GET("/auth/dev", h.DevLogin)
	e.GET("/auth/logout", h.Logout)
	e.GET("/auth/google/login", h.GoogleLogin)
	e.GET("/auth/google/callback", h.GoogleCallback)
}

// DevLogin simulates a Google Login for development
// param: email
func (h *AuthHandler) DevLogin(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		email = h.Cfg.DevAuthEmail
	}

	// Environment Check: Only allow in development
	if h.Cfg.Env != "development" {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusForbidden, "Dev login disabled in production"))
	}
	// Simulate "Google ID" creation
	googleID := "dev-" + email
	name := "Dev User"
	avatar := "https://ui-avatars.com/api/?name=Dev+User&background=random"

	user, err := h.findOrCreateUser(c.Request().Context(), googleID, email, name, avatar)
	if err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to login"))
	}

	return h.setSessionAndRedirect(c, user.ID)
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	if !h.Cfg.HasGoogleAuth {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusServiceUnavailable, "Google OAuth is not configured. Please set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET in .env"))
	}

	state := uuid.New().String()
	baseURL := os.Getenv("BASE_URL")
	isSecure := h.Cfg.Env == "production" || strings.HasPrefix(baseURL, "https://")

	cookie := new(http.Cookie)
	cookie.Name = "oauth_state"
	cookie.Value = state
	cookie.Path = "/"
	cookie.MaxAge = 10 * 60 // 10 minutes
	cookie.HttpOnly = true
	cookie.Secure = c.Scheme() == "https" || isSecure
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	url := h.GoogleProvider.GetAuthCodeURL(state, c.Scheme(), c.Request().Host)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")

	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "States don't match or expired"))
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
		return handler.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Code exchange failed"))
	}

	gUser, err := h.GoogleProvider.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "User data fetch failed"))
	}

	user, err := h.findOrCreateUser(c.Request().Context(), gUser.ID, gUser.Email, gUser.Name, gUser.Picture)
	if err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to login"))
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
		return handler.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Session Store Missing"))
	}

	baseURL := os.Getenv("BASE_URL")
	isSecure := h.Cfg.Env == "production" || strings.HasPrefix(baseURL, "https://")

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   c.Scheme() == "https" || isSecure,
		SameSite: http.SameSiteLaxMode,
	}
	sess.Values["user_id"] = userID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return handler.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to save session"))
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	sess := customMiddleware.GetSession(c)
	if sess != nil {
		sess.Options.MaxAge = -1
		_ = sess.Save(c.Request(), c.Response())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
