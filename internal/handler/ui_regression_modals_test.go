package handler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func verifyUsesModalBase(t *testing.T, file string) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", file)
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", file, err)
	}
	if !strings.Contains(string(templateContent), `template "modal_base"`) {
		t.Errorf("%s missing expected modal_base component usage", file)
	}
}

func TestCreateListingModalTheme(t *testing.T) {
	verifyUsesModalBase(t, "modal_create_listing.html")
}

func TestEditListingModalTheme(t *testing.T) {
	verifyUsesModalBase(t, "modal_edit_listing.html")
}

func TestDetailModalTheme(t *testing.T) {
	verifyUsesModalBase(t, "modal_detail.html")
}

func TestProfileModalTheme(t *testing.T) {
	verifyUsesModalBase(t, "modal_profile.html")
}

func TestFeedbackModalTheme(t *testing.T) {
	verifyUsesModalBase(t, "modal_feedback.html")
}

// Now test that modal_base itself conforms to the UI standards
func TestModalBaseTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "components", "modal_base.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Skipf("Skipping modal_base check (maybe it doesn't exist?): %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("modal_base missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `data-modal-action="close"`) {
		t.Error("modal_base expected close button with attribute data-modal-action=\"close\", not found")
	}
}

