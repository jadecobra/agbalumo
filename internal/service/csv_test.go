package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestParseAndImport(t *testing.T) {
	svc := NewCSVService()
	ctx := context.Background()

	t.Run("Valid Import", func(t *testing.T) {
		csvContent := `title,type,description,origin,email,website
Test Biz,Business,Desc 1,Ghana,test@test.com,example.com
Test Svc,Service,Desc 2,Nigeria,svc@test.com,
`
		repo := &mock.MockListingRepository{}
		repo.On("FindByTitle", ctx, testifyMock.Anything).Return([]domain.Listing{}, nil)
		repo.On("Save", ctx, testifyMock.Anything).Return(nil).Times(2)

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
		repo.AssertExpectations(t)
	})

	t.Run("Missing Headers", func(t *testing.T) {
		csvContent := `title,description
Test,Desc`
		repo := &mock.MockListingRepository{}
		_, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err == nil {
			t.Error("Expected error for missing headers, got nil")
		}
		if !strings.Contains(err.Error(), "missing required header") {
			t.Errorf("Expected missing header error, got %v", err)
		}
	})

	t.Run("Partial Failure", func(t *testing.T) {
		csvContent := `title,type,description,origin,email,phone
Good,Business,Desc,Ghana,a@b.com,
BadTypeDefaultsToBusiness,InvalidType,Desc,Ghana,a@b.com,
MissingDesc,Business,,Ghana,a@b.com,
,Business,Desc,Ghana,a@b.com,
MissingOrigin,Business,Desc,,a@b.com,
MissingEmailHasPhone,Business,Desc,Ghana,,12345
MissingAllContact,Business,Desc,Ghana,,`

		repo := &mock.MockListingRepository{}
		repo.On("FindByTitle", ctx, testifyMock.Anything).Return([]domain.Listing{}, nil)
		// Expect 4 successful saves (Good, BadTypeDefaultsToBusiness, MissingOrigin, MissingEmailHasPhone)
		repo.On("Save", ctx, testifyMock.Anything).Return(nil).Times(4)

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.TotalProcessed != 7 {
			t.Errorf("Expected 7 total, got %d", result.TotalProcessed)
		}
		if result.SuccessCount != 4 {
			t.Errorf("Expected 4 success, got %d", result.SuccessCount)
		}
		if result.FailureCount != 3 {
			t.Errorf("Expected 3 failures, got %d", result.FailureCount)
		}
		repo.AssertExpectations(t)
	})

	t.Run("Duplicate Check", func(t *testing.T) {
		csvContent := `title,type,description,origin,email,phone,address,city
Dup Listing,Business,Same Desc,Ghana,dup@test.com,1234,123 St,Accra
`
		repo := &mock.MockListingRepository{}
		repo.On("FindByTitle", ctx, "Dup Listing").Return([]domain.Listing{
			{Title: "Dup Listing", Type: domain.Business, Description: "Same Desc", ContactEmail: "dup@test.com", OwnerOrigin: "Nigeria"},
		}, nil)
		// No Save expected!

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
		repo.AssertExpectations(t)
	})

	t.Run("Duplicate DB Error", func(t *testing.T) {
		csvContent := `title,type,description,email
DB Error Biz,Business,Desc,test@test.com
`
		repo := &mock.MockListingRepository{}
		repo.On("FindByTitle", ctx, "DB Error Biz").Return([]domain.Listing{}, errors.New("db connection failed"))

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.FailureCount != 1 {
			t.Errorf("Expected 1 failure, got %d", result.FailureCount)
		}
		if !strings.Contains(result.Errors[0], "Database error checking duplicates") {
			t.Errorf("Expected DB error string, got %s", result.Errors[0])
		}
		repo.AssertExpectations(t)
	})

	t.Run("Category Parsing Edge Cases", func(t *testing.T) {
		csvContent := `title,type,description,email
Lowercase Food,food,Desc,a@b.com
Random Type,Random,Desc,a@b.com
`
		repo := &mock.MockListingRepository{}
		repo.On("FindByTitle", ctx, testifyMock.Anything).Return([]domain.Listing{}, nil)
		repo.On("Save", ctx, testifyMock.Anything).Return(nil).Times(2)

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result.FailureCount != 0 {
			t.Errorf("Expected 0 failures, got %d", result.FailureCount)
		}
		repo.AssertExpectations(t)
	})
}
