package maintenance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSurfaceParity(t *testing.T) {
	tests := []struct {
		name           string
		cardContent    string
		modalContent   string
		wantViolations int
	}{
		{
			name:           "All Tokens Match",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white",
			wantViolations: 0,
		},
		{
			name:           "Card Missing Token",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40",
			wantViolations: 1,
		},
		{
			name:           "Modal Missing Token",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark text-white",
			wantViolations: 1,
		},
		{
			name:           "Card has forbidden neobrutalist border",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white border-2 border-stone-900",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white",
			wantViolations: 1,
		},
		{
			name:           "Modal has forbidden neobrutalist shadow",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]",
			wantViolations: 1,
		},
		{
			name:           "Card has forbidden legacy clay styling",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white bg-earth-clay",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white",
			wantViolations: 1,
		},
		{
			name:           "Both have forbidden styles",
			cardContent:    "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white border-2 border-stone-900",
			modalContent:   "bg-gradient-to-b from-earth-espresso to-earth-dark bg-earth-accent/20 text-earth-accent border border-earth-accent/40 text-white shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]",
			wantViolations: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			runSurfaceParityCase(t, tt.cardContent, tt.modalContent, tt.wantViolations)
		})
	}
}

func runSurfaceParityCase(t *testing.T, card, modal string, want int) {
	t.Helper()
	tmpDir := t.TempDir()
	partialsDir := filepath.Join(tmpDir, "ui/templates/partials")
	if err := os.MkdirAll(partialsDir, 0750); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(partialsDir, "listing_card.html"), []byte(card), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(partialsDir, "modal_detail.html"), []byte(modal), 0600); err != nil {
		t.Fatal(err)
	}

	violations, err := CheckSurfaceParity(tmpDir)
	if err != nil {
		t.Fatalf("CheckSurfaceParity failed: %v", err)
	}
	if len(violations) != want {
		t.Errorf("Expected %d violations, got %d", want, len(violations))
	}
}
