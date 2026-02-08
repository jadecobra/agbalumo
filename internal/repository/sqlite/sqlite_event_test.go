package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestSaveAndFindEvent(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second) // Truncate to match SQLite precision

	event := domain.Listing{
		ID:           "event-1",
		Title:        "Music Festival",
		OwnerOrigin:  "Ghana",
		Type:         domain.Event,
		IsActive:     true,
		CreatedAt:    now,
		ContactEmail: "fest@example.com",
		EventStart:   now.Add(24 * time.Hour),
		EventEnd:     now.Add(48 * time.Hour),
	}

	// 1. Save Event
	if err := repo.Save(ctx, event); err != nil {
		t.Fatalf("Failed to save event: %v", err)
	}

	// 2. Find Event
	found, err := repo.FindByID(ctx, "event-1")
	if err != nil {
		t.Fatalf("Failed to find event: %v", err)
	}

	// 3. Verify Fields
	if !found.EventStart.Equal(event.EventStart) {
		t.Errorf("Expected EventStart %v, got %v", event.EventStart, found.EventStart)
	}
	if !found.EventEnd.Equal(event.EventEnd) {
		t.Errorf("Expected EventEnd %v, got %v", event.EventEnd, found.EventEnd)
	}
}
