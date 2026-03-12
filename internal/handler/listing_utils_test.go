package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"  ", ""},
		{"example.com", "https://example.com"},
		{"http://test.com", "http://test.com"},
		{"https://secure.com", "https://secure.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeURL(tt.input))
		})
	}
}

func TestIsImageError(t *testing.T) {
	assert.False(t, IsImageError(nil))
	assert.True(t, IsImageError(errors.New("File size exceeds")))
	assert.True(t, IsImageError(errors.New("Invalid file type")))
	assert.True(t, IsImageError(errors.New("Invalid or unsupported image")))
	assert.False(t, IsImageError(errors.New("Other error")))
}

func TestRenderImageErrorToast(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &ListingHandler{}
	
	err := h.renderImageErrorToast(c, echo.NewHTTPError(http.StatusBadRequest, "some error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "some error")
}
