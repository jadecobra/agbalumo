package feedback

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

type FeedbackHandler struct {
	Repo domain.FeedbackStore
}

func NewFeedbackHandler(repo domain.FeedbackStore) *FeedbackHandler {
	return &FeedbackHandler{Repo: repo}
}

// RegisterRoutes registers the feedback routes
func (h *FeedbackHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	feedbackGroup := e.Group("/feedback", authMw.RequireAuth)
	feedbackGroup.GET("/modal", h.HandleModal)
	feedbackGroup.POST("", h.HandleSubmit)
}

// HandleModal renders the feedback modal form
func (h *FeedbackHandler) HandleModal(c echo.Context) error {
	return c.Render(http.StatusOK, "modal_feedback.html", nil)
}

// HandleSubmit processes the feedback form submission
func (h *FeedbackHandler) HandleSubmit(c echo.Context) error {
	u, ok := user.GetUser(c)
	if !ok || u == nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
	}

	contentType := c.QueryParam("type")
	if contentType == "" {
		contentType = c.FormValue("type")
	}

	content := c.FormValue("content")

	// Validate
	if content == "" {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusBadRequest, "Content is required"))
	}
	if contentType == "" {
		contentType = string(domain.FeedbackTypeOther)
	}

	fb := domain.Feedback{
		ID:        uuid.New().String(),
		UserID:    u.ID,
		Type:      domain.FeedbackType(contentType),
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := h.Repo.SaveFeedback(c.Request().Context(), fb); err != nil {
		return ui.RespondError(c, echo.NewHTTPError(http.StatusInternalServerError, "Failed to save feedback"))
	}

	// Return success message or close modal
	return c.HTML(http.StatusOK, `
		<div class="flex flex-col items-center justify-center p-8 space-y-4 text-center animate-in fade-in zoom-in-95 duration-300">
			<div class="h-16 w-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mb-2">
				<span class="material-symbols-outlined text-4xl text-green-600 dark:text-green-400">check_circle</span>
			</div>
			<h3 class="text-xl font-bold text-stone-800 dark:text-white">Thank You!</h3>
			<p class="text-stone-500 dark:text-stone-400 max-w-xs">Your feedback has been received and helps us improve agbalumo.</p>
			<button hx-on:click="this.closest('dialog').close(); this.closest('dialog').remove();" class="mt-4 px-6 py-2 bg-stone-100 dark:bg-stone-800 hover:bg-stone-200 dark:hover:bg-stone-700 text-stone-700 dark:text-stone-300 rounded-full font-bold text-sm transition-colors">
				Close
			</button>
		</div>
	`)
}
