package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

type dashboardData struct {
	ClaimRequests   []domain.ClaimRequest
	FeedbackCounts  map[domain.FeedbackType]int
	ListingGrowth   []domain.DailyMetric
	UserGrowth      []domain.DailyMetric
	Feedbacks       []domain.Feedback
	Categories      []domain.CategoryData
	Users           []domain.User
	UserCount       int
	ListingCount    int
	AdaDiscoveryAvg float64
}

// HandleDashboard renders the admin dashboard.
func (h *AdminHandler) HandleDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.loadDashboardData(ctx, c)
	if err != nil {
		return ui.RespondError(c, err)
	}

	// Get Flash Messages
	sess := customMiddleware.GetSession(c)
	var flashMsg interface{}
	if sess != nil {
		if flashes := sess.Flashes(domain.FlashMessageKey); len(flashes) > 0 {
			flashMsg = flashes[0]
			_ = sess.Save(c.Request(), c.Response())
		}
	}

	return c.Render(http.StatusOK, "admin_dashboard.html", map[string]interface{}{
		"ClaimRequests":  data.ClaimRequests,
		"UserCount":      data.UserCount,
		"FeedbackCounts": data.FeedbackCounts,
		"ListingGrowth":  data.ListingGrowth,
		"UserGrowth":     data.UserGrowth,
		"Feedbacks":      data.Feedbacks,
		"User":           c.Get(domain.CtxKeyUser),

		"FlashMessage":    flashMsg,
		"ListingCount":    data.ListingCount,
		"Categories":      data.Categories,
		"Users":           data.Users,
		"AdaDiscoveryAvg": data.AdaDiscoveryAvg,
	})
}

func (h *AdminHandler) loadDashboardData(ctx context.Context, c echo.Context) (dashboardData, error) {
	var data dashboardData
	var err error

	data.ClaimRequests, err = h.App.DB.GetPendingClaimRequests(ctx)
	if err != nil {
		return data, err
	}

	data.UserCount, err = h.App.DB.GetUserCount(ctx)
	if err != nil {
		return data, err
	}

	data.FeedbackCounts, err = h.App.DB.GetFeedbackCounts(ctx)
	if err != nil {
		return data, err
	}

	data.ListingGrowth, err = h.App.DB.GetListingGrowth(ctx)
	if err != nil {
		return data, err
	}

	data.UserGrowth, err = h.App.DB.GetUserGrowth(ctx)
	if err != nil {
		return data, err
	}

	data.Feedbacks, err = h.App.DB.GetAllFeedback(ctx)
	if err != nil {
		return data, err
	}

	counts, _ := h.App.DB.GetCounts(ctx)
	for _, count := range counts {
		data.ListingCount += count
	}

	data.Categories, err = h.App.CategorizationSvc.GetCategories(ctx, domain.CategoryFilter{})
	if err != nil {
		c.Logger().Errorf("failed to get categories from service: %v", err)
		data.Categories = []domain.CategoryData{}
	}

	data.Users, err = h.App.DB.GetAllUsers(ctx, 10, 0)
	if err != nil {
		c.Logger().Errorf("failed to get users: %v", err)
		data.Users = []domain.User{}
	}

	// Fetch Ada Metrics (Last 24h)
	since := time.Now().Add(-24 * time.Hour)
	data.AdaDiscoveryAvg, err = h.App.DB.GetAverageValue(ctx, "discovery_success", since)
	if err != nil {
		c.Logger().Errorf("failed to get Ada metrics: %v", err)
	}

	return data, nil
}
