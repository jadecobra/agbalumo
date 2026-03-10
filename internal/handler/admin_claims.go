package handler

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

// HandleApproveClaim approves a user's claim request and transfers listing ownership.
func (h *AdminHandler) HandleApproveClaim(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.Repo.UpdateClaimRequestStatus(ctx, id, domain.ClaimStatusApproved); err != nil {
		return c.String(http.StatusNotFound, "Claim request not found")
	}

	return c.NoContent(http.StatusOK)
}

// HandleRejectClaim rejects a user's claim request.
func (h *AdminHandler) HandleRejectClaim(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.Repo.UpdateClaimRequestStatus(ctx, id, domain.ClaimStatusRejected); err != nil {
		return c.String(http.StatusNotFound, "Claim request not found")
	}

	return c.NoContent(http.StatusOK)
}
