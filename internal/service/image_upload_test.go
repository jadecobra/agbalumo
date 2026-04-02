package service_test

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestLocalImageService_UploadImage(t *testing.T) {
	svc, tempDir := setupTestImageService(t)

	pngData := createValidPNG()
	assert.NotNil(t, pngData, "failed to create test PNG")

	fileHeader := createMultipartImageRequest(t, "image", "test.png", pngData)
	path, err := svc.UploadImage(context.Background(), fileHeader, "listing-123")

	assert.NoError(t, err)
	assert.Contains(t, path, "/static/uploads/listing-123.webp")

	savedPath := filepath.Join(tempDir, "listing-123.webp")
	_, err = os.Stat(savedPath)
	assert.NoError(t, err)
}

func TestLocalImageService_UploadImage_Validation(t *testing.T) {
	svc, _ := setupTestImageService(t, func(s *service.LocalImageService) {
		s.MaxUploadSize = 100
	})

	path, err := svc.UploadImage(context.Background(), nil, "start")
	assert.NoError(t, err)
	assert.Equal(t, "", path)

	header := &multipart.FileHeader{Filename: "large.png", Size: 1000}
	_, err = svc.UploadImage(context.Background(), header, "large")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File size exceeds")

	tmpfile, _ := os.CreateTemp("", "test.txt")
	_, _ = tmpfile.WriteString("This is a text file")
	_ = tmpfile.Close()
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	fileHeader := createMultipartImageRequest(t, "image", "test.txt", []byte("This is a text file"))

	_, err = svc.UploadImage(context.Background(), fileHeader, "invalid-type")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid file type")
}

func TestLocalImageService_UploadImage_JPEG(t *testing.T) {
	svc, _ := setupTestImageService(t)

	jpegData := createValidJPEG()
	assert.NotNil(t, jpegData, "failed to create test JPEG")

	fileHeader := createMultipartImageRequest(t, "image", "test.jpg", jpegData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "jpeg-listing")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}

func TestLocalImageService_UploadImage_GIF(t *testing.T) {
	svc, _ := setupTestImageService(t)

	gifData := createValidGIF()
	assert.NotNil(t, gifData, "failed to create test GIF")

	fileHeader := createMultipartImageRequest(t, "image", "test.gif", gifData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "gif-listing")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}

func TestLocalImageService_UploadImage_Compression(t *testing.T) {
	svc, _ := setupTestImageService(t)
	svc.InitialQuality = 50

	pngData := createValidPNG()
	fileHeader := createMultipartImageRequest(t, "image", "test.png", pngData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "compress-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")

	savedPath := filepath.Join(svc.UploadDir, "compress-test.webp")
	info, err := os.Stat(savedPath)
	assert.NoError(t, err)
	assert.True(t, info.Size() < 1000, "compressed image should be very small")
}

func TestLocalImageService_UploadImage_LargeImageCompression(t *testing.T) {
	svc, _ := setupTestImageService(t, func(s *service.LocalImageService) {
		s.MaxUploadSize = 5 * 1024 * 1024
		s.InitialQuality = 85
	})

	pngData := createValidPNG()
	fileHeader := createMultipartImageRequest(t, "image", "large.png", pngData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "large-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}

func TestLocalImageService_UploadImage_Downscale(t *testing.T) {
	svc, _ := setupTestImageService(t, func(s *service.LocalImageService) {
		s.MaxFileSize = 500
		s.InitialQuality = 20
		s.MinQuality = 10
	})

	pngData := createValidPNG()
	fileHeader := createMultipartImageRequest(t, "image", "downscale.png", pngData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "downscale-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}

func TestLocalImageService_UploadImage_Errors(t *testing.T) {
	t.Run("mkdir error", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "not-a-dir")
		assert.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		svc, _ := setupTestImageService(t, func(s *service.LocalImageService) {
			s.UploadDir = filepath.Join(tmpFile.Name(), "subdir")
		})

		pngData := createValidPNG()
		header := createMultipartImageRequest(t, "image", "test.png", pngData)

		_, err = svc.UploadImage(context.Background(), header, "err")
		assert.Error(t, err)
	})
}
