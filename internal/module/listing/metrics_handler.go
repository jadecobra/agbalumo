package listing

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MetricRequest struct {
	Event    string                 `json:"event"`
	Value    float64                `json:"value"`
	Metadata map[string]interface{} `json:"metadata"`
}

// HandleMetricsIngestion receives frontend metrics and logs/saves them.
func (h *ListingHandler) HandleMetricsIngestion(c echo.Context) error {
	var req MetricRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if req.Event == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event name is required"})
	}

	// Capture IP or UserAgent if needed from context
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["ua"] = c.Request().UserAgent()

	h.App.MetricsSvc.LogAndSave(c.Request().Context(), req.Event, req.Value, req.Metadata)

	return c.NoContent(http.StatusOK)
}
