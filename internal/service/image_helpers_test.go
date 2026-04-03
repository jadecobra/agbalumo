package service

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createValidPNG() []byte {
	return createCustomPNG(10, 10)
}

func createCustomPNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some noise to make it harder to compress
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), uint8((x + y) % 256), 255})
		}
	}
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

func createValidGIF() []byte {
	img := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{color.White})
	var buf bytes.Buffer
	_ = gif.EncodeAll(&buf, &gif.GIF{
		Image: []*image.Paletted{img},
		Delay: []int{0},
	})
	return buf.Bytes()
}

func setupTestImageService(t *testing.T, mutators ...func(*LocalImageService)) (*LocalImageService, string) {
	t.Helper()
	tempDir := t.TempDir()
	svc := &LocalImageService{
		UploadDir:      tempDir,
		MaxUploadSize:  1024 * 1024,
		MaxFileSize:    200 * 1024,
		InitialQuality: 80,
		MinQuality:     20,
	}
	for _, m := range mutators {
		if m != nil {
			m(svc)
		}
	}
	return svc, svc.UploadDir
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
