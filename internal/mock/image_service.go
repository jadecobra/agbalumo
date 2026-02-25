package mock

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	args := m.Called(ctx, file, listingID)
	return args.String(0), args.Error(1)
}

func (m *MockImageService) DeleteImage(ctx context.Context, imageURL string) error {
	args := m.Called(ctx, imageURL)
	return args.Error(0)
}
