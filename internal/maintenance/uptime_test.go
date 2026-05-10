package maintenance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckServerUptime(t *testing.T) {
	t.Run("server is up", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		err := CheckServerUptime(ts.URL)
		assert.NoError(t, err)
	})

	t.Run("server returns error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		err := CheckServerUptime(ts.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-OK status")
	})

	t.Run("server is unreachable", func(t *testing.T) {
		err := CheckServerUptime("http://localhost:1") // highly likely to be unreachable
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unreachable")
	})

	t.Run("empty URL", func(t *testing.T) {
		err := CheckServerUptime("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL is empty")
	})
}
