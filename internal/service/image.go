package service

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// ImageService defines the interface for handling image uploads
type ImageService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error)
}

// LocalImageService handles saving images to the local filesystem
type LocalImageService struct {
	UploadDir     string
	MaxUploadSize int64
}

// NewLocalImageService creates a new instance with default settings
func NewLocalImageService() *LocalImageService {
	return &LocalImageService{
		UploadDir:     "ui/static/uploads",
		MaxUploadSize: 5 * 1024 * 1024, // 5MB
	}
}

// UploadImage validates and saves the uploaded image
func (s *LocalImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	if file == nil {
		return "", nil // No file to upload
	}

	if file.Size > s.MaxUploadSize {
		return "", echo.NewHTTPError(http.StatusBadRequest, "File size exceeds 5MB limit")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 1. Validate File Content (Magic Bytes)
	buff := make([]byte, 512)
	_, err = src.Read(buff)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/webp" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid file type. Only JPEG, PNG, and WebP are allowed.")
	}

	// Reset file pointer after reading magic bytes
	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// Ensure directory exists
	if err := os.MkdirAll(s.UploadDir, 0755); err != nil {
		return "", err
	}

	// Determine Extension
	ext := ".jpg"
	if fileType == "image/png" {
		ext = ".png"
	} else if fileType == "image/webp" {
		ext = ".webp"
	}

	filename := listingID + ext
	dstPath := filepath.Join(s.UploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return web-accessible path
	return "/static/uploads/" + filename, nil
}
