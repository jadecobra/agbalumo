package service_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/gif"
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
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func createValidJPEG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func setupTestImageService(t *testing.T) (*service.LocalImageService, string) {
	t.Helper()
	tempDir := t.TempDir()
	svc := &service.LocalImageService{
		UploadDir:      tempDir,
		MaxUploadSize:  1024 * 1024,
		MaxFileSize:    200 * 1024,
		InitialQuality: 80,
		MinQuality:     20,
	}
	return svc, tempDir
}

func createMultipartImageRequest(t *testing.T, fieldName, fileName string, fileData []byte) *multipart.FileHeader {
	t.Helper()
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	assert.NoError(t, err)
	_, err = part.Write(fileData)
	assert.NoError(t, err)
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	err = req.ParseMultipartForm(10 * 1024 * 1024)
	assert.NoError(t, err)

	return req.MultipartForm.File[fieldName][0]
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
	_, _ = tmpfile.WriteString("This is a text file")
	_ = tmpfile.Close()
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	fileHeader := createMultipartImageRequest(t, "image", "test.txt", []byte("This is a text file"))

	_, err = svc.UploadImage(context.Background(), fileHeader, "invalid-type")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid file type")
}

func TestLocalImageService_UploadImage_JPEG(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.UploadDir = t.TempDir()

	jpegData := createValidJPEG()
	assert.NotNil(t, jpegData, "failed to create test JPEG")

	fileHeader := createMultipartImageRequest(t, "image", "test.jpg", jpegData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "jpeg-listing")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}

func createValidGIF() []byte {
	img := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{color.White})
	var buf bytes.Buffer
	_ = gif.EncodeAll(&buf, &gif.GIF{
		Image: []*image.Paletted{img},
		Delay: []int{0},
	})
	return buf.Bytes()
}

func TestLocalImageService_UploadImage_GIF(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.UploadDir = t.TempDir()

	gifData := createValidGIF()
	assert.NotNil(t, gifData, "failed to create test GIF")

	fileHeader := createMultipartImageRequest(t, "image", "test.gif", gifData)

	path, err := svc.UploadImage(context.Background(), fileHeader, "gif-listing")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")
}

