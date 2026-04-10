package listing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleMetricsIngestion(t *testing.T) {
	t.Parallel()
	app, cleanup := testutil.SetupTestAppEnv(t)
	defer cleanup()

	h := NewListingHandler(app)
	e := echo.New()

	t.Run("Valid Metric", func(t *testing.T) {
		reqBody := MetricRequest{
			Event: "test_event",
			Value: 12.3,
			Metadata: map[string]interface{}{
				"foo": "bar",
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/metrics", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Mock Metrics Service
		mockMetrics := app.MetricsSvc.(*testutil.MockMetricsService)
		mockMetrics.On("LogAndSave", mock.Anything, "test_event", 12.3, mock.Anything).Once()

		err := h.HandleMetricsIngestion(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockMetrics.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/metrics", bytes.NewReader([]byte("not json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.HandleMetricsIngestion(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
