package service

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/testutil"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestBackgroundService_ExpireListings(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)
	// Seed an expired listing
	expiredListing := domain.Listing{
		ID:          "exp1",
		Title:       "Expired",
		Status:      domain.ListingStatusApproved,
		IsActive:    true,
		Deadline:    time.Now().Add(-24 * time.Hour),
		OwnerOrigin: "Nigeria",
		Type:        domain.Request,
	}
	_ = repo.Save(context.Background(), expiredListing)

	service := NewBackgroundService(repo, nil)

	// Since expireListings is private but we are in package service, we can call it.
	service.expireListings(context.Background())

	// Verify it's now inactive
	l, err := repo.FindByID(context.Background(), "exp1")
	if err != nil {
		t.Fatalf("Failed to find listing: %v", err)
	}
	if l.IsActive {
		t.Errorf("Expected status to be inactive (IsActive=false)")
	}
}

func TestBackgroundService_ExpireListings_Error(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)
	service := NewBackgroundService(repo, nil)

	// Passing a canceled context to simulate a database query error or timeout
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	service.expireListings(ctx)
	// We expect this to print an error log and return early, achieving coverage of err != nil branch.
}

func TestBackgroundService_StartTicker_Cancels(t *testing.T) {
	t.Parallel()
	repo := testutil.SetupTestRepository(t)

	service := NewBackgroundService(repo, nil)
	ctx, cancel := context.WithCancel(context.Background())

	// Run Ticker in goroutine
	done := make(chan bool)
	go func() {
		service.StartTicker(ctx)
		done <- true
	}()

	// Allow it to run for a tiny bit to trigger the initial "run once"
	time.Sleep(10 * time.Millisecond)

	// Cancel context to stop ticker
	cancel()

	// Wait for return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("StartTicker did not exit after context cancellation")
	}
}
