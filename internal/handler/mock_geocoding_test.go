package handler_test

import "context"

type MockGeocodingService struct {
	GetCityFunc func(ctx context.Context, address string) (string, error)
}

func (m *MockGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	if m.GetCityFunc != nil {
		return m.GetCityFunc(ctx, address)
	}
	return "", nil
}
