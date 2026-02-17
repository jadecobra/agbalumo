package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestMockListingRepository(t *testing.T) {
	m := &mock.MockListingRepository{}
	ctx := context.Background()

	// Test Save behavior
	m.On("Save", ctx, testifyMock.Anything).Return(nil).Once()
	if err := m.Save(ctx, domain.Listing{}); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test Save Error
	m.On("Save", ctx, testifyMock.Anything).Return(errors.New("save error")).Once()
	if err := m.Save(ctx, domain.Listing{}); err == nil {
		t.Error("Expected error from Save")
	}

	// Test FindByID
	m.On("FindByID", ctx, "found").Return(domain.Listing{ID: "found"}, nil)
	if l, _ := m.FindByID(ctx, "found"); l.ID != "found" {
		t.Error("Expected ID to be 'found'")
	}

	// Assert
	m.AssertExpectations(t)
}
