package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminDashboardFooterPosition(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_dashboard.html")}

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	// Seed some feedback
	_ = env.App.DB.SaveFeedback(context.Background(), domain.Feedback{
		ID:        "fb-1",
		Content:   "Great App bug",
		Type:      "Issue",
		CreatedAt: time.Now(),
	})
	_ = env.App.DB.SaveFeedback(context.Background(), domain.Feedback{
		ID:        "fb-2",
		Content:   "Another feedback",
		Type:      "Feature",
		CreatedAt: time.Now(),
	})

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

	// All feedback items must appear BEFORE the footer
	feedbacks := []string{"Great App bug", "Another feedback"}
	for _, content := range feedbacks {
		fbIdx := strings.Index(body, content)
		if fbIdx == -1 {
			t.Errorf("Feedback content '%s' not found", content)
			continue
		}
		if fbIdx > footerIdx {
			t.Errorf("Regression found: Feedback content '%s' found AFTER footer (fbIdx: %d, footerIdx: %d)", content, fbIdx, footerIdx)
		}
	}
}

func TestMetricCardsHaveModalTriggers(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_dashboard.html")}

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Business A", Type: domain.Business, IsActive: true})
	h := admin.NewAdminHandler(env.App)

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
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_dashboard.html")}

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

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

func TestAdminDashboard_FlashMessages(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_dashboard.html")}

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin", nil)
	c.Echo().Renderer = e.Renderer

	// Set up a session with a flash message
	sess, _ := testutil.GetAuthSession(c)
	sess.AddFlash("Success message", domain.FlashMessageKey)
	_ = sess.Save(c.Request(), rec)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()
	assert.Contains(t, body, "Success message")
}

func TestAdminListings_ModalTrigger(t *testing.T) {
	t.Parallel()
	e := echo.New()
	e.Renderer = &testutil.RealTemplateRenderer{Templates: testutil.NewRealTemplateForPage(t, "admin_listings.html")}

	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "listing1", Title: "Business A", Type: domain.Business, IsActive: true})

	h := admin.NewAdminHandler(env.App)

	c, rec := testutil.SetupAdminContext(http.MethodGet, "/admin/listings", nil)
	c.Echo().Renderer = e.Renderer

	if err := h.HandleAllListings(c); err != nil {
		t.Fatalf("HandleAllListings failed: %v", err)
	}

	body := rec.Body.String()

	if strings.Contains(body, `target="_blank"`) && strings.Contains(body, `href="/listings/listing1"`) {
		t.Error("Listing title link should not use target='_blank' because it's a raw modal fragment")
	}

	if !strings.Contains(body, `hx-get="/listings/listing1"`) {
		t.Error("Expected listing title link to trigger HTMX modal using hx-get")
	}
}
