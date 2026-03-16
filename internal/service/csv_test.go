package service

import (
	"context"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
)


func TestParseAndImport(t *testing.T) {
	svc := NewCSVService()
	ctx := context.Background()

	t.Run("Valid Import", func(t *testing.T) {
		csvContent := `title,type,description,origin,email,website
Test Biz,Business,Desc 1,Ghana,test@test.com,example.com
Test Svc,Service,Desc 2,Nigeria,svc@test.com,
`
		repo := setupTestRepo(t)

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.TotalProcessed != 2 {
			t.Errorf("Expected 2 total, got %d", result.TotalProcessed)
		}
		if result.SuccessCount != 2 {
			t.Errorf("Expected 2 success, got %d", result.SuccessCount)
		}
		if result.FailureCount != 0 {
			t.Errorf("Expected 0 failures, got %d", result.FailureCount)
		}
	})

	t.Run("Missing Headers", func(t *testing.T) {
		csvContent := `title,description
Test,Desc`
		repo := setupTestRepo(t)
		_, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err == nil {
			t.Error("Expected error for missing headers, got nil")
		}
		if !strings.Contains(err.Error(), "missing required header") {
			t.Errorf("Expected missing header error, got %v", err)
		}
	})

	t.Run("Partial Failure", func(t *testing.T) {
		csvContent := `title,type,description,origin,email,phone,website
Good,Business,Desc,Ghana,a@b.com,,
BadTypeDefaultsToBusiness,InvalidType,Desc,Ghana,a@b.com,,
MissingDesc,Business,,Ghana,a@b.com,,
,Business,Desc,Ghana,a@b.com,,
MissingOrigin,Business,Desc,,a@b.com,,
MissingEmailHasPhone,Business,Desc,Ghana,,12345,
MissingAllContact,Business,Desc,Ghana,,,
MissingEmailPhoneHasWebsite,Business,Desc,Ghana,,,example.com`

		repo := setupTestRepo(t)

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.TotalProcessed != 8 {
			t.Errorf("Expected 8 total, got %d", result.TotalProcessed)
		}
		if result.SuccessCount != 5 {
			t.Errorf("Expected 5 success, got %d", result.SuccessCount)
		}
		if result.FailureCount != 3 {
			t.Errorf("Expected 3 failures, got %d", result.FailureCount)
		}
	})

	t.Run("Duplicate Check", func(t *testing.T) {
		csvContent := `title,type,description,origin,email,phone,address,city
Dup Listing,Business,Same Desc,Ghana,dup@test.com,1234,123 St,Accra
`
		repo := setupTestRepo(t)
		// Seed duplicate
		_ = repo.Save(ctx, domain.Listing{Title: "Dup Listing", Type: domain.Business, Description: "Same Desc", ContactEmail: "dup@test.com", OwnerOrigin: "Nigeria"})

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.FailureCount != 1 {
			t.Errorf("Expected 1 failure, got %d. Errors: %v", result.FailureCount, result.Errors)
		}
		if result.SuccessCount != 0 {
			t.Errorf("Expected 0 success, got %d", result.SuccessCount)
		}
		if !strings.Contains(result.Errors[0], ">2 fields match") {
			t.Errorf("Expected >2 fields match error, got %s", result.Errors[0])
		}
	})

	// Duplicate DB Error and Save Error are hard to force with SQLite without mocks.
	// Since we're moving towards integration tests, we'll skip these or use real constraints if possible.

	t.Run("Category Parsing Edge Cases", func(t *testing.T) {
		csvContent := `title,type,description,email
Lowercase Food,food,Desc,a@b.com
Random Type,Random,Desc,a@b.com
Dynamic Church,church,Desc,a@b.com
UPPERCASE CHURCH,CHURCH,Desc,a@b.com
`
		repo := setupTestRepo(t)

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.FailureCount != 0 {
			t.Errorf("Expected 0 failures, got %d", result.FailureCount)
		}

		listings, err := repo.FindByTitle(ctx, "Dynamic Church")
		if err != nil || len(listings) == 0 {
			t.Fatalf("Failed to find Dynamic Church listing")
		}
		if listings[0].Type != domain.Category("Church") {
			t.Errorf("Expected category 'Church', got %q", listings[0].Type)
		}

		listingsUpper, err := repo.FindByTitle(ctx, "UPPERCASE CHURCH")
		if err != nil || len(listingsUpper) == 0 {
			t.Fatalf("Failed to find UPPERCASE CHURCH listing")
		}
		if listingsUpper[0].Type != domain.Category("Church") {
			t.Errorf("Expected category 'Church', got %q", listingsUpper[0].Type)
		}
	})

	t.Run("Geocoding Fallback", func(t *testing.T) {
		csvContent := `title,type,description,email,address
Geo Hub,Business,Test Geocode,test@geo.com,"1600 Amphitheatre Parkway, Mountain View, CA"
`
		repo := setupTestRepo(t)
		mockGeo := &mockGeocodingService{
			GetCityFunc: func(ctx context.Context, addr string) (string, error) {
				if strings.Contains(addr, "Mountain View") {
					return "Mountain View", nil
				}
				return "", nil
			},
		}
		svcWithGeo := NewCSVService()
		svcWithGeo.Geocoding = mockGeo

		result, err := svcWithGeo.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.SuccessCount != 1 {
			t.Fatalf("Expected 1 success, got %d. Errors: %v", result.SuccessCount, result.Errors)
		}

		// Verify city was populated
		listings, _ := repo.FindByTitle(ctx, "Geo Hub")
		if listings[0].City != "Mountain View" {
			t.Errorf("Expected city 'Mountain View', got %q", listings[0].City)
		}
	})
}

type mockGeocodingService struct {
	GetCityFunc func(ctx context.Context, address string) (string, error)
}

func (m *mockGeocodingService) GetCity(ctx context.Context, address string) (string, error) {
	return m.GetCityFunc(ctx, address)
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
		repo := setupTestRepo(t)
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

func TestGenerateCSV(t *testing.T) {
	svc := NewCSVService()
	ctx := context.Background()

	listings := []domain.Listing{
		{
			ID:           "test-1",
			Title:        "Listing 1",
			Type:         domain.Business,
			Description:  "Desc 1",
			City:         "Lagos",
			ContactEmail: "l1@example.com",
		},
		{
			ID:           "test-2",
			Title:        "Listing 2",
			Type:         domain.Job,
			Description:  "Desc 2",
			Company:      "Tech Co",
			ContactEmail: "l2@example.com",
		},
	}

	reader, err := svc.GenerateCSV(ctx, listings)
	if err != nil {
		t.Fatalf("GenerateCSV failed: %v", err)
	}

	// Read and verify
	importSvc := NewCSVService()
	// Since GenerateCSV is a stream, we can read it all
	// We'll verify it by counting lines or checking content
	importResult, err := importSvc.ParseAndImport(ctx, reader, setupTestRepo(t))
	if err != nil {
		t.Fatalf("Failed to parse generated CSV: %v", err)
	}

	if importResult.SuccessCount != len(listings) {
		t.Errorf("Expected %d successful imports, got %d. Errors: %v", len(listings), importResult.SuccessCount, importResult.Errors)
	}
}
