package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func setupCSVTest(t *testing.T) (*CSVService, context.Context, domain.ListingRepository) {
	t.Helper()
	svc := NewCSVService()
	ctx := context.Background()
	repo := testutil.SetupTestRepository(t)
	return svc, ctx, repo
}

func TestParseAndImport_Table(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		csv      string
		total    int
		success  int
		failures int
		wantErr  bool
	}{
		{
			name: "Valid Listings",
			csv: `title,type,description,origin,email,website
Test Biz,Business,Desc 1,Ghana,test@test.com,example.com
Test Svc,Service,Desc 2,Nigeria,svc@test.com,
`,
			total:    2,
			success:  2,
			failures: 0,
		},
		{
			name: "Partial Failures",
			csv: `title,type,description,origin,email,phone,website
Good,Business,Desc,Ghana,a@b.com,,
BadTypeDefaultsToBusiness,InvalidType,Desc,Ghana,a@b.com,,
MissingDesc,Business,,Ghana,a@b.com,,
,Business,Desc,Ghana,a@b.com,,
MissingOrigin,Business,Desc,,a@b.com,,
MissingEmailHasPhone,Business,Desc,Ghana,,12345,
MissingAllContact,Business,Desc,Ghana,,,
MissingEmailPhoneHasWebsite,Business,Desc,Ghana,,,example.com`,
			total:    8,
			success:  5,
			failures: 3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc, ctx, repo := setupCSVTest(t)
			result, err := svc.ParseAndImport(ctx, strings.NewReader(tt.csv), repo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.total, result.TotalProcessed)
			assert.Equal(t, tt.success, result.SuccessCount)
			assert.Equal(t, tt.failures, result.FailureCount)
		})
	}
}

func TestParseAndImport_MissingHeaders(t *testing.T) {
	t.Parallel()
	svc, ctx, repo := setupCSVTest(t)
	csvContent := `title,description
Test,Desc`
	_, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required header")
}

func TestParseAndImport_Duplicate(t *testing.T) {
	t.Parallel()
	svc := NewCSVService()
	ctx := context.Background()
	csvContent := `title,type,description,origin,email,phone,address,city
Dup Listing,Business,Same Desc,Ghana,dup@test.com,1234,123 St,Accra
`
	repo := testutil.SetupTestRepository(t)
	// Seed duplicate
	_ = repo.Save(ctx, domain.Listing{Title: "Dup Listing", Type: domain.Business, Description: "Same Desc", ContactEmail: "dup@test.com", OwnerOrigin: "Nigeria"})

	result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.FailureCount)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Contains(t, result.Errors[0], ">2 fields match")
}

func TestParseAndImport_CategoryEdgeCases(t *testing.T) {
	t.Parallel()
	svc := NewCSVService()
	ctx := context.Background()
	csvContent := `title,type,description,email
Lowercase Food,food,Desc,a@b.com
Random Type,Random,Desc,a@b.com
Dynamic Church,church,Desc,a@b.com
UPPERCASE CHURCH,CHURCH,Desc,a@b.com
`
	repo := testutil.SetupTestRepository(t)

	result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.FailureCount)

	listings, err := repo.FindByTitle(ctx, "Dynamic Church")
	assert.NoError(t, err)
	assert.NotEmpty(t, listings)
	assert.Equal(t, domain.Category("Church"), listings[0].Type)

	listingsUpper, err := repo.FindByTitle(ctx, "UPPERCASE CHURCH")
	assert.NoError(t, err)
	assert.NotEmpty(t, listingsUpper)
	assert.Equal(t, domain.Category("Church"), listingsUpper[0].Type)
}

func TestParseAndImport_Geocoding(t *testing.T) {
	t.Parallel()
	svc := NewCSVService()
	ctx := context.Background()
	csvContent := `title,type,description,email,address
Geo Hub,Business,Test Geocode,test@geo.com,"1600 Amphitheatre Parkway, Mountain View, CA"
`
	repo := testutil.SetupTestRepository(t)
	mockGeo := &mockGeocodingService{
		GetCityFunc: func(ctx context.Context, addr string) (string, error) {
			if strings.Contains(addr, "Mountain View") {
				return "Mountain View", nil
			}
			return "", nil
		},
	}
	svc.Geocoding = mockGeo

	result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)

	// Verify city was populated
	listings, _ := repo.FindByTitle(ctx, "Geo Hub")
	assert.Equal(t, "Mountain View", listings[0].City)
}

