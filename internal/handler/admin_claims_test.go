package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestHandleApproveClaim(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "claim1", domain.ClaimStatusApproved).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/admin/claims/claim1/approve", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("claim1")

	if assert.NoError(t, h.HandleApproveClaim(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleApproveClaim_Error(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "bad", domain.ClaimStatusApproved).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/admin/claims/bad/approve", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bad")

	_ = h.HandleApproveClaim(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleRejectClaim(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "claim1", domain.ClaimStatusRejected).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/admin/claims/claim1/reject", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("claim1")

	if assert.NoError(t, h.HandleRejectClaim(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestHandleRejectClaim_Error(t *testing.T) {
	e := echo.New()
	mockRepo := &mock.MockListingRepository{}
	h := handler.NewAdminHandler(mockRepo, nil, &config.Config{})

	mockRepo.On("UpdateClaimRequestStatus", testifyMock.Anything, "bad", domain.ClaimStatusRejected).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/admin/claims/bad/reject", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("bad")

	_ = h.HandleRejectClaim(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
