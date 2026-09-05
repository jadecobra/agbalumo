package ui_test

import (
	"bytes"
	"testing"

	"github.com/jadecobra/agbalumo/internal/module"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadMeta_OpenGraphTags(t *testing.T) {
	renderer, err := ui.NewTemplateRenderer(
		"../../ui/templates/*.html",
		"../../ui/templates/partials/*.html",
		"../../ui/templates/components/*.html",
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		data     module.BaseViewData
		contains []string
	}{
		{
			name: "default_fallback_tags",
			data: module.BaseViewData{
				CSRF: "test-csrf",
			},
			contains: []string{
				`<meta property="og:type" content="website"`,
				`<meta property="og:site_name" content="agbalumo"`,
				`<meta property="og:title" content="agbalumo - Find African Food in <60s"`,
				`<meta property="og:description" content="Find authentic Nigerian and West African food in under 60 seconds. Verified Dallas spots, directions, and direct contact."`,
				`<meta property="og:url" content="https://agbalumo.com"`,
				`<meta property="og:image" content="https://agbalumo.com/static/icons/favicon.png"`,
				`<meta name="twitter:card" content="summary_large_image"`,
				`<meta name="twitter:title" content="agbalumo - Find African Food in <60s"`,
				`<meta name="twitter:description" content="Find authentic Nigerian and West African food in under 60 seconds. Verified Dallas spots, directions, and direct contact."`,
				`<meta name="twitter:image" content="https://agbalumo.com/static/icons/favicon.png"`,
			},
		},
		{
			name: "custom_meta_tags",
			data: module.BaseViewData{
				CSRF:            "test-csrf",
				MetaTitle:       "Authentic Suya in Dallas | agbalumo",
				MetaDescription: "Top 3 suya spots in Dallas rated for authentic spice and minimal wait time.",
				MetaImage:       "https://agbalumo.com/static/uploads/suya.jpg",
				MetaURL:         "https://agbalumo.com/?q=suya&city=Dallas",
				MetaType:        "restaurant",
			},
			contains: []string{
				`<meta property="og:type" content="restaurant"`,
				`<meta property="og:title" content="Authentic Suya in Dallas | agbalumo"`,
				`<meta property="og:description" content="Top 3 suya spots in Dallas rated for authentic spice and minimal wait time."`,
				`<meta property="og:url" content="https://agbalumo.com/?q=suya&amp;city=Dallas"`,
				`<meta property="og:image" content="https://agbalumo.com/static/uploads/suya.jpg"`,
				`<meta name="twitter:title" content="Authentic Suya in Dallas | agbalumo"`,
				`<meta name="twitter:description" content="Top 3 suya spots in Dallas rated for authentic spice and minimal wait time."`,
				`<meta name="twitter:image" content="https://agbalumo.com/static/uploads/suya.jpg"`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := renderer.RenderDefinition(buf, "head_meta", tc.data)
			require.NoError(t, err)

			output := buf.String()
			for _, expected := range tc.contains {
				assert.Contains(t, output, expected)
			}
		})
	}
}
