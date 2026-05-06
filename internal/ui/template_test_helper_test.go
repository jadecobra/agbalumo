package ui_test

import (
	"bytes"
	"testing"

	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/stretchr/testify/require"
)

// renderPartial loads the real project templates and executes a named partial with data.
// This enables testing partials in isolation without a server.
func renderPartial(t *testing.T, name string, data interface{}) string {
	renderer, err := ui.NewTemplateRenderer("../../ui/templates/*.html", "../../ui/templates/partials/*.html", "../../ui/templates/components/*.html")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = renderer.RenderDefinition(&buf, name, data)
	require.NoError(t, err)
	return buf.String()
}
