package auth

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	App            *env.AppEnv
	GoogleProvider GoogleProvider
}

func NewAuthHandler(app *env.AppEnv) *AuthHandler {
	var googleProvider GoogleProvider
	if app.Cfg.MockAuth {
		googleProvider = &MockGoogleProvider{
			Email: "test@agbalumo.com",
			Name:  "Test User",
		}
	} else {
		googleProvider = NewRealGoogleProvider()
	}

	return &AuthHandler{
		App:            app,
		GoogleProvider: googleProvider,
	}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	e.GET("/auth/dev", h.DevLogin)
	e.GET("/auth/logout", h.Logout)
	e.GET("/auth/google/login", h.GoogleLogin)
	e.GET("/auth/google/callback", h.GoogleCallback)
}

type MockGoogleProvider struct {
	Email string
	Name  string
}

func (p *MockGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + host + "/auth/google/callback?state=" + state + "&code=mock-code"
}

func (p *MockGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "mock-token"}, nil // #nosec - testing only
}

func (p *MockGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	return &GoogleUser{
		ID:    "mock-" + p.Email,
		Email: p.Email,
		Name:  p.Name,
	}, nil
}
