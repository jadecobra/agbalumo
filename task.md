# Decision Log

- **2026-04-28**: Implementing Hours-to-Pulse LLM Pipeline for accurate closed indicators.
  - **The Problem**: Regex parsing of operating hours causes false closures for weekend schedules.
  - **The Decision**: Add `structured_hours` field to listings, backed by Gemini extraction during enrichment.
  - **Complexity Kill-Switch**: Pre-computing JSON hours at enrichment time protects the search latency budget (< 100ms impact).

# Execution Plan

## Phase 2: Autonomous Execution Loop (TDD)
- [x] Create database migration `internal/repository/sqlite/migrations/010_structured_hours.sql`.
- [x] Update `internal/domain/listing.go` to include `StructuredHours string \`json:"structured_hours" form:"structured_hours"\``.
- [x] Update `internal/repository/sqlite/queries.go`, `sqlite_listing_read.go`, and `sqlite_listing_write.go` to support `structured_hours`.
- [x] Create a Gemini API extraction client in `internal/service/gemini_extractor.go` to map raw text to JSON.
- [x] Update `scraper_job.go` to populate `structured_hours` via Gemini upon website parsing.
- [x] Update `ComputeIsOpen` in `hours_parser.go` to prefer the `structured_hours` schema.

## Phase 3: Audit & Resilience
- [x] Run `go test ./internal/service/ -run TestComputeIsOpen` to verify logic.
- [x] Run `go run ./cmd/verify precommit` for full regression testing.

