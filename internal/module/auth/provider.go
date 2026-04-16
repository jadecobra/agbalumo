package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

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
	clientID := os.Getenv(domain.EnvKeyGoogleClientID)
	clientSecret := os.Getenv(domain.EnvKeyGoogleClientSecret)

	if clientID == "" || clientSecret == "" {
		fmt.Fprintf(os.Stderr, "WARNING: %s or %s not set. OAuth will fail.\n", domain.EnvKeyGoogleClientID, domain.EnvKeyGoogleClientSecret)
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
	baseURL := os.Getenv(config.EnvBaseURL)
	if baseURL != "" {
		return fmt.Sprintf("%s/auth/google/callback", baseURL)
	}

	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if redirectURL != "" {
		return redirectURL
	}

	if scheme == "" {
		scheme = "http"
		if strings.HasSuffix(host, ":8443") || os.Getenv("AGBALUMO_ENV") == "production" {
			scheme = "https"
		}
	}

	return fmt.Sprintf("%s://%s/auth/google/callback", scheme, host)
}

func (p *RealGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	cfg := *p.config
	cfg.RedirectURL = p.getRedirectURL(scheme, host)
	return cfg.AuthCodeURL(state)
}

func (p *RealGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
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

	resp, err := client.Do(req) // #nosec G704 - Standard Google OAuth2 userinfo request
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
