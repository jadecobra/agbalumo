-- Rebuild FTS5 to include address and backfill state/country (006)

-- 1. Rebuild FTS5 with address
DROP TABLE IF EXISTS listings_fts;
-- STATEMENT
CREATE VIRTUAL TABLE IF NOT EXISTS listings_fts USING fts5(
    title, description, city, state, country, address,
    content=listings,
    content_rowid=rowid,
    tokenize='trigram'
);
-- STATEMENT
-- 2. Update Triggers to include address
DROP TRIGGER IF EXISTS listings_ai;
-- STATEMENT
CREATE TRIGGER listings_ai AFTER INSERT ON listings BEGIN
    INSERT INTO listings_fts(rowid, title, description, city, state, country, address)
    VALUES (new.rowid, new.title, new.description, new.city, new.state, new.country, new.address);
END;
-- STATEMENT
DROP TRIGGER IF EXISTS listings_ad;
-- STATEMENT
CREATE TRIGGER listings_ad AFTER DELETE ON listings BEGIN
    INSERT INTO listings_fts(listings_fts, rowid, title, description, city, state, country, address)
    VALUES ('delete', old.rowid, old.title, old.description, old.city, old.state, old.country, old.address);
END;
-- STATEMENT
DROP TRIGGER IF EXISTS listings_au;
-- STATEMENT
CREATE TRIGGER listings_au AFTER UPDATE ON listings BEGIN
    INSERT INTO listings_fts(listings_fts, rowid, title, description, city, state, country, address)
    VALUES ('delete', old.rowid, old.title, old.description, old.city, old.state, old.country, old.address);
    INSERT INTO listings_fts(rowid, title, description, city, state, country, address)
    VALUES (new.rowid, new.title, new.description, new.city, new.state, new.country, new.address);
END;
-- STATEMENT
INSERT INTO listings_fts(listings_fts) VALUES('rebuild');
