package domain_test

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"testing"
)

// TestListingStore_IsSegregated verifies that ListingStore is now composed of
// ListingReader and ListingWriter interfaces.
func TestListingStore_IsSegregated(t *testing.T) {
	var _ domain.ListingReader = (domain.ListingStore)(nil)
	var _ domain.ListingWriter = (domain.ListingStore)(nil)
}
