package service_test

import (
	"context"
	"mime/multipart" // Added by instruction
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestNewLocalImageService(t *testing.T) {
	svc := service.NewLocalImageService()

	assert.NotNil(t, svc)
	assert.Equal(t, "ui/static/uploads", svc.UploadDir)
	assert.Equal(t, int64(5*1024*1024), svc.MaxUploadSize)
}

func TestLocalImageService_UploadImage(t *testing.T) {
	// Setup temporary directory
	tempDir := t.TempDir()
	svc := &service.LocalImageService{
		UploadDir:     tempDir,
		MaxUploadSize: 1024 * 1024,
	}

	// Create a dummy image file
	tmpfile, err := os.CreateTemp("", "test-*.png")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write PNG magic bytes
	_, err = tmpfile.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	assert.NoError(t, err)
	tmpfile.Close()

	// Open the file to create multipart.File
	file, err := os.Open(tmpfile.Name())
	assert.NoError(t, err)
	defer file.Close()

	// Create multipart.FileHeader
	// We need to set the internal 'content' of FileHeader which is not exported/accessible directly
	// without using multipart.Writer.
	// Standard way is to create a form file.

	// Real multipart setup
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.png")
	assert.NoError(t, err)

	// Write magic bytes
	_, err = part.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	assert.NoError(t, err)
	writer.Close()

	// Now read it back?
	// The Service expects *multipart.FileHeader.
	// We can get it from an http.Request.

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// ParseMultipartForm
	err = req.ParseMultipartForm(1024)
	assert.NoError(t, err)

	fileHeader := req.MultipartForm.File["image"][0]

	// Execute
	path, err := svc.UploadImage(context.Background(), fileHeader, "listing-123")

	// Verify
	assert.NoError(t, err)
	assert.Contains(t, path, "/static/uploads/listing-123.png")

	// Verify file existence
	savedPath := filepath.Join(tempDir, "listing-123.png")
	_, err = os.Stat(savedPath)
	assert.NoError(t, err)
}

func TestLocalImageService_UploadImage_Validation(t *testing.T) {
	svc := service.NewLocalImageService()
	svc.MaxUploadSize = 100 // 100 bytes for testing

	// 1. Nil File
	path, err := svc.UploadImage(context.Background(), nil, "start")
	assert.NoError(t, err)
	assert.Equal(t, "", path)

	// 2. File Too Large
	header := &multipart.FileHeader{
		Filename: "large.png",
		Size:     1000,
	}
	_, err = svc.UploadImage(context.Background(), header, "large")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File size exceeds")

	// 3. Invalid Type (Text file)
	// Setup text file
	tmpfile, err := os.CreateTemp("", "test.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("This is a text file")
	tmpfile.Close()

	file, err := os.Open(tmpfile.Name())
	assert.NoError(t, err)
	defer file.Close()

	// Create multipart form
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.txt")
	assert.NoError(t, err)
	part.Write([]byte("This is a text file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(1024)
	fileHeader := req.MultipartForm.File["image"][0]

	_, err = svc.UploadImage(context.Background(), fileHeader, "invalid-type")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid file type")
}

func TestLocalImageService_UploadImage_WebP(t *testing.T) {
	svc := service.NewLocalImageService()
	svc.UploadDir = t.TempDir()

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.webp")
	assert.NoError(t, err)

	// Write WebP magic bytes: RIFF....WEBP
	// RIFF + size + WEBP + VP8
	// 52 49 46 46 (RIFF)
	// .. .. .. .. (Size)
	// 57 45 42 50 (WEBP)
	// 56 50 38 20 (VP8 )
	magic := []byte{
		0x52, 0x49, 0x46, 0x46,
		0x00, 0x00, 0x00, 0x00,
		0x57, 0x45, 0x42, 0x50,
		0x56, 0x50, 0x38, 0x20,
	}
	part.Write(magic)
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(1024)
	fileHeader := req.MultipartForm.File["image"][0]

	path, err := svc.UploadImage(context.Background(), fileHeader, "webp-listing")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}
