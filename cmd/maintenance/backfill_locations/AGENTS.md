# Backfill Locations Tool Context

A specialized maintenance script for populating missing location metadata.

## Operational Constraints
*   **Heuristic Extraction**: Extracts state codes from address strings using regular expressions.
*   **Idempotency**: The script is safe to run multiple times; it checks and populates default country and missing state fields only.

## Safety
*   **Backup**: Run a database backup before executing this script in production.
