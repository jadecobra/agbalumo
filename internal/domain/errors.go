package domain

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found in the repository.
	ErrUserNotFound = errors.New("user not found")
	// ErrListingNotFound is returned when a listing is not found.
	ErrListingNotFound = errors.New("listing not found")
	// ErrCategoryNotFound is returned when a category is not found.
	ErrCategoryNotFound = errors.New("category not found")
	// ErrCategoryInactive is returned when a category is inactive.
	ErrCategoryInactive = errors.New("category is inactive")
	// ErrClaimNotFound is returned when a claim record is not found.
	ErrClaimNotFound = errors.New("claim record not found")
	// ErrListingOwned is returned when attempting to claim an already owned listing.
	ErrListingOwned = errors.New("listing is already owned")
	// ErrLoginRequired is returned when an operation requires an authenticated user.
	ErrLoginRequired = errors.New("Login required")
	// ErrUserIDRequired is returned when a user ID is missing.
	ErrUserIDRequired = errors.New("user ID is required")
	// ErrInvalidCategoryType is returned when a category type is invalid.
	ErrInvalidCategoryType = errors.New("invalid category type")
	// ErrListingNotClaimable is returned when a listing's category is not claimable.
	ErrListingNotClaimable = errors.New("listing type cannot be claimed")
	// ErrPendingClaimExists is returned when a user already has a pending claim.
	ErrPendingClaimExists = errors.New("you already have a pending claim for this listing")
	// ErrFailedToSaveClaim is returned when a claim record cannot be persisted.
	ErrFailedToSaveClaim = errors.New("failed to save claim request")
)
