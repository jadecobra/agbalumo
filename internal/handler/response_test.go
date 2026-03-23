package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/labstack/echo/v4"
)

func TestRespondJSONError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errMessage := "this is a test error"
	err := handler.RespondJSONError(c, http.StatusBadRequest, errMessage)
	if err != nil {
		t.Fatalf("expected no error from RespondJSONError, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp handler.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if resp.Error != errMessage {
		t.Errorf("expected error message %q, got %q", errMessage, resp.Error)
	}
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected code in JSON %d, got %d", http.StatusBadRequest, resp.Code)
	}
}
