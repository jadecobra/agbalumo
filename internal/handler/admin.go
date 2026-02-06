package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	Repo domain.ListingRepository
}

func NewAdminHandler(repo domain.ListingRepository) *AdminHandler {
	return &AdminHandler{Repo: repo}
}

// Middleware
func (h *AdminHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("admin_session")
		if err != nil || cookie.Value != "authenticated" {
			return c.Redirect(http.StatusFound, "/admin/login")
		}
		return next(c)
	}
}

// GET /admin/login
func (h *AdminHandler) HandleLoginView(c echo.Context) error {
	return c.Render(http.StatusOK, "admin_login.html", nil)
}

// POST /admin/login
func (h *AdminHandler) HandleLoginAction(c echo.Context) error {
	code := c.FormValue("code")
	// Externalize Admin Code
	adminCode := os.Getenv("ADMIN_ACCESS_CODE")
	if adminCode == "" {
		adminCode = "agbalumo2024" // Fallback (or log warning)
	}

	if code == adminCode {
		cookie := new(http.Cookie)
		cookie.Name = "admin_session"
		cookie.Value = "authenticated"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = c.Scheme() == "https" || os.Getenv("AGBALUMO_ENV") == "production"
		cookie.SameSite = http.SameSiteLaxMode
		c.SetCookie(cookie)
		return c.Redirect(http.StatusFound, "/admin")
	}

	return c.Render(http.StatusOK, "admin_login.html", map[string]interface{}{
		"Error": "Invalid Access Code",
	})
}

// GET /admin
func (h *AdminHandler) HandleDashboard(c echo.Context) error {
	listings, err := h.Repo.FindAll(c.Request().Context(), "", "", true)
	if err != nil {
		return RespondError(c, err)
	}
	return c.Render(http.StatusOK, "admin_dashboard.html", map[string]interface{}{
		"Listings": listings,
	})
}

// DELETE /admin/listings/:id
func (h *AdminHandler) HandleDelete(c echo.Context) error {
	id := c.Param("id")
	l, err := h.Repo.FindByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "Listing not found")
	}
	l.IsActive = false
	if err := h.Repo.Save(c.Request().Context(), l); err != nil {
		return RespondError(c, err)
	}

	return c.String(http.StatusOK, "")
}
