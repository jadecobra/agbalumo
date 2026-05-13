package maintenance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckHitboxes_Small(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprint(w, `
			<html>
				<body>
					<button style="width: 20px; height: 20px;">Small</button>
					<button style="width: 50px; height: 50px;">Large</button>
				</body>
			</html>
		`)
	}))
	defer ts.Close()

	violations, err := CheckHitboxes(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	foundSmall := false
	for _, v := range violations {
		if v.Reason == "Touch target too small (min 44x44)" {
			foundSmall = true
			break
		}
	}
	assert.True(t, foundSmall, "Should have found a small hitbox violation")
}

func TestCheckHitboxes_Blocked(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprint(w, `
			<html>
				<body>
					<button style="width: 100px; height: 100px; position: absolute; top: 0; left: 0;">Click Me</button>
					<div class="overlay" style="position: absolute; top: 0; left: 0; width: 100px; height: 100px; z-index: 10;"></div>
				</body>
			</html>
		`)
	}))
	defer ts.Close()

	violations, err := CheckHitboxes(ts.URL)
	assert.NoError(t, err)

	foundBlocked := false
	for _, v := range violations {
		if v.Reason == "Interaction blocked by overlay" {
			foundBlocked = true
			break
		}
	}
	assert.True(t, foundBlocked, "Should have found a blocked interaction violation")
}
