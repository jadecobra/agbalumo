package module

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
)

func TestPopulateBase(t *testing.T) {
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	// Seed some categories
	cat := domain.CategoryData{ID: "c1", Name: "Test Cat", Active: true}
	_ = env.App.DB.SaveCategory(context.Background(), cat)

	h := &BaseHandler{App: env.App}

	tests := []struct {
		setupContext func(echo.Context)
		wantUser     *domain.User
		name         string
	}{
		{
			name:         "no user in context",
			setupContext: func(c echo.Context) {},
			wantUser:     nil,
		},
		{
			name: "user in context as pointer",
			setupContext: func(c echo.Context) {
				c.Set("User", &domain.User{ID: "u1", Name: "Pointer User"})
			},
			wantUser: &domain.User{ID: "u1", Name: "Pointer User"},
		},
		{
			name: "user in context as value",
			setupContext: func(c echo.Context) {
				c.Set("User", domain.User{ID: "u2", Name: "Value User"})
			},
			wantUser: &domain.User{ID: "u2", Name: "Value User"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			tt.setupContext(c)

			data := h.PopulateBase(c)

			if tt.wantUser == nil {
				if data.User != nil {
					t.Errorf("expected User to be nil, got %v", data.User)
				}
			} else {
				// This will fail to compile or fail type assertion if BaseViewData.User is interface{}
				// or if it's not correctly handled in PopulateBase
				user, ok := data.User.(*domain.User)
				if !ok {
					t.Errorf("expected User to be *domain.User, got %T", data.User)
					return
				}
				if user.ID != tt.wantUser.ID {
					t.Errorf("expected User ID %s, got %s", tt.wantUser.ID, user.ID)
				}
				if user.Name != tt.wantUser.Name {
					t.Errorf("expected User Name %s, got %s", tt.wantUser.Name, user.Name)
				}
			}
		})
	}
}
