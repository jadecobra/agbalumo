package domain

import "context"

// GeocodingService defines the contract for converting addresses to location metadata.
type GeocodingService interface {
	GetCity(ctx context.Context, address string) (string, error)
}
