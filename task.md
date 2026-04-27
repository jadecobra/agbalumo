# Decision Log

- **2026-04-27**: Initialized task to investigate and fix unpopulated `menuURL` fields in production listings. Need to verify current scraper behavior and update it to look for variations like `/menu`, `/order`, or `/order-online`.
- **2026-04-27**: Confirmed via browser audit of agbalumo.com that listings like Native Restaurant & Lounge, AGEGE BUKKAH, and Joloff lack menu URLs.
- **2026-04-27**: Identified that `internal/service/scraper.go` only checks `href` attributes for keywords. Propose updating it to inspect anchor text (e.g., `<a href="/foo">Order Online</a>`) to capture missed menus without introducing excessive DOM tree parsing latency.
- **2026-04-27**: Proactively rejected adding `goquery` or similar full DOM tree parsers to maintain minimum memory and latency footprint in the background processing job.
- **2026-04-27**: Plan to fix menu URL enrichment pipeline. Adding `enrichment_attempted_at` column to track attempts and avoid hourly retries of failed sites. Updating User-Agent to avoid WAF blocks.

# Execution Plan

- [ ] **Phase 1: Database Migration**
  - [ ] Create `internal/repository/sqlite/migrations/005_enrichment_tracking.sql`
- [ ] **Phase 2: Domain & Repository Updates (TDD)**
  - [ ] Add `EnrichmentAttemptedAt` to `internal/domain/listing.go`
  - [ ] Update `internal/repository/sqlite/queries.go`
  - [ ] Update `internal/repository/sqlite/sqlite_listing_write.go`
  - [ ] Update `internal/repository/sqlite/sqlite_listing_read.go`
- [ ] **Phase 3: Scraper Job Update (TDD)**
  - [ ] Update `enrichSingle` in `internal/service/scraper_job.go`
- [ ] **Phase 4: FindEnrichmentTargets Update (TDD)**
  - [ ] Update WHERE clause in `internal/repository/sqlite/sqlite_listing_read.go`
- [ ] **Phase 5: User-Agent Update (TDD)**
  - [ ] Update User-Agent in `internal/service/scraper.go`

