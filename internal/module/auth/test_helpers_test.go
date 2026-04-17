package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module/auth"
	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// MockGoogleProvider
type MockGoogleProvider struct {
	testifyMock.Mock
}

func (m *MockGoogleProvider) GetAuthCodeURL(state string, scheme string, host string) string {
	args := m.Called(state, scheme, host)
	return args.String(0)
}

func (m *MockGoogleProvider) Exchange(ctx context.Context, code string, scheme string, host string) (*oauth2.Token, error) {
	args := m.Called(ctx, code, scheme, host)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockGoogleProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*auth.GoogleUser, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.GoogleUser), args.Error(1)
}

func setupAuthContext(method, url string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}
	req := httptest.NewRequest(method, url, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	store := sessions.NewCookieStore([]byte("secret"))
	sess, _ := store.Get(req, "session-name")
	c.Set("session", sess)

	return c, rec
}

func performRegistration(t *testing.T, app *env.AppEnv, payload map[string]string) *httptest.ResponseRecorder {
	c, rec := setupAuthContext(http.MethodGet, "/auth/google/callback?state=random-state&code=valid-code")
	req := c.Request()
	req.AddCookie(&http.Cookie{Name: domain.SessionKeyOAuthState, Value: "random-state"})


	mockProvider := &MockGoogleProvider{}
	h := auth.NewAuthHandler(app)
	h.GoogleProvider = mockProvider

	token := &oauth2.Token{AccessToken: "token"}
	gUser := &auth.GoogleUser{
		ID:      payload["id"],
		Email:   payload["email"],
		Name:    payload["name"],
		Picture: payload["picture"],
	}

	mockProvider.On("Exchange", testifyMock.Anything, "valid-code", "http", "example.com").Return(token, nil)
	mockProvider.On("GetUserInfo", testifyMock.Anything, token).Return(gUser, nil)

	_ = h.GoogleCallback(c)

	return rec
}
