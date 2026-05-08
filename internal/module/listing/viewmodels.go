package listing

import (
	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/module"
)

// SaveButtonViewModel is the typed data for the save/unsave toggle fragment.
type SaveButtonViewModel struct {
	ListingID string
	IsSaved   bool
}

// MetricRequest represents the payload for frontend metric ingestion.
type MetricRequest struct {
	Metadata map[string]any `json:"metadata"`
	Event    string         `json:"event"`
	Value    float64        `json:"value"`
}

// HomeViewModel represents the data for the home page.
type HomeViewModel struct {
	SavedIDs         map[string]bool
	Source           string
	Query            string
	City             string
	FilterType       string
	GoogleMapsApiKey string
	Featured         []domain.Listing
	Locations        []domain.Location
	Listings         []domain.Listing
	module.BaseViewData
	Pagination Pagination
	Radius     float64
	TotalCount int
}

// ListingFragmentViewModel represents the data for the listing list fragment.
type ListingFragmentViewModel struct {
	User       interface{}
	SavedIDs   map[string]bool
	Query      string
	City       string
	FilterType string
	Source     string
	Listings   []domain.Listing
	Featured   []domain.Listing
	Pagination Pagination
	Radius     float64
}

// DetailViewModel represents the data for the listing detail modal.
type DetailViewModel struct {
	GoogleMapsApiKey string
	Category         domain.CategoryData
	SavedIDs         map[string]bool
	module.BaseViewData
	Listing  domain.Listing
	CanClaim bool
}

// EditViewModel represents the data for the listing edit modal.
type EditViewModel struct {
	TargetID         string
	Source           string
	GoogleMapsApiKey string
	module.BaseViewData
	Listing domain.Listing
}

// ProfileViewModel represents the data for the user profile page.
type ProfileViewModel struct {
	SavedIDs         map[string]bool
	GoogleMapsApiKey string
	Listings         []domain.Listing
	module.BaseViewData
}

// SavedListingsViewModel represents the data for the saved listings page.
type SavedListingsViewModel struct {
	SavedIDs   map[string]bool
	Source     string
	City       string
	FilterType string
	Query      string
	Listings   []domain.Listing
	Featured   []domain.Listing
	module.BaseViewData
	Pagination Pagination
	Radius     float64
}
