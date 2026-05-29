package listing

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/ui"
	"github.com/stretchr/testify/assert"
)

const testInvalid = "invalid"

func TestParseDeadline(t *testing.T) {
	t.Parallel()
	req := &ListingFormRequest{DeadlineDate: "2024-12-31"}
	l := &domain.Listing{Type: domain.Request}

	err := parseDeadline(req, l)
	assert.NoError(t, err)
	assert.Equal(t, 2024, l.Deadline.Year())

	req.DeadlineDate = testInvalid
	err = parseDeadline(req, l)
	assert.Error(t, err)
}

func TestParseEventDates(t *testing.T) {
	t.Parallel()
	req := &ListingFormRequest{
		EventStart: "2024-12-01T10:00",
		EventEnd:   "2024-12-01T12:00",
	}
	l := &domain.Listing{Type: domain.Event}

	err := parseEventDates(req, l)
	assert.NoError(t, err)
	assert.Equal(t, 2024, l.EventStart.Year())

	req.EventStart = testInvalid
	err = parseEventDates(req, l)
	assert.Error(t, err)

	req.EventStart = "2024-12-01T10:00"
	req.EventEnd = testInvalid
	err = parseEventDates(req, l)
	assert.Error(t, err)
}

func TestParseJobStartDate(t *testing.T) {
	t.Parallel()
	req := &ListingFormRequest{JobStartDate: "2024-12-01T09:00"}
	l := &domain.Listing{Type: domain.Job}

	err := parseJobStartDate(req, l)
	assert.NoError(t, err)
	assert.Equal(t, 2024, l.JobStartDate.Year())

	req.JobStartDate = testInvalid
	err = parseJobStartDate(req, l)
	assert.Error(t, err)
}

func TestListingFormDefaultType(t *testing.T) {
	t.Parallel()
	renderer, err := ui.NewTemplateRenderer("../../../ui/templates/*.html")
	assert.NoError(t, err)

	tmpl := template.New("listing_form_type_origin.html").Funcs(renderer.GetFuncMap())
	_, err = tmpl.ParseFiles(
		"../../../ui/templates/components/listing_form_type_origin.html",
		"../../../ui/templates/components/custom_country_options.html",
	)
	assert.NoError(t, err)

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Categories": []domain.CategoryData{
			{Name: "Food"},
			{Name: "Business"},
		},
		"SelectedType": "",
	}
	err = tmpl.ExecuteTemplate(&buf, "listing_form_type_origin", data)
	assert.NoError(t, err)

	html := buf.String()
	// Should default to "Food" instead of "Business"
	assert.Contains(t, html, `value="Food"`)
	assert.Contains(t, html, `<span class="dropdown-display pointer-events-none">Food</span>`)
}
