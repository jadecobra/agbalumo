package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func createMultipartFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		t.Fatal(err)
	}
	part.Write(content)
	writer.Close()

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Parse
	err = req.ParseMultipartForm(1024)
	if err != nil {
		t.Fatal(err)
	}
	return req.MultipartForm.File["image"][0]
}

func TestLocalImageService_UploadImage_Success(t *testing.T) {
	tempDir := t.TempDir()
	service := &LocalImageService{
		UploadDir:     tempDir,
		MaxUploadSize: 1024 * 1024,
	}

	// Valid JPEG magic bytes
	content := []byte("\xFF\xD8\xFF\xE0\x00\x10JFIF\x00\x01\x01\x01\x00H\x00H\x00\x00")
	fileHeader := createMultipartFileHeader(t, "test.jpg", content)

	path, err := service.UploadImage(context.Background(), fileHeader, "listing-123")
	assert.NoError(t, err)
	assert.Contains(t, path, "/static/uploads/listing-123.jpg")

	// Verify file exists
	_, err = os.Stat(filepath.Join(tempDir, "listing-123.jpg"))
	assert.NoError(t, err)
}

func TestLocalImageService_UploadImage_InvalidType(t *testing.T) {
	tempDir := t.TempDir()
	service := &LocalImageService{UploadDir: tempDir, MaxUploadSize: 1024}

	content := []byte("plain text file disguised as jpg")
	fileHeader := createMultipartFileHeader(t, "test.jpg", content)

	_, err := service.UploadImage(context.Background(), fileHeader, "listing-123")
	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, he.Code)
}
