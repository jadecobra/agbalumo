package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/labstack/echo/v4"
)

type FeedbackHandler struct {
	Repo domain.ListingRepository
}

func NewFeedbackHandler(repo domain.ListingRepository) *FeedbackHandler {
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
	return c.HTML(http.StatusOK, "<div class='p-4 text-green-600 font-bold'>Thank you for your feedback!</div>")
}
