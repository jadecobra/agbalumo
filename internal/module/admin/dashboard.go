package admin

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
	customMiddleware "github.com/jadecobra/agbalumo/internal/middleware"
	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/labstack/echo/v4"
)

type AdminDashboardViewModel struct {
	ClaimRequests  []domain.ClaimRequest
	FeedbackCounts map[domain.FeedbackType]int
	ListingGrowth  []domain.DailyMetric
	UserGrowth     []domain.DailyMetric
	Feedbacks      []domain.Feedback
	Users          []domain.User
	Categories     []domain.CategoryData
	FlashMessage   interface{}
	module.BaseViewData
	AdaDiscoveryAvg float64
	UserCount       int
	ListingCount    int
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

	vm := AdminDashboardViewModel{
		BaseViewData:    h.PopulateBase(c),
		ClaimRequests:   data.ClaimRequests,
		FeedbackCounts:  data.FeedbackCounts,
		ListingGrowth:   data.ListingGrowth,
		UserGrowth:      data.UserGrowth,
		Feedbacks:       data.Feedbacks,
		Users:           data.Users,
		Categories:      data.Categories,
		UserCount:       data.UserCount,
		ListingCount:    data.ListingCount,
		AdaDiscoveryAvg: data.AdaDiscoveryAvg,
		FlashMessage:    flashMsg,
	}

	return h.RenderTyped(c, "admin_dashboard.html", vm)
}

func (h *AdminHandler) loadDashboardData(ctx context.Context, c echo.Context) (AdminDashboardViewModel, error) {
	var data AdminDashboardViewModel
	var err error

	data.ClaimRequests, err = h.App.DB.GetPendingClaimRequests(ctx)
	if err != nil {
		return data, err
	}

	data.Feedbacks, err = h.App.DB.GetAllFeedback(ctx)
	if err != nil {
		return data, err
	}

	metrics, err := h.App.DB.GetDashboardMetrics(ctx)
	if err != nil {
		return data, err
	}
	data.UserCount = metrics.UserCount
	data.FeedbackCounts = metrics.FeedbackCounts
	data.ListingGrowth = metrics.ListingGrowth
	data.UserGrowth = metrics.UserGrowth
	data.ListingCount = metrics.ListingCount
	data.AdaDiscoveryAvg = metrics.AdaDiscoveryAvg

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

	return data, nil
}
