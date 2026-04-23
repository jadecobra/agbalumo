# Service Constraints

- Pure business logic — no HTTP imports, no DB imports
- CSV operations are split: `csv.go` (parsing), `csv_export.go` (generation), `csv_duplicate.go` (dedup)
- All business rules that validate domain objects belong in `internal/domain/`, not here
- Services accept interfaces (e.g., `domain.ListingStore`) — never concrete repository types
