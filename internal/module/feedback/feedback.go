package feedback

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/infra/env"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/module/user"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

type FeedbackHandler struct {
	module.BaseHandler
}

func NewFeedbackHandler(app *env.AppEnv) *FeedbackHandler {
	return &FeedbackHandler{
		BaseHandler: module.BaseHandler{App: app},
	}
}

type FeedbackViewModel struct {
	Path string
	module.BaseViewData
}

// RegisterRoutes registers the feedback routes
func (h *FeedbackHandler) RegisterRoutes(e *echo.Echo, authMw domain.AuthMiddleware) {
	// Public group: feedback is intentionally low-friction to support acquisition
	// (e.g. anonymous triers from seeded links in Option 1 closed-network blasts).
	// Global OptionalAuth middleware still populates User context when present.
	feedbackGroup := e.Group("/feedback")
	feedbackGroup.GET("/modal", h.HandleFeedbackForm)
	feedbackGroup.POST("", h.HandleSubmit)
}

// HandleFeedbackForm renders the feedback modal form
func (h *FeedbackHandler) HandleFeedbackForm(c echo.Context) error {
	vm := FeedbackViewModel{
		BaseViewData: h.PopulateBase(c),
		Path:         c.Request().URL.Path,
	}
	return h.RenderTyped(c, "modal_feedback.html", vm)
}

// HandleSubmit processes the feedback form submission.
// Auth is optional: anonymous submissions (e.g. from seeded acquisition links) are allowed
// with UserID left empty. Global OptionalAuth populates context when a user is signed in.
func (h *FeedbackHandler) HandleSubmit(c echo.Context) error {
	userID := ""
	if u, ok := user.GetUser(c); ok && u != nil {
		userID = u.ID
	}

	contentType := c.QueryParam(domain.FieldType)
	if contentType == "" {
		contentType = c.FormValue(domain.FieldType)
	}

	content := c.FormValue(domain.FieldContent)

	// Validate
	if content == "" {
		return ui.RespondErrorMsg(c, http.StatusBadRequest, "Content is required")
	}
	if contentType == "" {
		contentType = string(domain.FeedbackTypeOther)
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()
	hasher := sha256.New()
	hasher.Write([]byte(ip + "|" + ua))
	fingerprint := hex.EncodeToString(hasher.Sum(nil))

	fb := domain.Feedback{
		ID:          uuid.New().String(),
		UserID:      userID,
		Type:        domain.FeedbackType(contentType),
		Content:     content,
		CreatedAt:   time.Now(),
		Fingerprint: fingerprint,
		Resolved:    false,
	}

	if err := h.App.DB.SaveFeedback(c.Request().Context(), fb); err != nil {
		return ui.RespondErrorMsg(c, http.StatusInternalServerError, "Failed to save feedback")
	}

	// Return success message or close modal
	return c.HTML(http.StatusOK, `
		<div class="flex flex-col items-center justify-center p-8 space-y-4 text-center animate-in fade-in zoom-in-95 duration-300">
			<div class="h-16 w-16 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mb-2">
				<span class="material-symbols-outlined text-4xl text-green-600 dark:text-green-400">check_circle</span>
			</div>
			<h3 class="text-xl font-bold text-stone-800 dark:text-white">Thank You!</h3>
			<p class="text-stone-500 dark:text-stone-400 max-w-xs">Your feedback has been received and helps us improve agbalumo.</p>
			<button data-modal-action="close" class="mt-4 px-6 py-2 bg-stone-100 dark:bg-stone-800 hover:bg-stone-200 dark:hover:bg-stone-700 text-stone-700 dark:text-stone-300 rounded-full font-bold text-sm transition-colors">
				Close
			</button>
		</div>
	`)
}
