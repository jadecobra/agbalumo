package domain_test

import (
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
)

// TestListingStore_IsSegregated verifies that ListingStore is now composed of
// ListingReader and ListingWriter interfaces.
func TestListingStore_IsSegregated(t *testing.T) {
	t.Parallel()
	var _ domain.ListingReader = (domain.ListingStore)(nil)
	var _ domain.ListingWriter = (domain.ListingStore)(nil)
}
