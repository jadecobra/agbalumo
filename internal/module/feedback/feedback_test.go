package feedback

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedbackHandler_HandleModal(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/feedback/modal", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := NewFeedbackHandler(app)

	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}

	err := h.HandleModal(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_Success(t *testing.T) {
	t.Parallel()
	e := echo.New()
	formData := url.Values{}
	formData.Set("type", "Issue")
	formData.Set("content", "This is a bug.")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user in context (pointer)
	mockUser := &domain.User{ID: "user1"}
	c.Set("User", mockUser)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := NewFeedbackHandler(app)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "check_circle")

	// Verify feedback in DB
	feedbacks, err := app.DB.GetAllFeedback(c.Request().Context())
	require.NoError(t, err)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, "user1", feedbacks[0].UserID)
	assert.Equal(t, domain.FeedbackTypeIssue, feedbacks[0].Type)
	assert.Equal(t, "This is a bug.", feedbacks[0].Content)
}

func TestFeedbackHandler_HandleSubmit_NoAuth(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}
	req := httptest.NewRequest(http.MethodPost, "/feedback", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := NewFeedbackHandler(app)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_EmptyContent(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.TestRenderer{Templates: testutil.NewMainTemplate()}
	formData := url.Values{}
	formData.Set("type", "Issue")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", domain.User{ID: "user1"})

	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()
	h := NewFeedbackHandler(app)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
