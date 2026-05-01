CREATE TABLE IF NOT EXISTS saved_listings (
    user_id TEXT NOT NULL,
    listing_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, listing_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE CASCADE
);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_saved_listings_user ON saved_listings(user_id);
