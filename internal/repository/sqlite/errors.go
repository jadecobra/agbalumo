package sqlite

import "errors"

const (
	errListingNotFound      = "listing not found"
	errClaimRequestNotFound = "claim request not found"
)

var (
	ErrListingNotFound      = errors.New(errListingNotFound)
	ErrClaimRequestNotFound = errors.New(errClaimRequestNotFound)
)
