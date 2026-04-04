package testutil

import "context"

type MockGeocodingService struct{}

func (m *MockGeocodingService) GetCoordinates(ctx context.Context, address string) (float64, float64, error) {
	return 0.0, 0.0, nil
}
func (m *MockGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	return "Lagos", nil
}
