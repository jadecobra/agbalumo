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
	"github.com/jadecobra/agbalumo/internal/util"
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

	if err := s.validateUpload(file); err != nil {
		return "", err
	}

	img, err := s.decodeImage(file)
	if err != nil {
		return "", err
	}

	err = util.SafeMkdir(s.UploadDir)
	if err != nil {
		return "", err
	}

	buf, err := s.compressAndScaleImage(img)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(listingID) + ".webp"
	dstPath := filepath.Join(s.UploadDir, filename)
	if err := util.SafeWriteFile(dstPath, buf.Bytes()); err != nil {
		return "", err
	}

	return "/static/uploads/" + filename, nil
}

func (s *LocalImageService) validateUpload(file *multipart.FileHeader) error {
	if file.Size > s.MaxUploadSize {
		return echo.NewHTTPError(http.StatusBadRequest, "File size exceeds 1MB limit")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	buff := make([]byte, 512)
	if _, err = src.Read(buff); err != nil {
		return err
	}

	if !strings.HasPrefix(http.DetectContentType(buff), "image/") {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file type. Only image files are allowed.")
	}
	return nil
}

func (s *LocalImageService) decodeImage(file *multipart.FileHeader) (image.Image, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = src.Close() }()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid or unsupported image file")
	}
	return img, nil
}

func (s *LocalImageService) compressAndScaleImage(img image.Image) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	quality := s.InitialQuality

	for {
		buf.Reset()
		ebuf, err := s.encodeToBuffer(img, quality)
		if err != nil {
			return nil, err
		}
		buf = *ebuf

		if buf.Len() <= int(s.MaxFileSize) || quality <= s.MinQuality {
			break
		}
		quality -= 10
	}

	if buf.Len() > int(s.MaxFileSize) {
		return s.downscaleAndRecompress(img)
	}
	return &buf, nil
}

func (s *LocalImageService) downscaleAndRecompress(img image.Image) (*bytes.Buffer, error) {
	buf, err := s.encodeToBuffer(img, s.MinQuality)
	if err != nil {
		return nil, err
	}

	scale := math.Sqrt(float64(s.MaxFileSize) / float64(buf.Len()))
	if scale < 0.25 {
		scale = 0.25
	}

	b := img.Bounds()
	newWidth, newHeight := int(float64(b.Dx())*scale), int(float64(b.Dy())*scale)
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			newImg.Set(x, y, img.At(x*b.Dx()/newWidth, y*b.Dy()/newHeight))
		}
	}

	return s.encodeToBuffer(newImg, s.MinQuality)
}

func (s *LocalImageService) encodeToBuffer(img image.Image, quality int) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, webp.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return &buf, nil
}

// DeleteImage removes the image from the local filesystem
func (s *LocalImageService) DeleteImage(ctx context.Context, imageURL string) error {
	if imageURL == "" {
		return nil
	}

	filename := filepath.Base(imageURL)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	dstPath := filepath.Join(s.UploadDir, filename)
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

	ebuf, err := s.encodeToBuffer(img, s.InitialQuality)
	if err != nil {
		return nil, err
	}
	return ebuf, nil
}

// ConvertToWebP converts any image to WebP format
func (s *LocalImageService) ConvertToWebP(src io.Reader) (io.Reader, error) {
	return s.CompressImage(src)
}
