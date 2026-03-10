package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestAdminDashboardFooterPosition(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplateForPage(t, "admin_dashboard.html")}

	// Create a mock repo with some feedback
	mockRepo := &mock.MockListingRepository{}

	feedbacks := []domain.Feedback{
		{
			ID:        "fb-1",
			Content:   "Great App bug",
			Type:      "Issue",
			CreatedAt: time.Now(),
		},
		{
			ID:        "fb-2",
			Content:   "Another feedback",
			Type:      "Feature",
			CreatedAt: time.Now(),
		},
	}

	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{domain.FeedbackTypeIssue: 1}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return(feedbacks, nil)
	mockRepo.On("GetAllUsers", testifyMock.Anything, 10, 0).Return([]domain.User{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{domain.Business: 10}, nil)

	h := handler.NewAdminHandler(mockRepo, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()

	// The footer has class "footer-fruit"
	footerIdx := strings.Index(body, "footer-fruit")
	if footerIdx == -1 {
		t.Fatal("Footer with class 'footer-fruit' not found in rendered HTML")
	}
	t.Logf("Footer index: %d", footerIdx)

	// All feedback items must appear BEFORE the footer
	for _, fb := range feedbacks {
		fbIdx := strings.Index(body, fb.Content)
		if fbIdx == -1 {
			t.Errorf("Feedback content '%s' not found", fb.Content)
			// Log around where we expect it to be
			start := max(0, footerIdx-2000)
			end := min(len(body), footerIdx+500)
			t.Logf("Body snippet around footer:\n%s", body[start:end])
			continue
		}
		t.Logf("Feedback '%s' index: %d", fb.Content, fbIdx)
		if fbIdx > footerIdx {
			t.Errorf("Regression found: Feedback content '%s' found AFTER footer (fbIdx: %d, footerIdx: %d)", fb.Content, fbIdx, footerIdx)
		}
	}
}

func TestMetricCardsHaveModalTriggers(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplateForPage(t, "admin_dashboard.html")}

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(42, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil)
	mockRepo.On("GetAllUsers", testifyMock.Anything, 10, 0).Return([]domain.User{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{domain.Business: 5}, nil)

	h := handler.NewAdminHandler(mockRepo, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()

	// Total Listings metric → /admin/listings link
	if !strings.Contains(body, `href="/admin/listings"`) {
		t.Error("Expected Total Listings metric card to link to /admin/listings")
	}

	// Pending metric → moderationModal
	if !strings.Contains(body, `data-modal-target="moderationModal"`) {
		t.Error("Expected Pending metric card to have data-modal-target=\"moderationModal\"")
	}

	// Total Users metric → usersModal
	if !strings.Contains(body, `data-modal-target="usersModal"`) {
		t.Error("Expected Total Users metric card to have data-modal-target=\"usersModal\"")
	}

	// Metric cards should be clickable (have open-modal action)
	if !strings.Contains(body, `data-action="open-modal"`) {
		t.Error("Expected metric cards to have data-action=\"open-modal\"")
	}
}

func TestCategoryModalExists(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewRealTemplateForPage(t, "admin_dashboard.html")}

	mockRepo := &mock.MockListingRepository{}
	mockRepo.On("GetPendingClaimRequests", testifyMock.Anything).Return([]domain.ClaimRequest{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return([]domain.Feedback{}, nil)
	mockRepo.On("GetAllUsers", testifyMock.Anything, 10, 0).Return([]domain.User{}, nil)
	mockRepo.On("GetCounts", testifyMock.Anything).Return(map[domain.Category]int{domain.Business: 3}, nil)
	mockRepo.On("GetCategories", testifyMock.Anything, testifyMock.Anything).Return([]domain.CategoryData{
		{ID: "cat-1", Name: "Business", IsSystem: true, Active: true},
		{ID: "cat-2", Name: "Crafts", IsSystem: false, Active: true},
	}, nil)

	h := handler.NewAdminHandler(mockRepo, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()

	// categoryModal div must exist
	if !strings.Contains(body, `id="categoryModal"`) {
		t.Error("Expected categoryModal div to exist in the rendered dashboard")
	}

	// Form must post to /admin/categories
	if !strings.Contains(body, `action="/admin/categories"`) {
		t.Error("Expected add-category form with action=\"/admin/categories\"")
	}

	// Name input must be present
	if !strings.Contains(body, `name="name"`) {
		t.Error("Expected category name input field with name=\"name\"")
	}

	// Categories button in admin tools grid must target categoryModal
	if !strings.Contains(body, `data-target="categoryModal"`) {
		t.Error("Expected Categories admin tool button to target categoryModal")
	}
}
