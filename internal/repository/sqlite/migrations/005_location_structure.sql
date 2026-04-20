-- Add state and country columns to listings (005)
ALTER TABLE listings ADD COLUMN state TEXT DEFAULT '';
-- STATEMENT
ALTER TABLE listings ADD COLUMN country TEXT DEFAULT 'USA';
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_listings_state_country ON listings(state, country);
-- STATEMENT
-- Rebuild FTS5 table to include state and country
DROP TABLE IF EXISTS listings_fts;
-- STATEMENT
CREATE VIRTUAL TABLE IF NOT EXISTS listings_fts USING fts5(
    title, description, city, state, country,
    content=listings,
    content_rowid=rowid,
    tokenize='trigram'
);
-- STATEMENT
-- Update triggers for FTS5
DROP TRIGGER IF EXISTS listings_ai;
-- STATEMENT
CREATE TRIGGER listings_ai AFTER INSERT ON listings BEGIN
    INSERT INTO listings_fts(rowid, title, description, city, state, country)
    VALUES (new.rowid, new.title, new.description, new.city, new.state, new.country);
END;
-- STATEMENT
DROP TRIGGER IF EXISTS listings_ad;
-- STATEMENT
CREATE TRIGGER listings_ad AFTER DELETE ON listings BEGIN
    INSERT INTO listings_fts(listings_fts, rowid, title, description, city, state, country)
    VALUES ('delete', old.rowid, old.title, old.description, old.city, old.state, old.country);
END;
-- STATEMENT
DROP TRIGGER IF EXISTS listings_au;
-- STATEMENT
CREATE TRIGGER listings_au AFTER UPDATE ON listings BEGIN
    INSERT INTO listings_fts(listings_fts, rowid, title, description, city, state, country)
    VALUES ('delete', old.rowid, old.title, old.description, old.city, old.state, old.country);
    INSERT INTO listings_fts(rowid, title, description, city, state, country)
    VALUES (new.rowid, new.title, new.description, new.city, new.state, new.country);
END;
-- STATEMENT
INSERT INTO listings_fts(listings_fts) VALUES('rebuild');
