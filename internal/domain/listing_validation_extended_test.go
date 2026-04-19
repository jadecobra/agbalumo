package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateDeadline(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
	longString := func(n int) string {
		return strings.Repeat("a", n)
	}

	baseListing := func() Listing {
		return Listing{
			ID:           "test-len",
			OwnerOrigin:  "Nigeria",
			Type:         Business,
			Title:        "Valid Title",
			Description:  "Valid Description",
			Address:      "Valid Address",
			City:         "Lagos",
			ContactEmail: "test@example.com",
			CreatedAt:    time.Now(),
		}
	}

	tests := []struct {
		mutate  func(*Listing)
		name    string
		wantErr bool
	}{
		{
			name:    "Title too long (>100)",
			mutate:  func(l *Listing) { l.Title = longString(101) },
			wantErr: true,
		},
		{
			name:    "Description too long (>2000)",
			mutate:  func(l *Listing) { l.Description = longString(2001) },
			wantErr: true,
		},
		{
			name: "Company too long (>100)",
			mutate: func(l *Listing) {
				l.Type = Job
				l.Company = longString(101)
				l.Skills = "Go"
				l.PayRange = "100k"
				l.JobApplyURL = "http://test.com"
				l.JobStartDate = time.Now().Add(24 * time.Hour)
			},
			wantErr: true,
		},
		{
			name:    "Address too long (>200)",
			mutate:  func(l *Listing) { l.Address = longString(201) },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := baseListing()
			tt.mutate(&l)
			err := l.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidate_Request_NoCreatedAt(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
