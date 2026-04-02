package service_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jadecobra/agbalumo/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestLocalImageService_DeleteImage(t *testing.T) {
	tempDir := t.TempDir()
	svc := service.NewLocalImageService(tempDir)

	listingID := "test-delete"
	filename := listingID + ".jpg"
	savedPath := filepath.Join(tempDir, filename)
	err := os.WriteFile(savedPath, []byte("dummy image data"), 0644)
	assert.NoError(t, err)

	imageURL := "/static/uploads/" + filename
	err = svc.DeleteImage(context.Background(), imageURL)
	assert.NoError(t, err)

	_, err = os.Stat(savedPath)
	assert.True(t, os.IsNotExist(err), "file should be deleted")

	err = svc.DeleteImage(context.Background(), "/static/uploads/non-existent.jpg")
	assert.NoError(t, err)
}

func TestLocalImageService_DeleteImage_EdgeCases(t *testing.T) {
	svc := service.NewLocalImageService(t.TempDir())

	err := svc.DeleteImage(context.Background(), "/static/uploads/../../etc/passwd")
	assert.NoError(t, err)

	err = svc.DeleteImage(context.Background(), "test.jpg")
	assert.NoError(t, err)
}
