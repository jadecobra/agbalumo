package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindListingWithNullEventDates(t *testing.T) {
	// 1. Setup in-memory DB
	repo, err := NewSQLiteRepository(":memory:")
	assert.NoError(t, err)

	// 2. Insert a listing with explicit NULL for event_start/event_end
	// We use direct SQL execution to simulate existing data or external modifications
	ctx := context.Background()
	query := `
	INSERT INTO listings (
		id, owner_id, title, description, type, owner_origin, 
		city, address, is_active, created_at, image_url, 
		contact_email, contact_phone, contact_whatsapp, website_url, 
		deadline, event_start, event_end
	) VALUES (
		'old-listing-123', 'owner-1', 'Old Listing', 'Desc', 'Business', 'Ghana',
		'Accra', 'Some Address', true, ?, 'img.jpg',
		'email@example.com', '', '', '',
		?, NULL, NULL
	)`
	
	_, err = repo.db.ExecContext(ctx, query, time.Now(), time.Now().Add(24*time.Hour))
	assert.NoError(t, err)

	// 3. Try to FindByID
	l, err := repo.FindByID(ctx, "old-listing-123")
	
	// Expectation: This should now succeed.
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, "old-listing-123", l.ID)
	// EventStart should be zero time because it was NULL
	assert.True(t, l.EventStart.IsZero())
	assert.True(t, l.EventEnd.IsZero())

    
    // We will assert NoError after we apply the fix.
    // For reproduction step, we just want to run this.
}
