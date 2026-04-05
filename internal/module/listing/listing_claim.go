package listing

import (
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleClaim processes a request to claim an unowned listing.
func (h *ListingHandler) HandleClaim(c echo.Context) error {
	u, ok := user.GetUser(c)
	if !ok {
		return c.Redirect(http.StatusFound, "/auth/google/login")
	}

	id := c.Param("id")

	_, err := h.App.ListingSvc.ClaimListing(c.Request().Context(), *u, id)
	if err != nil {
		switch err.Error() {
		case "listing not found":
			return ui.RespondError(c, echo.NewHTTPError(http.StatusNotFound, "Listing not found"))
		case "listing is already owned", "listing type cannot be claimed":
			return ui.RespondError(c, echo.NewHTTPError(http.StatusForbidden, err.Error()))
		case "you already have a pending claim for this listing":
			return ui.RespondError(c, echo.NewHTTPError(http.StatusConflict, err.Error()))
		default:
			return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit claim: "+err.Error()))
		}
	}

	// HTMX: replace the claim button with a pending-approval notice
	c.Response().Header().Set("Content-Type", "text/html")
	return c.HTML(http.StatusOK, `
		<div class="flex items-center gap-2 bg-earth-accent/10 border border-earth-accent/20 px-3 py-1.5">
			<span class="material-symbols-outlined text-[14px] text-earth-accent">pending</span>
			<span class="text-earth-accent text-xs font-bold uppercase tracking-widest">Claim Pending Review</span>
		</div>`)
}
