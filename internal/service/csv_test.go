package service

import (
	"context"
	"strings"
	"testing"

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
		csvContent := `title,type,description,origin,email
Good,Business,Desc,Ghana,a@b.com
Bad,InvalidType,Desc,Ghana,a@b.com
MissingDesc,Business,,Ghana,a@b.com`

		repo := &mock.MockListingRepository{}
		// Expect only 1 successful save
		repo.On("Save", ctx, testifyMock.Anything).Return(nil).Once()

		result, err := svc.ParseAndImport(ctx, strings.NewReader(csvContent), repo)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.TotalProcessed != 3 {
			t.Errorf("Expected 3 total, got %d", result.TotalProcessed)
		}
		if result.SuccessCount != 1 {
			t.Errorf("Expected 1 success (first row), got %d", result.SuccessCount)
		}
		if result.FailureCount != 2 {
			t.Errorf("Expected 2 failures, got %d", result.FailureCount)
		}
		repo.AssertExpectations(t)
	})
}
