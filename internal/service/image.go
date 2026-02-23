package service

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
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
	UploadDir          string
	MaxUploadSize      int64
	CompressionQuality int // 1-100, higher = better quality, larger file
}

// NewLocalImageService creates a new instance with default settings
func NewLocalImageService() *LocalImageService {
	return &LocalImageService{
		UploadDir:          "ui/static/uploads",
		MaxUploadSize:      5 * 1024 * 1024, // 5MB
		CompressionQuality: 80,              // Good balance of quality and file size
	}
}

// UploadImage validates, compresses, and saves the uploaded image
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

	// 2. Decode the image
	img, _, err := image.Decode(src)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}

	// 3. Ensure directory exists
	if err := os.MkdirAll(s.UploadDir, 0755); err != nil {
		return "", err
	}

	// 4. Compress and save as JPEG (best compression)
	filename := listingID + ".jpg"
	dstPath := filepath.Join(s.UploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Encode with compression quality
	quality := s.CompressionQuality
	if quality <= 0 || quality > 100 {
		quality = 80
	}

	err = jpeg.Encode(dst, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return "", err
	}

	// Return web-accessible path
	return "/static/uploads/" + filename, nil
}

// CompressImage compresses an image buffer and returns compressed bytes
func (s *LocalImageService) CompressImage(src io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: s.CompressionQuality})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// PNGToJPEG converts a PNG image to JPEG
func (s *LocalImageService) PNGToJPEG(src io.Reader) (io.Reader, error) {
	img, err := png.Decode(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: s.CompressionQuality})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
