package listing

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

func (h *ListingHandler) HandleSaveToggle(c echo.Context) error {
	u, err := user.RequireUserAPI(c)
	if err != nil {
		return err
	}
	id := c.Param("id")
	ctx := c.Request().Context()

	saved, err := h.App.DB.IsListingSaved(ctx, u.ID, id)
	if err != nil {
		return ui.RespondError(c, err)
	}

	if saved {
		err = h.App.DB.UnsaveListing(ctx, u.ID, id)
	} else {
		err = h.App.DB.SaveListing(ctx, u.ID, id)
	}
	if err != nil {
		return ui.RespondError(c, err)
	}

	return c.Render(http.StatusOK, "save_button", map[string]interface{}{
		"ListingID": id,
		"IsSaved":   !saved,
	})
}
