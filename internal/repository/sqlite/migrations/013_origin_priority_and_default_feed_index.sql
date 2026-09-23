ALTER TABLE listings ADD COLUMN origin_priority INTEGER GENERATED ALWAYS AS (
    CASE WHEN (regional_specialty LIKE '%Nigerian%' OR owner_origin LIKE '%Nigerian%') THEN 0 ELSE 1 END
) VIRTUAL;
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_listings_default_feed ON listings(is_active, status, featured DESC, origin_priority ASC, heat_level DESC, rating DESC, created_at DESC, id ASC);
