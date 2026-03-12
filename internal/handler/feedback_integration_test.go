package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedbackHandler_HandleModal(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/feedback/modal", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	h := handler.NewFeedbackHandler(repo)
	
	// Use TestRenderer from listing_helpers_test.go
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}

	err := h.HandleModal(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_Success(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("type", "Issue")
	formData.Set("content", "This is a bug.")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := domain.User{ID: "user1"}
	c.Set("User", user)

	repo := handler.SetupTestRepository(t)

	h := handler.NewFeedbackHandler(repo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "check_circle")

	// Verify feedback in DB
	feedbacks, err := repo.GetAllFeedback(c.Request().Context())
	require.NoError(t, err)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, "user1", feedbacks[0].UserID)
	assert.Equal(t, domain.FeedbackTypeIssue, feedbacks[0].Type)
	assert.Equal(t, "This is a bug.", feedbacks[0].Content)
}

func TestFeedbackHandler_HandleSubmit_NoAuth(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	req := httptest.NewRequest(http.MethodPost, "/feedback", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	repo := handler.SetupTestRepository(t)
	h := handler.NewFeedbackHandler(repo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_InvalidUserType(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	formData := url.Values{}
	formData.Set("content", "test")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", "not a user struct")

	repo := handler.SetupTestRepository(t)
	h := handler.NewFeedbackHandler(repo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_EmptyContent(t *testing.T) {
	e := echo.New()
	e.Renderer = &TestRenderer{templates: NewMainTemplate()}
	formData := url.Values{}
	formData.Set("type", "Issue")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", domain.User{ID: "user1"})

	repo := handler.SetupTestRepository(t)
	h := handler.NewFeedbackHandler(repo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_DefaultType(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("content", "test content")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", domain.User{ID: "user1"})

	repo := handler.SetupTestRepository(t)

	h := handler.NewFeedbackHandler(repo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	feedbacks, err := repo.GetAllFeedback(c.Request().Context())
	require.NoError(t, err)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, domain.FeedbackTypeOther, feedbacks[0].Type)
}

func TestFeedbackHandler_HandleSubmit_TypeFromQueryParam(t *testing.T) {
	e := echo.New()
	formData := url.Values{}
	formData.Set("content", "test content")
	req := httptest.NewRequest(http.MethodPost, "/feedback?type=Feature", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", domain.User{ID: "user1"})

	repo := handler.SetupTestRepository(t)

	h := handler.NewFeedbackHandler(repo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	feedbacks, err := repo.GetAllFeedback(c.Request().Context())
	require.NoError(t, err)
	assert.Equal(t, 1, len(feedbacks))
	assert.Equal(t, domain.FeedbackTypeFeature, feedbacks[0].Type)
}
