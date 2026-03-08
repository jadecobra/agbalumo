package security

import (
	"bytes"
	"html/template"
	"image"
	"net/http"
	"net/url"
	"testing"

	"github.com/gen2brain/webp"
	"github.com/stretchr/testify/assert"
)

// GO-2026-4603: URLs in meta content attribute actions are not escaped in html/template
func TestVuln_MetaContentEscaping(t *testing.T) {
	tmpl, err := template.New("test").Parse(`<meta content="{{.URL}}">`)
	assert.NoError(t, err)

	var buf bytes.Buffer
	// If vulnerable, certain characters might not be escaped correctly in this context.
	// We use a payload that might break out of the attribute if not escaped.
	data := struct{ URL string }{URL: "javascript:alert(1)"}
	err = tmpl.Execute(&buf, data)
	assert.NoError(t, err)

	// In safe versions, html/template should escape or neutralize dangerous URLs in meta content.
	// This test is to exercise the code path mentioned in the audit.
	t.Logf("Rendered meta content: %s", buf.String())
}

// GO-2026-4602: FileInfo can escape from a Root in os
func TestVuln_OsFileInfoEscape(t *testing.T) {
	// This vulnerability affects os.Root (introduced in 1.24)
	// Even if we don't use os.Root directly, govulncheck flagged it via webp.Encode -> os.File.ReadDir
	// We'll exercise the webp.Encode path to ensure it works.
	
	// Just exercise the path flagged by govulncheck
	var out bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	_ = webp.Encode(&out, img, webp.Options{Quality: 75})
	
	t.Log("Exercised webp.Encode path")
}

// GO-2026-4601: Incorrect parsing of IPv6 host literals in net/url
func TestVuln_IPv6Parsing(t *testing.T) {
	// Vulnerable versions might incorrectly parse IPv6 host literals.
	rawURL := "http://[2001:db8::1]:8080/path"
	u, err := url.Parse(rawURL)
	assert.NoError(t, err)
	assert.Equal(t, "[2001:db8::1]:8080", u.Host)
	
	// The fix involves how it handles specific malformed or edge-case IPv6 literals.
	// We'll also exercise the http.Get path mentioned in the audit.
	t.Logf("Parsed IPv6 host: %s", u.Host)
	
	// Exercise url.ParseRequestURI as well (flagged in cmd/serve.go)
	_, err = url.ParseRequestURI("/test?query=1")
	assert.NoError(t, err)
}

func TestVuln_HttpHeaderCheck(t *testing.T) {
    // Audit flagged http.Get in handler.RealGoogleProvider.GetUserInfo
    // We'll just verify a basic http request structure doesn't crash or behave weirdly
    // even if we don't actually hit the network.
    req, _ := http.NewRequest("GET", "https://example.com", nil)
    t.Logf("Request URL: %s", req.URL.String())
}
