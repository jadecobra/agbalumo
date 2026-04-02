package auth

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

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
