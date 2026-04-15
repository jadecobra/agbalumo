package service

import (
	"context"
	"github.com/jadecobra/agbalumo/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCoverageBoost_ImageService(t *testing.T) {
	t.Parallel()
	t.Run("NewLocalImageService empty path", func(t *testing.T) {
		t.Parallel()
		svc := NewLocalImageService("")
		assert.Equal(t, "ui/static/uploads", svc.UploadDir)
	})

	t.Run("DeleteImage variations", func(t *testing.T) {
		t.Parallel()
		svc, _ := setupTestImageService(t, nil)

		// Empty URL
		err := svc.DeleteImage(context.Background(), "")
		assert.NoError(t, err)

		// URL with query param
		err = svc.DeleteImage(context.Background(), "/static/uploads/test.webp?v=123")
		assert.NoError(t, err)

		// Non-existent file
		err = svc.DeleteImage(context.Background(), "/static/uploads/non-existent.webp")
		assert.NoError(t, err)
	})
}

func TestCoverageBoost_Geocoding(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK","results":[{"address_components":[{"long_name":"Test City","types":["locality"]}]}]}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewGoogleGeocodingService("fake-key")
	svc.BaseURL = server.URL

	// Valid address (mocked)
	city, err := svc.GetCity(context.Background(), "123 Test St")
	assert.NoError(t, err)
	assert.Equal(t, "Test City", city)

	// ZERO_RESULTS
	mux.HandleFunc("/zero", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ZERO_RESULTS"}`))
	})
	svc.BaseURL = server.URL + "/zero"
	city, err = svc.GetCity(context.Background(), "nowhere")
	assert.NoError(t, err)
	assert.Equal(t, "", city)
}

func TestCoverageBoost_Background(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)
	svc := NewBackgroundService(repo, nil)
	svc.Interval = 1 * time.Millisecond // Fast ticker

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	go func() {
		svc.StartTicker(ctx)
		done <- true
	}()

	// Wait for at least one tick and then stop via context timeout
	<-done
}
