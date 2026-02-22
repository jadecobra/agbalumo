package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestFeedbackHandler_HandleModal(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/feedback/modal", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewFeedbackHandler(nil)
	e.Renderer = &mock.MockRenderer{}

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

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("SaveFeedback", testifyMock.Anything, testifyMock.MatchedBy(func(f domain.Feedback) bool {
		return f.UserID == "user1" && f.Type == domain.FeedbackTypeIssue && f.Content == "This is a bug."
	})).Return(nil)

	h := NewFeedbackHandler(mockRepo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "check_circle")
	mockRepo.AssertExpectations(t)
}

func TestFeedbackHandler_HandleSubmit_NoAuth(t *testing.T) {
	e := echo.New()
	e.Renderer = &mock.MockRenderer{}
	req := httptest.NewRequest(http.MethodPost, "/feedback", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewFeedbackHandler(nil)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_InvalidUserType(t *testing.T) {
	e := echo.New()
	e.Renderer = &mock.MockRenderer{}
	formData := url.Values{}
	formData.Set("content", "test")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", "not a user struct")

	h := NewFeedbackHandler(nil)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestFeedbackHandler_HandleSubmit_EmptyContent(t *testing.T) {
	e := echo.New()
	e.Renderer = &mock.MockRenderer{}
	formData := url.Values{}
	formData.Set("type", "Issue")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", domain.User{ID: "user1"})

	h := NewFeedbackHandler(nil)

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

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("SaveFeedback", testifyMock.Anything, testifyMock.MatchedBy(func(f domain.Feedback) bool {
		return f.Type == domain.FeedbackTypeOther
	})).Return(nil)

	h := NewFeedbackHandler(mockRepo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
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

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("SaveFeedback", testifyMock.Anything, testifyMock.MatchedBy(func(f domain.Feedback) bool {
		return f.Type == domain.FeedbackTypeFeature
	})).Return(nil)

	h := NewFeedbackHandler(mockRepo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestFeedbackHandler_HandleSubmit_SaveError(t *testing.T) {
	e := echo.New()
	e.Renderer = &mock.MockRenderer{}
	formData := url.Values{}
	formData.Set("content", "test content")
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(formData.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("User", domain.User{ID: "user1"})

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("SaveFeedback", testifyMock.Anything, testifyMock.Anything).Return(assert.AnError)

	h := NewFeedbackHandler(mockRepo)

	err := h.HandleSubmit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockRepo.AssertExpectations(t)
}
