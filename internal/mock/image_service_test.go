package mock_test

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/jadecobra/agbalumo/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestMockImageService_UploadImage(t *testing.T) {
	service := new(mock.MockImageService)
	ctx := context.Background()
	file := &multipart.FileHeader{Filename: "test.jpg"}

	service.On("UploadImage", ctx, file, "listing1").Return("url1", nil)

	url, err := service.UploadImage(ctx, file, "listing1")
	assert.NoError(t, err)
	assert.Equal(t, "url1", url)
	service.AssertExpectations(t)
}

func TestMockImageService_DeleteImage(t *testing.T) {
	service := new(mock.MockImageService)
	ctx := context.Background()

	service.On("DeleteImage", ctx, "url1").Return(nil)

	err := service.DeleteImage(ctx, "url1")
	assert.NoError(t, err)
	service.AssertExpectations(t)
}
