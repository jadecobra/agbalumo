package service

import (
	"context"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func (s *CSVService) isDuplicate(ctx context.Context, repo domain.ListingStore, listing *domain.Listing) (bool, error) {
	existingListings, err := repo.FindByTitle(ctx, listing.Title)
	if err != nil {
		return false, err
	}

	for _, ex := range existingListings {
		if s.matchFieldsCount(ex, listing) > 2 {
			return true, nil
		}
	}
	return false, nil
}

func (s *CSVService) matchFieldsCount(ex domain.Listing, listing *domain.Listing) int {
	matches := 0
	check := func(s1, s2 string) {
		if s1 == s2 && s1 != "" {
			matches++
		}
	}

	if ex.Type == listing.Type {
		matches++
	}
	check(ex.Description, listing.Description)
	check(ex.OwnerOrigin, listing.OwnerOrigin)
	check(ex.ContactEmail, listing.ContactEmail)
	check(ex.ContactPhone, listing.ContactPhone)
	check(ex.ContactWhatsApp, listing.ContactWhatsApp)
	check(ex.Address, listing.Address)
	check(ex.City, listing.City)
	return matches
}
