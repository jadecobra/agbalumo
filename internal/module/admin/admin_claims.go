package admin

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

const errClaimRequestNotFound = "Claim request not found"

// HandleApproveClaim approves a user's claim request and transfers listing ownership.
func (h *AdminHandler) HandleApproveClaim(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.App.DB.UpdateClaimRequestStatus(ctx, id, domain.ClaimStatusApproved); err != nil {
		return c.String(http.StatusNotFound, errClaimRequestNotFound)
	}

	return c.NoContent(http.StatusOK)
}

// HandleRejectClaim rejects a user's claim request.
func (h *AdminHandler) HandleRejectClaim(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.App.DB.UpdateClaimRequestStatus(ctx, id, domain.ClaimStatusRejected); err != nil {
		return c.String(http.StatusNotFound, errClaimRequestNotFound)
	}

	return c.NoContent(http.StatusOK)
}
