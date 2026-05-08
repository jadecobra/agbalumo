# Backfill Locations Tool Context

A specialized maintenance script for populating missing location metadata.

## Operational Constraints
*   **Dry Run**: ALWAYS support a `--dry-run` flag that logs intended changes without writing to the database.
*   **Idempotency**: The script must be safe to run multiple times. Check for existing data before updating.
*   **Rate Limiting**: When calling external geocoding APIs, implement backoff and rate limiting to avoid quota exhaustion.

## Safety
*   **Backup**: Run a database backup before executing this script in production.
