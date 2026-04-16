package listing

import (
	"github.com/jadecobra/agbalumo/internal/common"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"

	"errors"
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
		if errors.Is(err, domain.ErrListingNotFound) {
			return ui.RespondErrorMsg(c, http.StatusNotFound, domain.ErrListingNotFound.Error())
		}
		if errors.Is(err, domain.ErrListingOwned) || errors.Is(err, domain.ErrListingNotClaimable) {
			return ui.RespondErrorMsg(c, http.StatusForbidden, err.Error())
		}
		if errors.Is(err, domain.ErrPendingClaimExists) {
			return ui.RespondErrorMsg(c, http.StatusConflict, err.Error())
		}
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, "Failed to submit claim: "+err.Error())
	}

	// HTMX: replace the claim button with a pending-approval notice
	c.Response().Header().Set(common.HeaderContentType, common.MimeHTML)
	return c.HTML(http.StatusOK, `
		<div class="flex items-center gap-2 bg-earth-accent/10 border border-earth-accent/20 px-3 py-1.5">
			<span class="material-symbols-outlined text-[14px] text-earth-accent">pending</span>
			<span class="text-earth-accent text-xs font-bold uppercase tracking-widest">Claim Pending Review</span>
		</div>`)
}
