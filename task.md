# Decision Log

- **2026-04-27**: Integrating Google Places API to backfill rating and review metrics for food listings.
- **2026-04-27**: **Architecture Tradeoff**: Leveraging the existing `BackgroundService` to automatically execute a 30-day polling cycle. Listings with `rating_updated_at` either empty or older than 30 days will be continuously processed in small batches (e.g., 5 per tick) to preserve strict API quota limits.

# Execution Plan

## Phase 2: TDD & Implementation
- [ ] Create Google Places API Client in `internal/service/google_places.go` with strict field masking
- [ ] Create `RatingEnricherJob` in `internal/service/rating_enricher.go` mimicking `ScraperJob`
- [ ] Wire `RatingEnricherJob` into `internal/service/background.go` to process eligible listings per tick
- [ ] Write integration test validating API payload binding and parsing

## Phase 3: Audit & Resilience
- [ ] Audit DB state mapping against 30-day thresholds
