package listing_test

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module/listing"
	"github.com/jadecobra/agbalumo/internal/testutil"
)

func init() {
	testutil.ListingServiceConstructor = func(repo domain.ListingRepository) domain.ListingService {
		return listing.NewListingService(repo, repo, repo)
	}
}
