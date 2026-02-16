package mock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/mock"
)

func TestMockListingRepository(t *testing.T) {
	m := &mock.MockListingRepository{}
	ctx := context.Background()

	// Test default behavior (nil functions)
	if err := m.Save(ctx, domain.Listing{}); err != nil {
		t.Errorf("Expected nil error when function is nil, got %v", err)
	}
	if l, err := m.FindByID(ctx, "id"); err != nil || l.ID != "" {
		t.Errorf("Expected empty listing and nil error, got %v, %v", l, err)
	}
	if l, err := m.FindAll(ctx, "", "", false); err != nil || l != nil {
		t.Errorf("Expected nil listings and nil error, got %v, %v", l, err)
	}

	// Test defined behavior
	m.SaveFn = func(ctx context.Context, l domain.Listing) error {
		return errors.New("save error")
	}
	if err := m.Save(ctx, domain.Listing{}); err == nil {
		t.Error("Expected error from SaveFn")
	}

	m.FindByIDFn = func(ctx context.Context, id string) (domain.Listing, error) {
		return domain.Listing{ID: "found"}, nil
	}
	if l, _ := m.FindByID(ctx, "id"); l.ID != "found" {
		t.Error("Expected ID to be 'found'")
	}

	// Cover other methods
	m.SaveUser(ctx, domain.User{})
	m.FindUserByGoogleID(ctx, "gid")
	m.FindUserByID(ctx, "uid")
	m.FindAllByOwner(ctx, "oid")
	m.Delete(ctx, "id")
	m.GetCounts(ctx)
	m.ExpireListings(ctx)
}
