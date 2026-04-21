package admin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// HandleAddCategory processes the form submission to add a new category
func (h *AdminHandler) HandleAddCategory(c echo.Context) error {
	ctx := c.Request().Context()
	name := strings.TrimSpace(c.FormValue(domain.FieldName))
	if name == "" {
		return c.Redirect(http.StatusFound, domain.PathAdmin)
	}

	if existing, err := h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{ActiveOnly: false}); err == nil {
		if hasDuplicateCategory(existing, name) {
			return h.redirectWithFlash(c, fmt.Sprintf("Category '%s' already exists!", name), domain.PathAdmin)
		}
	}

	claimable := c.FormValue(domain.FieldClaimable) == "true"
	now := time.Now()
	cat := domain.CategoryData{
		ID:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Name:      name,
		Claimable: claimable,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.App.DB.SaveCategory(ctx, cat); err != nil {
		c.Logger().Errorf("failed to save custom category: %v", err)
	}

	return h.redirectWithFlash(c, "Category added successfully!", domain.PathAdmin)
}

func hasDuplicateCategory(existing []domain.CategoryData, name string) bool {
	for _, cat := range existing {
		if strings.EqualFold(cat.Name, name) {
			return true
		}
	}
	return false
}