type mockGeocodingService struct {
	GetCityFunc func(ctx context.Context, address string) (string, error)
	GeocodeFunc func(ctx context.Context, address string) (float64, float64, error)
}

func (m *mockGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	return m.GetCityFunc(ctx, address)
}

func (m *mockGeocodingService) Geocode(ctx context.Context, address string) (float64, float64, error) {
	if m.GeocodeFunc != nil {
		return m.GeocodeFunc(ctx, address)
	}
	return 0, 0, nil
}

func FuzzParseAndImport(f *testing.F) {
	svc := NewCSVService()
	ctx := context.Background()

	// Seed some inputs
	f.Add("title,type,description,origin,email\nTest Biz,Business,Description,Nigeria,test@test.com")
	f.Add("title,type,description\nTest Biz,Business,Description")
	f.Add("invalid csv data")
	f.Add("")

	f.Fuzz(func(t *testing.T, data string) {
		repo := testutil.SetupTestRepository(t)
		_, _ = svc.ParseAndImport(ctx, strings.NewReader(data), repo)
	})
}

func FuzzParseCategory(f *testing.F) {
	f.Add("Business")
	f.Add("food")
	f.Add("Random")
	f.Add("")
	f.Add("🌟")
	f.Fuzz(func(t *testing.T, typeStr string) {
		cat := parseCategory(typeStr)
		if cat == "" {
			t.Error("category should never be empty")
		}
	})
}

func TestGenerateCSV_Empty(t *testing.T) {
	t.Parallel()
	svc := NewCSVService()
	ctx := context.Background()

	reader, err := svc.GenerateCSV(ctx, []domain.Listing{})
	assert.NoError(t, err)

	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	assert.NoError(t, err)
	assert.Equal(t, "ID", headers[0])

	_, err = csvReader.Read()
	assert.Equal(t, io.EOF, err)
}

func TestGenerateCSV_FullMapping(t *testing.T) {
	t.Parallel()
	svc := NewCSVService()
	ctx := context.Background()
	now := time.Now().Truncate(time.Second)

	listings := []domain.Listing{
		{
			ID:                    "full-1",
			Title:                 "Full Title",
			Type:                  domain.Business,
			Description:           "Full Desc",
			City:                  "City",
			Address:               "Addr",
			OwnerOrigin:           "Origin",
			ContactEmail:          "e@e.com",
			ContactPhone:          "123",
			ContactWhatsApp:       "456",
			WebsiteURL:            "w.com",
			CreatedAt:             now,
			Status:                domain.ListingStatusApproved,
			IsActive:              true,
			Featured:              true,
			Company:               "Co",
			PayRange:              "1-2",
			Skills:                "Go",
			JobApplyURL:           "j.com",
			JobStartDate:          now,
			EventStart:            now,
			EventEnd:              now.Add(time.Hour),
			Deadline:              now.Add(24 * time.Hour),
			EnrichmentAttemptedAt: &now,
		},
	}

	reader, err := svc.GenerateCSV(ctx, listings)
	assert.NoError(t, err)

	csvReader := csv.NewReader(reader)
	_, _ = csvReader.Read() // Skip header
	row, err := csvReader.Read()
	assert.NoError(t, err)
	assert.Equal(t, "full-1", row[0])
	assert.Equal(t, "Full Title", row[1])
	assert.Equal(t, "Business", row[2])
	assert.Equal(t, "Full Desc", row[3])
	assert.Equal(t, now.Format(time.RFC3339), row[23])
}

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("forced read error")
}

func TestParseAndImport_ReadError(t *testing.T) {
	t.Parallel()
	svc := NewCSVService()
	ctx := context.Background()
	repo := testutil.SetupTestRepository(t)

	_, err := svc.ParseAndImport(ctx, &errReader{}, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read CSV header")
}
