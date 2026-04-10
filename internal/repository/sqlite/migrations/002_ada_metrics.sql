-- 002_ada_metrics.sql
-- Migration to add the metrics table for tracking Ada's Time-to-Comfort.
CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    value REAL,
    metadata TEXT, -- JSON blob for additional details
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_metrics_event_type ON metrics(event_type);
-- STATEMENT
CREATE INDEX IF NOT EXISTS idx_metrics_created_at ON metrics(created_at);
