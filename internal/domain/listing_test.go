package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateDeadline(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		deadline time.Time
		typeStr  Category
		wantErr  error
	}{
		{
			name:     "Valid Deadline (89 days)",
			deadline: now.Add(89 * 24 * time.Hour),
			typeStr:  Request,
			wantErr:  nil,
		},
		{
			name:     "Valid Deadline (Exactly 90 days)",
			deadline: now.Add(90 * 24 * time.Hour),
			typeStr:  Request,
			wantErr:  nil,
		},
		{
			name:     "Invalid Deadline (90 days + 1 second)",
			deadline: now.Add(90*24*time.Hour + time.Second),
			typeStr:  Request,
			wantErr:  ErrInvalidDeadline,
		},
		{
			name:     "Ignore Deadline for Business Type",
			deadline: now.Add(100 * 24 * time.Hour),
			typeStr:  Business,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Listing{
				ID:           "test-id",
				OwnerOrigin:  "Nigeria",
				Type:         tt.typeStr,
				Title:        "Test Title",
				ContactEmail: "test@example.com",
				Address:      "123 Valid St", // Satisfy address requirement for Business/Food
				CreatedAt:    now,
				Deadline:     tt.deadline,
				IsActive:     true,
			}
			err := l.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContactRequirement(t *testing.T) {
	listing := Listing{
		ID:          "2",
		OwnerOrigin: "Ghana",
		Type:        Business,
		Title:       "Jollof Place",
		CreatedAt:   time.Now(),
		IsActive:    true,
		Address:     "123 St",
	}

	// No contact info
	err := listing.Validate()
	assert.ErrorIs(t, err, ErrMissingContact)

	// With Email only
	listing.ContactEmail = "jollof@example.com"
	err = listing.Validate()
	assert.NoError(t, err)

	// With WhatsApp only
	listing.ContactEmail = ""
	listing.ContactWhatsApp = "+233555555555"
	err = listing.Validate()
	assert.NoError(t, err)

	// With Phone only (New)
	listing.ContactWhatsApp = ""
	listing.ContactPhone = "+2348000000000"
	err = listing.Validate()
	assert.NoError(t, err)
}

func TestOriginValidation(t *testing.T) {
	tests := []struct {
		name    string
		origin  string
		wantErr error
	}{
		{
			name:    "Valid Origin - Nigeria",
			origin:  "Nigeria",
			wantErr: nil,
		},
		{
			name:    "Valid Origin - Ghana",
			origin:  "Ghana",
			wantErr: nil,
		},
		{
			name:    "Valid Origin - Senegal",
			origin:  "Senegal",
			wantErr: nil,
		},
		{
			name:    "Invalid Origin - France",
			origin:  "France",
			wantErr: ErrInvalidOrigin,
		},
		{
			name:    "Invalid Origin - USA",
			origin:  "USA",
			wantErr: ErrInvalidOrigin,
		},
		{
			name:    "Empty Origin",
			origin:  "",
			wantErr: ErrMissingOrigin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Listing{
				ID:           "3",
				Type:         Product,
				Title:        "Shea Butter",
				ContactEmail: "shea@example.com",
				OwnerOrigin:  tt.origin,
				CreatedAt:    time.Now(),
			}
			err := l.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequestRequiresOrigin(t *testing.T) {
	// 10x Engineer Standard: Strict TDD - Requests MUST have an origin.
	// This test asserts that the exemption in the current code is removed.

	l := Listing{
		ID:           "req-1",
		Type:         Request,
		Title:        "Looking for Palm Wine",
		ContactEmail: "seeker@example.com",
		OwnerOrigin:  "", // Deliberately empty
		CreatedAt:    time.Now(),
		Deadline:     time.Now().Add(24 * time.Hour),
		IsActive:     true,
	}

	err := l.Validate()
	assert.ErrorIs(t, err, ErrMissingOrigin, "Requests must require an OwnerOrigin as per spec.md")
}

func BenchmarkValidate(b *testing.B) {
	l := Listing{
		ID:           "bench-1",
		OwnerOrigin:  "Nigeria",
		Type:         Business,
		Title:        "Benchmark Business",
		ContactEmail: "bench@example.com",
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Validate()
	}
}

func TestAddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		lType   Category
		address string
		wantErr bool
	}{
		{
			name:    "Business requires address",
			lType:   Business,
			address: "",
			wantErr: true,
		},
		{
			name:    "Business with address is valid",
			lType:   Business,
			address: "123 Market St",
			wantErr: false,
		},
		{
			name:    "Food requires address",
			lType:   Food,
			address: "",
			wantErr: true,
		},
		{
			name:    "Service does NOT require address",
			lType:   Service,
			address: "",
			wantErr: false,
		},
		{
			name:    "Service with address is also valid",
			lType:   Service,
			address: "Home Office",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Listing{
				ID:           "test-addr",
				OwnerOrigin:  "Ghana",
				Type:         tt.lType,
				Title:        "Test Biz",
				ContactEmail: "test@example.com",
				Address:      tt.address,
				CreatedAt:    time.Now(),
				IsActive:     true,
			}
			err := l.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHoursOfOperationField(t *testing.T) {
	// This test ensures the field exists and can be set.
	l := Listing{
		ID:               "test-hours",
		OwnerOrigin:      "Togo",
		Type:             Business,
		Title:            "Hours Test",
		ContactEmail:     "hours@example.com",
		Address:          "Main St",
		HoursOfOperation: "Mon-Fri 9-5",
		CreatedAt:        time.Now(),
	}
	
	assert.Equal(t, "Mon-Fri 9-5", l.HoursOfOperation)
}

func TestHoursOfOperationRestriction(t *testing.T) {
	tests := []struct {
		name     string
		lType    Category
		hours    string
		wantErr  bool
	}{
		{ "Business with hours", Business, "9-5", false },
		{ "Service with hours", Service, "9-5", false },
		{ "Food with hours", Food, "9-5", false },
		{ "Product with hours", Product, "9-5", true }, // Should fail
		{ "Event with hours", Event, "9-5", true },   // Should fail
		{ "Job with hours", Job, "9-5", true },       // Should fail
		{ "Request with hours", Request, "9-5", true }, // Should fail
		{ "Product without hours", Product, "", false },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Listing{
				ID:               "test-restrict",
				OwnerOrigin:      "Ghana",
				Type:             tt.lType,
				Title:            "Test",
				ContactEmail:     "test@example.com",
				// Satisfy other requirements
				Company:          "Acme", // for Job
				Skills:           "Go",   // for Job
				JobStartDate:     time.Now().Add(24*time.Hour), // for Job
				EventStart:       time.Now().Add(24*time.Hour), // for Event
				EventEnd:         time.Now().Add(25*time.Hour), // for Event
				Address:          "123 St", // for Business/Food
				HoursOfOperation: tt.hours,
				CreatedAt:        time.Now(),
			}
			err := l.Validate()
			if tt.wantErr {
				assert.Error(t, err, "Expected error for type %s with hours", tt.lType)
			} else {
				assert.NoError(t, err, "Expected no error for type %s with hours", tt.lType)
			}
		})
	}
}


func TestListing_Validate_Length(t *testing.T) {
	longString := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = 'a'
		}
		return string(b)
	}

	tests := []struct {
		name    string
		listing Listing
		wantErr bool
	}{
		{
			name: "Title too long (>100)",
			listing: Listing{
				ID:          "1",
				OwnerOrigin: "Nigeria",
				Type:        Business,
				Title:       longString(101),
				Description: "Valid",
				Address:     "Valid",
				ContactEmail: "test@test.com",
			},
			wantErr: true,
		},
		{
			name: "Description too long (>2000)",
			listing: Listing{
				ID:          "2",
				OwnerOrigin: "Nigeria",
				Type:        Business,
				Title:       "Valid",
				Description: longString(2001),
				Address:     "Valid",
				ContactEmail: "test@test.com",
			},
			wantErr: true,
		},
		{
			name: "Company too long (>100)",
			listing: Listing{
				ID:           "3",
				OwnerOrigin:  "Nigeria",
				Type:         Job,
				Title:        "Valid",
				Description:  "Valid",
				Company:      longString(101),
				Skills:       "Go",
				PayRange:     "100k",
				JobStartDate: time.Now().Add(24 * time.Hour),
				JobApplyURL:  "http://test.com",
				City:         "Lagos",
				ContactEmail: "test@test.com",
			},
			wantErr: true,
		},
		{
			name: "Address too long (>200)",
			listing: Listing{
				ID:          "4",
				OwnerOrigin: "Nigeria",
				Type:        Business,
				Title:       "Valid",
				Description: "Valid",
				Address:     longString(201),
				ContactEmail: "test@test.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.listing.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
