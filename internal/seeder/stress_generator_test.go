package seeder_test

import (
	"testing"
	"github.com/jadecobra/agbalumo/internal/seeder"
)

func TestGenerateStressListings(t *testing.T) {
	count := 100
	listings := seeder.GenerateStressListings(count)

	if len(listings) != count {
		t.Fatalf("expected %d listings, got %d", count, len(listings))
	}

	for i, l := range listings {
		if err := l.Validate(); err != nil {
			t.Errorf("listing %d failed validation: %v\n%+v", i, err, l)
		}
	}
}
