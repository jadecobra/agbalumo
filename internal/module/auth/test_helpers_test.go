package auth_test

import (
	"context"
	"net/http/httptest"

	"github.com/gorilla/sessions"
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
