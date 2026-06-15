package admin

import (
	"net/http"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

// HandleResolveFeedback resolves a user feedback entry.
func (h *AdminHandler) HandleResolveFeedback(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.App.DB.ResolveFeedback(ctx, id); err != nil {
		return ui.RespondError(c, err)
	}

	feedbacks, err := h.App.DB.GetAllFeedback(ctx)
	if err != nil {
		return ui.RespondError(c, err)
	}

	var targetFeedback *domain.Feedback
	for i := range feedbacks {
		if feedbacks[i].ID == id {
			targetFeedback = &feedbacks[i]
			break
		}
	}

	if targetFeedback == nil {
		return c.String(http.StatusNotFound, "Feedback not found")
	}

	return h.RenderTyped(c, "admin_feedback_item", targetFeedback)
}
