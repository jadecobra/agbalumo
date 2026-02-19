package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

type FeedbackHandler struct {
	Repo domain.FeedbackStore
}

func NewFeedbackHandler(repo domain.FeedbackStore) *FeedbackHandler {
	return &FeedbackHandler{Repo: repo}
}

// HandleModal renders the feedback modal form
func (h *FeedbackHandler) HandleModal(c echo.Context) error {
	// user := c.Get("User")
	// If user not logged in, maybe redirect or show login?
	// But modal is usually triggered from UI where user is known or checks are done.
	return c.Render(http.StatusOK, "modal_feedback.html", nil)
}

// HandleSubmit processes the feedback form submission
func (h *FeedbackHandler) HandleSubmit(c echo.Context) error {
	user, ok := c.Get("User").(domain.User)
	if !ok || user.ID == "" {
		return c.String(http.StatusUnauthorized, "Login required")
	}

	contentType := c.QueryParam("type")
	if contentType == "" {
		contentType = c.FormValue("type")
	}

	content := c.FormValue("content")

	// Validate
	if content == "" {
		return c.String(http.StatusBadRequest, "Content is required")
	}
	if contentType == "" {
		contentType = string(domain.FeedbackTypeOther)
	}

	feedback := domain.Feedback{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Type:      domain.FeedbackType(contentType),
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := h.Repo.SaveFeedback(c.Request().Context(), feedback); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save feedback")
	}

	// Return success message or close modal
	// We can return a small HTML snippet to replace the form or Close the modal
	// For now, let's just return a success message or a specialized template
	// Return success message
	return c.HTML(http.StatusOK, `
		<div class="flex flex-col items-center justify-center p-8 space-y-4 text-center animate-in fade-in zoom-in-95 duration-300">
			<div class="h-16 w-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mb-2">
				<span class="material-symbols-outlined text-4xl text-green-600 dark:text-green-400">check_circle</span>
			</div>
			<h3 class="text-xl font-bold text-stone-800 dark:text-white">Thank You!</h3>
			<p class="text-stone-500 dark:text-stone-400 max-w-xs">Your feedback has been received and helps us improve Agbalumo.</p>
			<button onclick="this.closest('dialog').close(); this.closest('dialog').remove();" class="mt-4 px-6 py-2 bg-stone-100 dark:bg-stone-800 hover:bg-stone-200 dark:hover:bg-stone-700 text-stone-700 dark:text-stone-300 rounded-full font-bold text-sm transition-colors">
				Close
			</button>
		</div>
	`)
}
