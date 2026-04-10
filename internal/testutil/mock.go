package testutil

import (
	"context"
	"io"
	"mime/multipart"
	"strings"

	"github.com/jadecobra/agbalumo/internal/domain"
	testifyMock "github.com/stretchr/testify/mock"
)

type MockGeocodingService struct {
	testifyMock.Mock
}

func (m *MockGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	args := m.Called(ctx, address)
	return args.String(0), args.Error(1)
}

type MockImageService struct {
	testifyMock.Mock
}

func (m *MockImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	args := m.Called(ctx, file, listingID)
	return args.String(0), args.Error(1)
}
func (m *MockImageService) DeleteImage(ctx context.Context, imageURL string) error {
	args := m.Called(ctx, imageURL)
	return args.Error(0)
}

// StubImageService is a non-panicking fake for use in general integration tests.
type StubImageService struct{}

func (m *StubImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	return "http://example.com/image.png", nil
}
func (m *StubImageService) DeleteImage(ctx context.Context, imageURL string) error {
	return nil
}

type MockCSVService struct{}

func (m *MockCSVService) GenerateCSV(ctx context.Context, listings []domain.Listing) (io.Reader, error) {
	return strings.NewReader(""), nil
}
func (m *MockCSVService) ParseAndImport(ctx context.Context, r io.Reader, repo domain.ListingStore) (*domain.BulkUploadResult, error) {
	return &domain.BulkUploadResult{TotalProcessed: 1, SuccessCount: 1}, nil
}

type MockListingService struct{}

func (m *MockListingService) ClaimListing(ctx context.Context, user domain.User, listingID string) (domain.ClaimRequest, error) {
	return domain.ClaimRequest{}, nil
}

type MockCategorizationService struct{}

func (m *MockCategorizationService) GetActiveCategories(ctx context.Context) ([]domain.CategoryData, error) {
	return []domain.CategoryData{}, nil
}

func (m *MockCategorizationService) GetCategories(ctx context.Context, filter domain.CategoryFilter) ([]domain.CategoryData, error) {
	return []domain.CategoryData{}, nil
}

type MockMetricsService struct {
	testifyMock.Mock
}

func (m *MockMetricsService) LogAndSave(ctx context.Context, eventType string, value float64, metadata map[string]interface{}) {
	m.Called(ctx, eventType, value, metadata)
}
