package handler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdminDashboardTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_dashboard.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_dashboard.html: %v", err)
	}

	content := string(templateContent)

	componentsDir := filepath.Join(projectRoot, "ui", "templates", "components")
	files, _ := os.ReadDir(componentsDir)
	for _, f := range files {
		compContent, _ := os.ReadFile(filepath.Join(componentsDir, f.Name()))
		content += string(compContent)
	}

	if !strings.Contains(content, `bg-earth-sand`) {
		t.Error("Admin dashboard metrics card missing semantic sand styling bg-earth-sand")
	}

	if !strings.Contains(content, `bg-earth-dark`) {
		t.Error("Admin dashboard page missing dark theme background bg-earth-dark")
	}

	if !strings.Contains(content, `text-earth-dark`) {
		t.Error("Admin dashboard text missing semantic dark text color text-earth-dark")
	}
}

func TestAdminListingsTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_listings.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_listings.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark flex-1`) {
		t.Error("Admin listings missing expected base dark theme wrapper classes (bg-earth-dark flex-1)")
	}

	if !strings.Contains(content, `divide-white/10`) {
		t.Error("Admin listings table missing translucent divide styling")
	}
}

func TestAdminUsersTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_users.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_users.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark min-h-screen`) {
		t.Error("Admin users missing expected base dark theme wrapper classes")
	}

	if !strings.Contains(content, `divide-white/10`) {
		t.Error("Admin users table missing translucent divide styling")
	}
}

func TestAdminLoginTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_login.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_login.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark font-sans`) {
		t.Error("Admin login body missing expected dark theme classes")
	}

	if !strings.Contains(content, `border-b border-white/20`) {
		t.Error("Admin login input missing border-bottom styling")
	}
}

func TestAdminDashboardModalCloseButtons(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	componentsDir := filepath.Join(projectRoot, "ui", "templates", "components")
	files, err := os.ReadDir(componentsDir)
	if err != nil {
		t.Fatalf("Failed to read components dir: %v", err)
	}

	closeCount := 0
	for _, f := range files {
		content, _ := os.ReadFile(filepath.Join(componentsDir, f.Name()))
		closeCount += strings.Count(string(content), `"Label" "Close"`)
	}

	if closeCount < 4 {
		t.Errorf("admin components should have at least 4 CLOSE bottom buttons (one per modal), found %d", closeCount)
	}
}

func TestAdminListingsUIElements(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "admin_listings.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read admin_listings.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, "pt-32") || !strings.Contains(content, "max-w-6xl") || !strings.Contains(content, "bg-earth-dark") {
		t.Error("Regression: admin_listings.html missing standard admin container classes (pt-32, max-w-6xl, bg-earth-dark)")
	}

	if !strings.Contains(content, "text-[10px]") || !strings.Contains(content, "uppercase") || !strings.Contains(content, "tracking-[0.3em]") {
		t.Error("Regression: admin_listings.html typography missing premium admin styling (text-[10px] uppercase tracking-[0.3em])")
	}

	hasHeaderTemplate := strings.Contains(content, "admin_listing_table_header")
	hasClasses := strings.Contains(content, "bg-white/5") && strings.Contains(content, "text-white/50")
	if !hasHeaderTemplate && !hasClasses {
		t.Error("Regression: admin_listings.html table headers missing premium dark styling (bg-white/5, text-white/50 or admin_listing_table_header component)")
	}
}

func TestAdminPaginationUI(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	// Verify admin_listings.html includes admin_pagination.html
	adminListingsPath := filepath.Join(projectRoot, "ui", "templates", "admin_listings.html")
	adminListingsContent, err := os.ReadFile(adminListingsPath)
	if err != nil {
		t.Fatalf("Failed to read admin_listings.html: %v", err)
	}
	if !strings.Contains(string(adminListingsContent), `template "admin_pagination.html"`) {
		t.Error("Admin listings missing expected admin_pagination.html inclusion")
	}

	// Verify admin_pagination.html does not use HTMX
	adminPaginationPath := filepath.Join(projectRoot, "ui", "templates", "partials", "admin_pagination.html")
	adminPaginationContent, err := os.ReadFile(adminPaginationPath)
	if err != nil {
		t.Fatalf("Failed to read admin_pagination.html: %v", err)
	}
	
	content := string(adminPaginationContent)
	if strings.Contains(content, "hx-get") {
		t.Error("Admin pagination should not use HTMX hx-get, it should use plain full-page reloads")
	}
	if !strings.Contains(content, `href="?page=`) {
		t.Error("Admin pagination missing basic href=\"?page=\" links")
	}
}
