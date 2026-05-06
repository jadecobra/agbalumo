package module

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
)

func TestPopulateBase(t *testing.T) {
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	// Use real categorization service instead of mock to test database integration
	env.App.CategorizationSvc = service.NewCategorizationService(env.App.DB, &domain.CategoryCache{})

	// Seed some categories so Categories slice is not empty
	cat := domain.CategoryData{ID: "c1", Name: "Test Cat", Active: true}
	err := env.App.DB.SaveCategory(context.Background(), cat)
	if err != nil {
		t.Fatalf("Failed to save category: %v", err)
	}

	h := &BaseHandler{App: env.App}

	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// PopulateBase should return BaseViewData
	data := h.PopulateBase(c)

	if data.Env == "" {
		t.Error("expected Env to be non-empty")
	}
	if len(data.Categories) == 0 {
		t.Error("expected Categories to be non-empty")
	}
}
