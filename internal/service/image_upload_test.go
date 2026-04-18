package service

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalImageService_UploadImage(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	svc, _ := setupTestImageService(t, func(s *LocalImageService) {
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

func TestLocalImageService_UploadImage_Formats(t *testing.T) {
	t.Parallel()
	svc, _ := setupTestImageService(t)

	tests := []struct {
		name     string
		filename string
		data     []byte
	}{
		{"JPEG", "test.jpg", createValidJPEG()},
		{"GIF", "test.gif", createValidGIF()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.data, "failed to create test data for "+tt.name)
			fileHeader := createMultipartImageRequest(t, "image", tt.filename, tt.data)

			path, err := svc.UploadImage(context.Background(), fileHeader, tt.name+"-listing")
			assert.NoError(t, err)
			assert.Contains(t, path, ".webp")
		})
	}
}

func TestLocalImageService_UploadImage_Compression(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	svc, _ := setupTestImageService(t, func(s *LocalImageService) {
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
	t.Parallel()
	svc, _ := setupTestImageService(t, func(s *LocalImageService) {
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

func TestLocalImageService_UploadImage_Triggers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		listingName string
		name        string
		maxFileSize int64
		quality     int
		minQual     int
		imgSize     int
	}{
		{
			listingName: "extreme-compress",
			name:        "ExtremeCompression",
			maxFileSize: 1024,
			quality:     90,
			minQual:     10,
			imgSize:     100,
		},
		{
			listingName: "downscale-trigger",
			name:        "Downscale",
			maxFileSize: 100,
			quality:     10,
			minQual:     5,
			imgSize:     200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _ := setupTestImageService(t, func(s *LocalImageService) {
				s.MaxFileSize = tt.maxFileSize
				s.InitialQuality = tt.quality
				s.MinQuality = tt.minQual
			})

			pngData := createCustomPNG(tt.imgSize, tt.imgSize)
			fileHeader := createMultipartImageRequest(t, "image", tt.name+".png", pngData)

			path, err := svc.UploadImage(context.Background(), fileHeader, tt.listingName)
			assert.NoError(t, err)
			assert.Contains(t, path, ".webp")
		})
	}
}

func TestLocalImageService_UploadImage_Errors(t *testing.T) {
	t.Parallel()
	// ... (rest of the error tests)
	t.Run("mkdir error", func(t *testing.T) {
		t.Parallel()
		tmpFile, err := os.CreateTemp("", "not-a-dir")
		assert.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		svc, _ := setupTestImageService(t, func(s *LocalImageService) {
			s.UploadDir = filepath.Join(tmpFile.Name(), "subdir")
		})

		pngData := createValidPNG()
		header := createMultipartImageRequest(t, "image", "test.png", pngData)

		_, err = svc.UploadImage(context.Background(), header, "err")
		assert.Error(t, err)
	})
}
