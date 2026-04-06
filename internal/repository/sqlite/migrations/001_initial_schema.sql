-- initial schema (001)
-- This file represents the baseline schema for the migration system.
-- It incorporates all previous raw string migrations into the initial definitions.

CREATE TABLE IF NOT EXISTS listings (
    id TEXT PRIMARY KEY,
    owner_id TEXT,
    title TEXT,
    description TEXT,
    type TEXT,
    owner_origin TEXT,
    is_active BOOLEAN,
    created_at DATETIME,
    image_url TEXT,
    contact_email TEXT,
    contact_phone TEXT,
    contact_whatsapp TEXT,
    website_url TEXT,
    deadline DATETIME,
    skills TEXT,
    job_start_date DATETIME,
    job_apply_url TEXT,
    company TEXT,
    pay_range TEXT,
    address TEXT,
    city TEXT,
    hours_of_operation TEXT DEFAULT '',
    event_start DATETIME,
    event_end DATETIME,
    status TEXT DEFAULT 'Approved',
    featured BOOLEAN DEFAULT 0
);
-- STATEMENT
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    google_id TEXT UNIQUE,
    email TEXT,
    name TEXT,
    avatar_url TEXT,
    created_at DATETIME,
    role TEXT DEFAULT 'User'
);
-- STATEMENT
CREATE TABLE IF NOT EXISTS feedback (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    type TEXT,
    content TEXT,
    created_at DATETIME
);
-- STATEMENT
-- Legacy Migrations (needed for existing databases)
-- These might fail on fresh DBs where the columns are already in CREATE TABLE, 
-- or on already-migrated DBs. We ignore errors for these in 001.
ALTER TABLE listings ADD COLUMN address TEXT;
-- STATEMENT
ALTER TABLE listings ADD COLUMN city TEXT;
-- STATEMENT
ALTER TABLE listings ADD COLUMN hours_of_operation TEXT DEFAULT '';
-- STATEMENT
UPDATE listings SET city = '' WHERE city IS NULL OR city = '';
-- STATEMENT
UPDATE listings SET city = '' WHERE city = 'Unknown';
-- STATEMENT
ALTER TABLE listings ADD COLUMN event_start DATETIME;
-- STATEMENT
ALTER TABLE listings ADD COLUMN event_end DATETIME;
-- STATEMENT
ALTER TABLE listings ADD COLUMN status TEXT DEFAULT 'Approved';
-- STATEMENT
ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'User';
-- STATEMENT
ALTER TABLE listings ADD COLUMN featured BOOLEAN DEFAULT 0;
-- STATEMENT
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_listings_owner_id ON listings(owner_id);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_listings_filter_sort ON listings(is_active, status, type, created_at DESC);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_listings_city ON listings(is_active, status, city);
-- STATEMENT
CREATE VIRTUAL TABLE IF NOT EXISTS listings_fts USING fts5(
    title, description, city,
    content=listings,
    content_rowid=rowid,
    tokenize='trigram'
);
-- STATEMENT
CREATE TRIGGER IF NOT EXISTS listings_ai AFTER INSERT ON listings BEGIN
    INSERT INTO listings_fts(rowid, title, description, city)
    VALUES (new.rowid, new.title, new.description, new.city);
END;
-- STATEMENT
CREATE TRIGGER IF NOT EXISTS listings_ad AFTER DELETE ON listings BEGIN
    INSERT INTO listings_fts(listings_fts, rowid, title, description, city)
    VALUES ('delete', old.rowid, old.title, old.description, old.city);
END;
-- STATEMENT
CREATE TRIGGER IF NOT EXISTS listings_au AFTER UPDATE ON listings BEGIN
    INSERT INTO listings_fts(listings_fts, rowid, title, description, city)
    VALUES ('delete', old.rowid, old.title, old.description, old.city);
    INSERT INTO listings_fts(rowid, title, description, city)
    VALUES (new.rowid, new.title, new.description, new.city);
END;
-- STATEMENT
INSERT INTO listings_fts(listings_fts) VALUES('rebuild');
-- STATEMENT
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    name TEXT,
    claimable BOOLEAN,
    is_system BOOLEAN,
    active BOOLEAN,
    requires_special_validation BOOLEAN,
    created_at DATETIME,
    updated_at DATETIME,
    active_fixed BOOLEAN DEFAULT 0
);
-- STATEMENT
CREATE TABLE IF NOT EXISTS claim_requests (
    id TEXT PRIMARY KEY,
    listing_id TEXT NOT NULL,
    listing_title TEXT,
    user_id TEXT NOT NULL,
    user_name TEXT,
    user_email TEXT,
    status TEXT NOT NULL DEFAULT 'Pending',
    created_at DATETIME
);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_claim_requests_user_listing ON claim_requests(user_id, listing_id);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_claim_requests_status ON claim_requests(status);
-- STATEMENT
ALTER TABLE categories ADD COLUMN active_fixed BOOLEAN DEFAULT 0;
-- STATEMENT
UPDATE categories SET active = 1, active_fixed = 1 WHERE active_fixed = 0;
