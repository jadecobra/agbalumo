package handler_test

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/handler"
	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/labstack/echo/v4"
	testifyMock "github.com/stretchr/testify/mock"
)

func NewAdminTemplate(t *testing.T) *template.Template {
	wd, _ := os.Getwd()
	projectRoot := filepath.Join(wd, "..", "..")

	funcMap := template.FuncMap{
		"mod":   func(i, j int) int { return i % j },
		"add":   func(i, j int) int { return i + j },
		"sub":   func(i, j int) int { return i - j },
		"split": strings.Split,
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, nil // simplified
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, nil // simplified
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"toJson": func(v interface{}) (template.JS, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return template.JS(b), nil
		},
		"isNew": func(createdAt time.Time) bool {
			if createdAt.IsZero() {
				return false
			}
			return time.Since(createdAt) < 7*24*time.Hour
		},
	}

	tmpl := template.New("base.html").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(
		filepath.Join(projectRoot, "ui", "templates", "base.html"),
		filepath.Join(projectRoot, "ui", "templates", "admin_dashboard.html"),
	)
	if err != nil {
		t.Fatalf("Failed to parse main templates: %v", err)
	}
	_, err = tmpl.ParseGlob(filepath.Join(projectRoot, "ui", "templates", "partials", "*.html"))
	if err != nil {
		t.Fatalf("Failed to parse partial templates: %v", err)
	}
	return tmpl
}

func TestAdminDashboardFooterPosition(t *testing.T) {
	e := echo.New()
	e.Renderer = &RealTemplateRenderer{templates: NewAdminTemplate(t)}

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

	mockRepo.On("GetPendingListings", testifyMock.Anything, 50, 0).Return([]domain.Listing{}, nil)
	mockRepo.On("GetUserCount", testifyMock.Anything).Return(5, nil)
	mockRepo.On("GetFeedbackCounts", testifyMock.Anything).Return(map[domain.FeedbackType]int{domain.FeedbackTypeIssue: 1}, nil)
	mockRepo.On("GetListingGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetUserGrowth", testifyMock.Anything).Return([]domain.DailyMetric{}, nil)
	mockRepo.On("GetAllFeedback", testifyMock.Anything).Return(feedbacks, nil)
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
