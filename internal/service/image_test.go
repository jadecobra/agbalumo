package service_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/stretchr/testify/assert"
)

func createValidPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func createValidJPEG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func TestNewLocalImageService(t *testing.T) {
	svc := service.NewLocalImageService("")

	assert.NotNil(t, svc)
	assert.Equal(t, "ui/static/uploads", svc.UploadDir)
	assert.Equal(t, int64(1*1024*1024), svc.MaxUploadSize)
	assert.Equal(t, int64(200*1024), svc.MaxFileSize)
	assert.Equal(t, 85, svc.InitialQuality)
	assert.Equal(t, 20, svc.MinQuality)
}

func TestLocalImageService_UploadImage(t *testing.T) {
	tempDir := t.TempDir()
	svc := &service.LocalImageService{
		UploadDir:      tempDir,
		MaxUploadSize:  1024 * 1024,
		InitialQuality: 80,
	}

	pngData := createValidPNG()
	assert.NotNil(t, pngData, "failed to create test PNG")

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.png")
	assert.NoError(t, err)
	_, err = part.Write(pngData)
	assert.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	err = req.ParseMultipartForm(1024)
	assert.NoError(t, err)

	fileHeader := req.MultipartForm.File["image"][0]
	path, err := svc.UploadImage(context.Background(), fileHeader, "listing-123")

	assert.NoError(t, err)
	assert.Contains(t, path, "/static/uploads/listing-123.jpg")

	savedPath := filepath.Join(tempDir, "listing-123.jpg")
	_, err = os.Stat(savedPath)
	assert.NoError(t, err)
}

func TestLocalImageService_UploadImage_Validation(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.MaxUploadSize = 100

	path, err := svc.UploadImage(context.Background(), nil, "start")
	assert.NoError(t, err)
	assert.Equal(t, "", path)

	header := &multipart.FileHeader{Filename: "large.png", Size: 1000}
	_, err = svc.UploadImage(context.Background(), header, "large")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "File size exceeds")

	tmpfile, _ := os.CreateTemp("", "test.txt")
	tmpfile.WriteString("This is a text file")
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.txt")
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

func TestLocalImageService_UploadImage_JPEG(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.UploadDir = t.TempDir()

	jpegData := createValidJPEG()
	assert.NotNil(t, jpegData, "failed to create test JPEG")

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write(jpegData)
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(1024)
	fileHeader := req.MultipartForm.File["image"][0]

	path, err := svc.UploadImage(context.Background(), fileHeader, "jpeg-listing")
	assert.NoError(t, err)
	assert.Contains(t, path, ".jpg")
}

func TestLocalImageService_UploadImage_Compression(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.UploadDir = t.TempDir()
	svc.InitialQuality = 50

	pngData := createValidPNG()
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.png")
	part.Write(pngData)
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(1024)
	fileHeader := req.MultipartForm.File["image"][0]

	path, err := svc.UploadImage(context.Background(), fileHeader, "compress-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".jpg")

	savedPath := filepath.Join(svc.UploadDir, "compress-test.jpg")
	info, err := os.Stat(savedPath)
	assert.NoError(t, err)
	assert.True(t, info.Size() < 1000, "compressed image should be very small")
}

func TestLocalImageService_UploadImage_Errors(t *testing.T) {
	t.Run("mkdir error", func(t *testing.T) {
		svc := &service.LocalImageService{
			UploadDir:      "/dev/null/nonexistent/path/that/fails",
			MaxUploadSize:  1024 * 1024,
			InitialQuality: 80,
		}

		pngData := createValidPNG()
		body := &strings.Builder{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("image", "test.png")
		part.Write(pngData)
		writer.Close()

		req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.ParseMultipartForm(1024)
		header := req.MultipartForm.File["image"][0]

		_, err := svc.UploadImage(context.Background(), header, "err")
		assert.Error(t, err)
	})

	t.Run("create error", func(t *testing.T) {
		svc := &service.LocalImageService{
			UploadDir:      "/",
			MaxUploadSize:  1024 * 1024,
			InitialQuality: 80,
		}

		pngData := createValidPNG()
		body := &strings.Builder{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("image", "test.png")
		part.Write(pngData)
		writer.Close()

		req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.ParseMultipartForm(1024)
		header := req.MultipartForm.File["image"][0]

		_, err := svc.UploadImage(context.Background(), header, "err")
		assert.Error(t, err)
	})
}

func TestLocalImageService_CompressImage(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.InitialQuality = 80

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var originalBuf bytes.Buffer
	_ = jpeg.Encode(&originalBuf, img, nil)

	compressed, err := svc.CompressImage(strings.NewReader(originalBuf.String()))
	assert.NoError(t, err)

	compressedData, _ := io.ReadAll(compressed)
	assert.True(t, len(compressedData) > 0, "compressed should have content")
}

func TestLocalImageService_PNGToJPEG(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.InitialQuality = 80

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var pngBuf bytes.Buffer
	_ = png.Encode(&pngBuf, img)

	converted, err := svc.PNGToJPEG(strings.NewReader(pngBuf.String()))
	assert.NoError(t, err)

	convertedData, _ := io.ReadAll(converted)
	_, _, err = image.Decode(bytes.NewReader(convertedData))
	assert.NoError(t, err)
}

func TestLocalImageService_UploadImage_LargeImageCompression(t *testing.T) {
	svc := &service.LocalImageService{
		UploadDir:      t.TempDir(),
		MaxUploadSize:  5 * 1024 * 1024,
		MaxFileSize:    200 * 1024, // 200KB target
		InitialQuality: 85,
		MinQuality:     20,
	}

	// Create a larger image (100x100) that will need compression
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with random-ish pattern to prevent efficient compression
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(x + y), 255})
		}
	}

	var imgBuf bytes.Buffer
	err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 95})
	assert.NoError(t, err)

	// Create multipart form with large image
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "large.jpg")
	part.Write(imgBuf.Bytes())
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(1024 * 1024)
	fileHeader := req.MultipartForm.File["image"][0]

	path, err := svc.UploadImage(context.Background(), fileHeader, "large-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".jpg")

	// Verify the final file is under 200KB
	savedPath := filepath.Join(svc.UploadDir, "large-test.jpg")
	info, err := os.Stat(savedPath)
	assert.NoError(t, err)
	assert.True(t, info.Size() <= 200*1024, "compressed image should be under 200KB but was %d bytes", info.Size())
}

func TestLocalImageService_UploadImage_DecodeError(t *testing.T) {
	svc := &service.LocalImageService{
		UploadDir:      t.TempDir(),
		MaxUploadSize:  1024 * 1024,
		MaxFileSize:    200 * 1024,
		InitialQuality: 85,
	}

	// Create a file with invalid image data
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "invalid.jpg")
	part.Write([]byte("not an image"))
	writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ParseMultipartForm(1024)
	fileHeader := req.MultipartForm.File["image"][0]

	_, err := svc.UploadImage(context.Background(), fileHeader, "decode-error")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid file type")
}

func TestLocalImageService_CompressImage_DecodeError(t *testing.T) {
	svc := service.NewLocalImageService("")

	_, err := svc.CompressImage(strings.NewReader("not an image"))
	assert.Error(t, err)
}

func TestLocalImageService_PNGToJPEG_DecodeError(t *testing.T) {
	svc := service.NewLocalImageService("")

	_, err := svc.PNGToJPEG(strings.NewReader("not a png"))
	assert.Error(t, err)
}
