package sqlite

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
)

func TestConcurrentWriteContention(t *testing.T) {
	dbPath := "write_contention.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	defer func() { _ = repo.Close() }()

	ctx := context.Background()
	concurrency := 20
	var wg sync.WaitGroup
	errs := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			listing := domain.Listing{
				ID:          fmt.Sprintf("listing-%d", id),
				Title:       fmt.Sprintf("Listing %d", id),
				Description: "Desc",
				Type:        domain.Service,
				OwnerOrigin: "Nigeria",
				City:        "Lagos",
				Status:      domain.ListingStatusApproved,
				CreatedAt:   time.Now(),
				OwnerID:     "user1",
			}
			err := repo.Save(ctx, listing)
			if err != nil {
				errs <- err
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	lockErrorCount := 0
	for err := range errs {
		if err != nil {
			t.Logf("Error during concurrent write: %v", err)
			lockErrorCount++
		}
	}

	if lockErrorCount > 0 {
		t.Logf("Found %d write errors (likely database is locked because currently MaxOpenConns=100 for EVERYTHING)", lockErrorCount)
	} else {
		t.Log("No write errors found. SQLite might be handling them with busy_timeout, but pool isolation is preferred for structural integrity.")
	}
}
