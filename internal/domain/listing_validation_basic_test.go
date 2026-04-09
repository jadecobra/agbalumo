package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOriginValidation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantErr error
		name    string
		origin  string
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
			t.Parallel()
			l := Listing{
				ID:           "3",
				Type:         Product,
				Title:        "Shea Butter",
				ContactEmail: "shea@example.com",
				City:         "Dakar",
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

func TestContactRequirement(t *testing.T) {
	t.Parallel()
	listing := Listing{
		ID:          "2",
		OwnerOrigin: "Ghana",
		Type:        Business,
		Title:       "Jollof Place",
		City:        "Accra",
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

	// With Phone only
	listing.ContactWhatsApp = ""
	listing.ContactPhone = "+2348000000000"
	err = listing.Validate()
	assert.NoError(t, err)

	// With Website only
	listing.ContactPhone = ""
	listing.WebsiteURL = "https://example.com"
	err = listing.Validate()
	assert.NoError(t, err)
}

func TestAddressValidation(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			l := Listing{
				ID:           "test-addr",
				OwnerOrigin:  "Ghana",
				Type:         tt.lType,
				Title:        "Test Biz",
				ContactEmail: "test@example.com",
				Address:      tt.address,
				City:         "Kumasi",
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

func TestCityRequirement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		lType   Category
		city    string
		address string
		wantErr bool
	}{
		{
			name:    "Business requires city",
			lType:   Business,
			city:    "",
			address: "123 St",
			wantErr: true,
		},
		{
			name:    "Food requires city",
			lType:   Food,
			city:    "",
			address: "123 St",
			wantErr: true,
		},
		{
			name:    "Event requires city",
			lType:   Event,
			city:    "",
			address: "123 St",
			wantErr: true,
		},
		{
			name:    "Valid Business with city",
			lType:   Business,
			city:    "Lagos",
			address: "123 St",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := Listing{
				ID:           "test-city",
				OwnerOrigin:  "Nigeria",
				Type:         tt.lType,
				Title:        "Test",
				ContactEmail: "test@example.com",
				Address:      tt.address,
				City:         tt.city,
				CreatedAt:    time.Now(),
				IsActive:     true,
				// Additional fields for Event
				EventStart: time.Now().Add(24 * time.Hour),
				EventEnd:   time.Now().Add(25 * time.Hour),
			}
			err := l.Validate()
			if tt.wantErr {
				assert.Error(t, err, "Expected error for type %s with empty city", tt.lType)
			} else {
				assert.NoError(t, err, "Expected no error for type %s with city", tt.lType)
			}
		})
	}
}
