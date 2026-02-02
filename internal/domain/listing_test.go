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
