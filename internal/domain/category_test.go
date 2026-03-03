package domain_test

import (
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// Dummy test just to ensure compile and basic defaults
func TestCategoryData(t *testing.T) {
	cd := domain.CategoryData{
		ID:        "Test",
		Name:      "Test",
		Claimable: true,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if cd.ID != "Test" {
		t.Errorf("Expected ID Test, got %s", cd.ID)
	}
}
