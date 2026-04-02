package service_test

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalImageService_CompressImage(t *testing.T) {
	svc, _ := setupTestImageService(t)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var originalBuf bytes.Buffer
	_ = jpeg.Encode(&originalBuf, img, nil)

	compressed, err := svc.CompressImage(strings.NewReader(originalBuf.String()))
	assert.NoError(t, err)

	compressedData, _ := io.ReadAll(compressed)
	assert.True(t, len(compressedData) > 0, "compressed should have content")
}

func TestLocalImageService_ConvertToWebP(t *testing.T) {
	svc, _ := setupTestImageService(t)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var pngBuf bytes.Buffer
	_ = png.Encode(&pngBuf, img)

	converted, err := svc.ConvertToWebP(strings.NewReader(pngBuf.String()))
	assert.NoError(t, err)

	convertedData, _ := io.ReadAll(converted)
	assert.True(t, len(convertedData) > 0, "converted should have content")
}

func TestLocalImageService_Errors(t *testing.T) {
	svc, _ := setupTestImageService(t)

	t.Run("CompressImage decode error", func(t *testing.T) {
		_, err := svc.CompressImage(strings.NewReader("not-an-image"))
		assert.Error(t, err)
	})

	t.Run("ConvertToWebP decode error", func(t *testing.T) {
		_, err := svc.ConvertToWebP(strings.NewReader("not-an-image"))
		assert.Error(t, err)
	})
}
