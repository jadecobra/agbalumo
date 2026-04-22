package admin_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/admin"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAdminDashboardFooterPosition(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	c, rec := testutil.SetupAdminIntegrationContext(t, http.MethodGet, domain.PathAdmin, nil, domain.TemplateAdminDashboard)
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

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()

	// The footer has class "footer-fruit"
	footerIdx := strings.Index(body, domain.ClassFooter)
	if footerIdx == -1 {
		t.Fatalf("Footer with class '%s' not found in rendered HTML", domain.ClassFooter)
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

func TestMetricCardsHaveHTMXTriggers(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "1", Title: "Business A", Type: domain.Business, IsActive: true})

	c, rec := testutil.SetupAdminIntegrationContext(t, http.MethodGet, domain.PathAdmin, nil, domain.TemplateAdminDashboard)
	h := admin.NewAdminHandler(env.App)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()

	// Total Listings metric → /admin/listings link
	if !strings.Contains(body, `href="`+domain.PathAdminListings+`"`) {
		t.Errorf("Expected Total Listings metric card to link to %s", domain.PathAdminListings)
	}

	// Pending metric → hx-get="/admin/modal/moderation"
	if !strings.Contains(body, `hx-get="/admin/modal/moderation"`) {
		t.Error("Expected Pending metric card to have hx-get=\"/admin/modal/moderation\"")
	}

	// Total Users metric → hx-get="/admin/modal/users"
	if !strings.Contains(body, `hx-get="/admin/modal/users"`) {
		t.Error("Expected Total Users metric card to have hx-get=\"/admin/modal/users\"")
	}

	// HXTarget should pointing to #admin-modal-container
	if !strings.Contains(body, `hx-target="#admin-modal-container"`) {
		t.Error("Expected metric cards to target #admin-modal-container")
	}
}

func TestCategoryModalTrigger(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	c, rec := testutil.SetupAdminIntegrationContext(t, http.MethodGet, domain.PathAdmin, nil, domain.TemplateAdminDashboard)
	h := admin.NewAdminHandler(env.App)

	if err := h.HandleDashboard(c); err != nil {
		t.Fatalf("HandleDashboard failed: %v", err)
	}

	body := rec.Body.String()

	// admin-modal-container div must exist
	if !strings.Contains(body, `id="admin-modal-container"`) {
		t.Error("Expected admin-modal-container div to exist in the rendered dashboard")
	}

	// Categories button in admin tools grid must have hx-get trigger
	if !strings.Contains(body, `hx-get="/admin/modal/category"`) {
		t.Error("Expected Categories admin tool button to have hx-get=\"/admin/modal/category\"")
	}
}

func TestCategoryModalFragment(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	c, rec := testutil.SetupAdminIntegrationContext(t, http.MethodGet, "/admin/modal/category", nil, "components/admin_modal_category.html")
	h := admin.NewAdminHandler(env.App)

	if err := h.HandleModalCategory(c); err != nil {
		t.Fatalf("HandleModalCategory failed: %v", err)
	}

	body := rec.Body.String()

	// categoryModal div must exist in the fragment
	if !strings.Contains(body, `id="`+domain.ModalCategory+`"`) {
		t.Errorf("Expected %s div to exist in the rendered fragment", domain.ModalCategory)
	}

	// Form must post to /admin/categories
	if !strings.Contains(body, `action="`+domain.PathAdminCategories+`"`) {
		t.Errorf("Expected add-category form with action=\"%s\"", domain.PathAdminCategories)
	}

	// Name input must be present
	if !strings.Contains(body, `name="name"`) {
		t.Error("Expected category name input field with name=\"name\"")
	}
}

func TestAdminDashboard_FlashMessages(t *testing.T) {
	t.Parallel()
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()

	c, rec := testutil.SetupAdminIntegrationContext(t, http.MethodGet, domain.PathAdmin, nil, domain.TemplateAdminDashboard)
	h := admin.NewAdminHandler(env.App)

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
	env := testutil.SetupTestModuleEnv(t)
	defer env.Cleanup()
	_ = env.App.DB.Save(context.Background(), domain.Listing{ID: "listing1", Title: "Business A", Type: domain.Business, IsActive: true})

	c, rec := testutil.SetupAdminIntegrationContext(t, http.MethodGet, domain.PathAdminListings, nil, domain.TemplateAdminListings)
	h := admin.NewAdminHandler(env.App)

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
