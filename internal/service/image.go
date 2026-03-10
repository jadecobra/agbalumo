package service

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/webp"
	"github.com/labstack/echo/v4"
	_ "golang.org/x/image/webp"
)

// LocalImageService handles saving images to the local filesystem
type LocalImageService struct {
	UploadDir      string
	MaxUploadSize  int64 // Maximum upload size in bytes
	MaxFileSize    int64 // Maximum final file size in bytes (200KB)
	InitialQuality int   // 1-100, starting compression quality
	MinQuality     int   // 1-100, minimum quality to try before downscaling
}

// NewLocalImageService creates a new instance with the specified upload directory.
func NewLocalImageService(uploadDir string) *LocalImageService {
	if uploadDir == "" {
		uploadDir = "ui/static/uploads"
	}
	return &LocalImageService{
		UploadDir:      uploadDir,
		MaxUploadSize:  1 * 1024 * 1024,
		MaxFileSize:    200 * 1024,
		InitialQuality: 85,
		MinQuality:     20,
	}
}

// UploadImage validates, compresses, and saves the uploaded image
func (s *LocalImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID string) (string, error) {
	if file == nil {
		return "", nil // No file to upload
	}

	if file.Size > s.MaxUploadSize {
		return "", echo.NewHTTPError(http.StatusBadRequest, "File size exceeds 1MB limit")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = src.Close() }()

	// 1. Validate File Content (Magic Bytes) - accept all image types
	buff := make([]byte, 512)
	_, err = src.Read(buff)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(buff)
	if !strings.HasPrefix(fileType, "image/") {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid file type. Only image files are allowed.")
	}

	// Reset file pointer after reading magic bytes
	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// 2. Decode the image
	img, _, err := image.Decode(src)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid or unsupported image file")
	}

	// 3. Ensure directory exists
	err = os.MkdirAll(s.UploadDir, 0755)
	if err != nil {
		return "", err
	}

	// 4. Compress and save as WebP with iterative compression to meet size target
	filename := listingID + ".webp"
	dstPath := filepath.Join(s.UploadDir, filename)

	// First attempt: encode with initial quality
	var buf bytes.Buffer
	quality := s.InitialQuality
	err = webp.Encode(&buf, img, webp.Options{Quality: quality})
	if err != nil {
		return "", err
	}

	// Iteratively reduce quality until we meet the target size
	for buf.Len() > int(s.MaxFileSize) && quality > s.MinQuality {
		buf.Reset()
		quality -= 10 // Reduce quality by 10% each iteration
		if quality < s.MinQuality {
			quality = s.MinQuality
		}
		err = webp.Encode(&buf, img, webp.Options{Quality: quality})
		if err != nil {
			return "", err
		}
	}

	// If still too large after minimum quality, downscale the image
	if buf.Len() > int(s.MaxFileSize) {
		// Calculate scale factor to reduce dimensions
		targetSize := int(s.MaxFileSize)
		currentSize := buf.Len()
		scale := float64(targetSize) / float64(currentSize)
		if scale < 1.0 {
			// Calculate new dimensions (square root because area scales with square)
			scale = math.Sqrt(scale)
			if scale < 0.25 {
				scale = 0.25 // Don't go below 25% of original size
			}

			b := img.Bounds()
			newWidth := int(float64(b.Dx()) * scale)
			newHeight := int(float64(b.Dy()) * scale)

			// Create downscaled image
			newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
			for y := 0; y < newHeight; y++ {
				for x := 0; x < newWidth; x++ {
					srcX := x * b.Dx() / newWidth
					srcY := y * b.Dy() / newHeight
					newImg.Set(x, y, img.At(srcX, srcY))
				}
			}
			img = newImg

			// Re-encode with minimum quality
			buf.Reset()
			err = webp.Encode(&buf, img, webp.Options{Quality: s.MinQuality})
			if err != nil {
				return "", err
			}
		}
	}

	// Write final compressed image to file
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = dst.Close() }()

	_, err = buf.WriteTo(dst)
	if err != nil {
		return "", err
	}

	// Return web-accessible path
	return "/static/uploads/" + filename, nil
}

// DeleteImage removes the image from the local filesystem
func (s *LocalImageService) DeleteImage(ctx context.Context, imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Assuming imageURL is like "/static/uploads/listing-123.jpg"
	// and we need to map it back to the local path.
	filename := filepath.Base(imageURL)
	// Strip query parameters if any (for cache busting)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	dstPath := filepath.Join(s.UploadDir, filename)

	// Check if file exists before trying to delete
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(dstPath)
}

// CompressImage compresses an image buffer and returns compressed WebP bytes
func (s *LocalImageService) CompressImage(src io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = webp.Encode(&buf, img, webp.Options{Quality: s.InitialQuality})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// ConvertToWebP converts any image to WebP format
func (s *LocalImageService) ConvertToWebP(src io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = webp.Encode(&buf, img, webp.Options{Quality: s.InitialQuality})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
