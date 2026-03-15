package handler_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCreateListingModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_create_listing.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_create_listing.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl p-6`) {
		t.Error("Create Listing modal missing expected dark theme wrapper classes")
	}

	hasWrapperClasses := strings.Contains(content, `bg-earth-sand/10 border border-white/20 p-1`)
	hasComponents := strings.Contains(content, `listing_form_title`) || strings.Contains(content, `listing_form_description`)
	if !hasWrapperClasses && !hasComponents {
		t.Error("Create Listing modal inputs missing new sharp border wrapper styling (or component calls)")
	}

	if strings.Contains(content, `multiple`) && strings.Contains(content, `type="file"`) {
		t.Error("Regression: Create Listing modal image input should NOT have 'multiple' attribute (single file upload only)")
	}
}

func TestEditListingModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_edit_listing.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_edit_listing.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Edit Listing modal missing expected dark theme wrapper classes")
	}

	hasWrapperClasses := strings.Contains(content, `bg-earth-sand/10 border border-white/20 p-1`)
	hasComponents := strings.Contains(content, `listing_form_title`) || strings.Contains(content, `listing_form_description`)
	if !hasWrapperClasses && !hasComponents {
		t.Error("Edit Listing modal inputs missing new sharp border wrapper styling (or component calls)")
	}
}

func TestCreateRequestModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_create_request.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_create_request.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl p-6`) {
		t.Error("Create Request modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `bg-earth-sand/10 border border-white/20 p-1`) {
		t.Error("Create Request modal inputs missing sharp border wrapper styling")
	}

	if !strings.Contains(content, `bg-earth-ochre hover:bg-earth-ochre-light`) {
		t.Error("Create Request modal button missing expected earth-ochre styling")
	}
}

func TestDetailModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_detail.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_detail.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Detail modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `<h2 class="text-2xl font-bold font-serif leading-tight">`) {
		t.Error("Detail modal title missing font-serif (Playfair Display) class")
	}
}

func TestProfileModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_profile.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_profile.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Profile modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `<h2 class="text-xl font-bold font-serif text-earth-cream`) {
		t.Error("Profile modal title missing font-serif class")
	}

	if !strings.Contains(content, `data-modal-action="close"`) {
		t.Error("Profile modal close button must use data-modal-action=\"close\" for reliable mobile touch handling")
	}

	if !strings.Contains(content, `id="profile-modal"`) {
		t.Error("Profile modal <dialog> must have id='profile-modal' for modals.js backdrop-click handling")
	}

	if strings.Contains(content, `h-full`) {
		t.Error("Profile modal <dialog> must not use h-full — it should be constrained so the backdrop is visible and clickable")
	}

	if !strings.Contains(content, `bg-white/10`) || !strings.Contains(content, `text-earth-cream`) {
		t.Error("Profile modal item count badge must use bg-white/10 and text-earth-cream for legibility on dark background")
	}

	signOutIdx := strings.Index(content, `/auth/logout`)
	if signOutIdx == -1 {
		t.Error("Profile modal Sign Out link (/auth/logout) not found")
	}
}

func TestFeedbackModalTheme(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_feedback.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_feedback.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `bg-earth-dark/95 backdrop-blur-xl border border-white/10`) {
		t.Error("Feedback modal missing expected dark theme wrapper classes")
	}

	if !strings.Contains(content, `border border-white/20 bg-white/5`) {
		t.Error("Feedback modal textarea missing translucent styling")
	}
}

func TestModalCloseButtons(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	partialsDir := filepath.Join(projectRoot, "ui", "templates", "partials")

	type modalCheck struct {
		file     string
		wantText string
		wantAttr string
	}

	checks := []modalCheck{
		{file: "modal_create_listing.html", wantText: "CLOSE", wantAttr: `data-modal-action="close"`},
		{file: "modal_create_request.html", wantText: "CLOSE", wantAttr: `data-modal-action="close"`},
		{file: "modal_profile.html", wantText: "CLOSE", wantAttr: `data-modal-action="close"`},
		{file: "modal_feedback.html", wantText: "CLOSE", wantAttr: `data-modal-action="close"`},
		{file: "modal_detail.html", wantText: "CLOSE", wantAttr: `data-modal-action="close"`},
		{file: "modal_edit_listing.html", wantText: "CLOSE", wantAttr: `data-modal-action="close"`},
	}

	for _, check := range checks {
		t.Run(check.file, func(t *testing.T) {
			path := filepath.Join(partialsDir, check.file)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", check.file, err)
			}
			body := string(content)

			if !strings.Contains(body, check.wantText) {
				t.Errorf("%s: expected CLOSE button with text %q, not found", check.file, check.wantText)
			}
			if !strings.Contains(body, check.wantAttr) {
				t.Errorf("%s: expected close button with attribute %q, not found", check.file, check.wantAttr)
			}
		})
	}
}

func TestModalCloseButtonStyle(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	partialsDir := filepath.Join(projectRoot, "ui", "templates", "partials")

	modalFiles := []string{
		"modal_create_listing.html",
		"modal_create_request.html",
	}

	for _, file := range modalFiles {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(partialsDir, file)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", file, err)
			}
			body := string(content)

			if !strings.Contains(body, "bg-earth-ochre") {
				t.Errorf("%s: CLOSE button missing bg-earth-ochre class (should match ASK/POST button style)", file)
			}
		})
	}
}

func TestModalNoOrphanIconOnlyCloseButton(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")
	partialsDir := filepath.Join(projectRoot, "ui", "templates", "partials")

	iconOnlyPattern := regexp.MustCompile(`(?s)<button[^>]*data-modal-action="close"[^>]*>\s*<span class="material-symbols-outlined[^"]*">close</span>\s*</button>`)

	modalFiles := []string{
		"modal_create_listing.html",
		"modal_create_request.html",
		"modal_profile.html",
		"modal_feedback.html",
		"modal_detail.html",
		"modal_edit_listing.html",
	}

	for _, file := range modalFiles {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(partialsDir, file)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s: %v", file, err)
			}
			if iconOnlyPattern.Match(content) {
				t.Errorf("%s: found icon-only close button with data-modal-action=\"close\" — replace with CLOSE text label for mobile accessibility", file)
			}
		})
	}
}

func TestDetailModalAddressLink(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	templatePath := filepath.Join(projectRoot, "ui", "templates", "partials", "modal_detail.html")
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read modal_detail.html: %v", err)
	}

	content := string(templateContent)

	if !strings.Contains(content, `<a href="https://www.google.com/maps/search/?api=1&query={{ urlquery .Listing.Address }},{{ urlquery .Listing.City }}"`) {
		t.Error("Detail modal address is not wrapped in a Google Maps link")
	}
}
