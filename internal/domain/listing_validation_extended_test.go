package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateDeadline(t *testing.T) {
	now := time.Now()

	tests := []struct {
		deadline time.Time
		wantErr  error
		name     string
		typeStr  Category
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
		{
			name:     "Deadline in Past",
			deadline: now.Add(-48 * time.Hour),
			typeStr:  Request,
			wantErr:  errors.New("deadline cannot be in the past"),
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
				Address:      "123 Valid St",
				City:         "Lagos",
				CreatedAt:    now,
				Deadline:     tt.deadline,
				IsActive:     true,
			}
			err := l.Validate()
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequestRequiresOrigin(t *testing.T) {
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

func TestHoursOfOperationRestriction(t *testing.T) {
	tests := []struct {
		name    string
		lType   Category
		hours   string
		wantErr bool
	}{
		{"Business with hours", Business, "9-5", false},
		{"Service with hours", Service, "9-5", false},
		{"Food with hours", Food, "9-5", false},
		{"Product with hours", Product, "9-5", true},
		{"Event with hours", Event, "9-5", true},
		{"Job with hours", Job, "9-5", true},
		{"Request with hours", Request, "9-5", true},
		{"Product without hours", Product, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Listing{
				ID:               "test-restrict",
				OwnerOrigin:      "Ghana",
				Type:             tt.lType,
				Title:            "Test",
				ContactEmail:     "test@example.com",
				Company:          "Acme",
				Skills:           "Go",
				JobStartDate:     time.Now().Add(24 * time.Hour),
				EventStart:       time.Now().Add(24 * time.Hour),
				EventEnd:         time.Now().Add(25 * time.Hour),
				Address:          "123 St",
				City:             "Lome",
				HoursOfOperation: tt.hours,
				CreatedAt:        time.Now(),
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
				ID:           "1",
				OwnerOrigin:  "Nigeria",
				Type:         Business,
				Title:        longString(101),
				Description:  "Valid",
				Address:      "Valid",
				City:         "Lagos",
				ContactEmail: "test@test.com",
			},
			wantErr: true,
		},
		{
			name: "Description too long (>2000)",
			listing: Listing{
				ID:           "2",
				OwnerOrigin:  "Nigeria",
				Type:         Business,
				Title:        "Valid",
				Description:  longString(2001),
				Address:      "Valid",
				City:         "Lagos",
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
				ID:           "4",
				OwnerOrigin:  "Nigeria",
				Type:         Business,
				Title:        "Valid",
				Description:  "Valid",
				Address:      longString(201),
				City:         "Lagos",
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

func TestValidate_Request_NoCreatedAt(t *testing.T) {
	l := Listing{
		Type:         Request,
		Title:        "Help",
		ContactEmail: "test@test.com",
		City:         "Lagos",
		Deadline:     time.Now().Add(24 * time.Hour),
		OwnerOrigin:  "Nigeria",
	}
	err := l.Validate()
	assert.NoError(t, err)
}

func TestValidate_Food_Success(t *testing.T) {
	l := Listing{
		Type:         Food,
		Title:        "Jollof",
		ContactEmail: "j@test.com",
		City:         "Accra",
		Address:      "Street 1",
		OwnerOrigin:  "Ghana",
		CreatedAt:    time.Now(),
	}
	err := l.Validate()
	assert.NoError(t, err)
}
