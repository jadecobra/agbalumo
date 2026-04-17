package listing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleMetricsIngestion(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	h := NewListingHandler(env.App)

	t.Run("Valid Metric", func(t *testing.T) {
		reqBody := MetricRequest{
			Event: "test_event",
			Value: 12.3,
			Metadata: map[string]interface{}{
				"foo": "bar",
			},
		}
		body, _ := json.Marshal(reqBody)
		c, rec := testutil.SetupModuleContext(http.MethodPost, "/api/metrics", bytes.NewReader(body))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		// Mock Metrics Service
		mockMetrics := env.App.MetricsSvc.(*testutil.MockMetricsService)
		mockMetrics.On("LogAndSave", mock.Anything, "test_event", 12.3, mock.Anything).Once()

		err := h.HandleMetricsIngestion(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockMetrics.AssertExpectations(t)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		c, rec := testutil.SetupModuleContext(http.MethodPost, "/api/metrics", bytes.NewReader([]byte("not json")))
		c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		err := h.HandleMetricsIngestion(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