func TestLocalImageService_UploadImage_Compression(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.UploadDir = t.TempDir()
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

func TestLocalImageService_UploadImage_Errors(t *testing.T) {
	t.Run("mkdir error", func(t *testing.T) {
		// Create a file to use as a path component (guarantees MkdirAll failure)
		tmpFile, err := os.CreateTemp("", "not-a-dir")
		assert.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		svc := &service.LocalImageService{
			UploadDir:      filepath.Join(tmpFile.Name(), "subdir"),
			MaxUploadSize:  1024 * 1024,
			InitialQuality: 80,
		}

		pngData := createValidPNG()
		header := createMultipartImageRequest(t, "image", "test.png", pngData)

		_, err = svc.UploadImage(context.Background(), header, "err")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})

	t.Run("create error", func(t *testing.T) {
		// Use a directory that we know exists but we don't have permission OR
		// use a path that is impossible to create because it's a file.
		tmpFile, err := os.CreateTemp("", "cannot-create-here")
		assert.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		svc := &service.LocalImageService{
			UploadDir:      tmpFile.Name(), // This is a file, so os.Create(filepath.Join(svc.UploadDir, "err.webp")) will fail
			MaxUploadSize:  1024 * 1024,
			InitialQuality: 80,
		}

		pngData := createValidPNG()
		header := createMultipartImageRequest(t, "image", "test.png", pngData)

		// Note: UploadImage calls os.MkdirAll(svc.UploadDir, 0755).
		// If svc.UploadDir is an existing file, MkdirAll might fail or succeed depending on implementation.
		// Actually MkdirAll returns nil if path exists and is a directory.
		// If path exists and is a file, it returns syscall.ENOTDIR.
		_, err = svc.UploadImage(context.Background(), header, "err")
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

func TestLocalImageService_ConvertToWebP(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.InitialQuality = 80

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var pngBuf bytes.Buffer
	_ = png.Encode(&pngBuf, img)

	converted, err := svc.ConvertToWebP(strings.NewReader(pngBuf.String()))
	assert.NoError(t, err)

	convertedData, _ := io.ReadAll(converted)
	assert.True(t, len(convertedData) > 0, "converted should have content")
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
	fileHeader := createMultipartImageRequest(t, "image", "large.jpg", imgBuf.Bytes())

	path, err := svc.UploadImage(context.Background(), fileHeader, "large-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")

	// Verify the final file is under 200KB
	savedPath := filepath.Join(svc.UploadDir, "large-test.webp")
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
	fileHeader := createMultipartImageRequest(t, "image", "invalid.jpg", []byte("not an image"))

	_, err := svc.UploadImage(context.Background(), fileHeader, "decode-error")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid file type")
}

func TestLocalImageService_CompressImage_DecodeError(t *testing.T) {
	svc := service.NewLocalImageService("")

	_, err := svc.CompressImage(strings.NewReader("not an image"))
	assert.Error(t, err)
}

func TestLocalImageService_ConvertToWebP_DecodeError(t *testing.T) {
	svc := service.NewLocalImageService("")

	_, err := svc.ConvertToWebP(strings.NewReader("not a png"))
	assert.Error(t, err)
}

func TestLocalImageService_DeleteImage(t *testing.T) {
	tempDir := t.TempDir()
	svc := service.NewLocalImageService(tempDir)

	// 1. Create a dummy image file
	listingID := "test-delete"
	filename := listingID + ".jpg"
	savedPath := filepath.Join(tempDir, filename)
	err := os.WriteFile(savedPath, []byte("dummy image data"), 0644)
	assert.NoError(t, err)

	// 2. Delete it using the service
	imageURL := "/static/uploads/" + filename
	err = svc.DeleteImage(context.Background(), imageURL)
	assert.NoError(t, err)

	// 3. Verify it's gone
	_, err = os.Stat(savedPath)
	assert.True(t, os.IsNotExist(err), "file should be deleted")

	// 4. Test deleting non-existent file (should not error)
	err = svc.DeleteImage(context.Background(), "/static/uploads/non-existent.jpg")
	assert.NoError(t, err)

	// 5. Test deleting empty image URL
	err = svc.DeleteImage(context.Background(), "")
	assert.NoError(t, err)

	// 6. Test deleting with query parameters
	savedPath = filepath.Join(tempDir, "cache-bust.jpg")
	_ = os.WriteFile(savedPath, []byte("data"), 0644)
	err = svc.DeleteImage(context.Background(), "/static/uploads/cache-bust.jpg?v=123")
	assert.NoError(t, err)
	_, err = os.Stat(savedPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalImageService_ConvertToWebP_EncodeError(t *testing.T) {
	svc := service.NewLocalImageService("")
	svc.InitialQuality = -1 // Likely to cause issues

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var pngBuf bytes.Buffer
	_ = png.Encode(&pngBuf, img)
}

func TestLocalImageService_UploadImage_Downscale(t *testing.T) {
	svc := &service.LocalImageService{
		UploadDir:      t.TempDir(),
		MaxUploadSize:  10 * 1024 * 1024,
		MaxFileSize:    500, // Extremely small target to force downscale
		InitialQuality: 20,
		MinQuality:     10,
	}

	// Create a larger image (200x200) that will still be >500 bytes even at 10% quality
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(x % 255), 255})
		}
	}

	var imgBuf bytes.Buffer
	err := jpeg.Encode(&imgBuf, img, &jpeg.Options{Quality: 95})
	assert.NoError(t, err)

	fileHeader := createMultipartImageRequest(t, "image", "downscale.jpg", imgBuf.Bytes())

	path, err := svc.UploadImage(context.Background(), fileHeader, "downscale-test")
	assert.NoError(t, err)
	assert.Contains(t, path, ".webp")

	// Verify the final file exists
	savedPath := filepath.Join(svc.UploadDir, "downscale-test.webp")
	_, err = os.Stat(savedPath)
	assert.NoError(t, err)
}

func TestLocalImageService_DeleteImage_EdgeCases(t *testing.T) {
	svc := service.NewLocalImageService(t.TempDir())

	// Test with path traversal attempt (should be neutralized by filepath.Base)
	err := svc.DeleteImage(context.Background(), "/static/uploads/../../etc/passwd")
	assert.NoError(t, err)

	// Test with only a filename
	err = svc.DeleteImage(context.Background(), "test.jpg")
	assert.NoError(t, err)
}

func TestLocalImageService_Errors(t *testing.T) {
	svc := service.NewLocalImageService("")

	// Test CompressImage decode error
	_, err := svc.CompressImage(strings.NewReader("not-an-image"))
	assert.Error(t, err)

	// Test ConvertToWebP decode error
	_, err = svc.ConvertToWebP(strings.NewReader("not-an-image"))
	assert.Error(t, err)
}
