# Decision Log

- **2026-04-27**: Implementing dynamic website logo fallback for listings lacking an ImageURL. This enhances visual utility without introducing significant latency.
- **2026-04-27**: Using Google's S2 Favicon service (`https://s2.googleusercontent.com/s2/favicons?domain=<hostname>&sz=256`) to fetch higher resolution icons dynamically via the template function.
- **2026-04-27**: Implementing Delivery Platform Badges for Food Listings. Storing normalized platform names (UberEats, DoorDash, Grubhub) as a JSON array of strings in the `Listing` struct and SQLite DB.
- **2026-04-27**: Decision made to use **Passive Badges** (non-clickable) instead of buttons, as we are only storing platform names and not specific deep links/URLs. This avoids dead ends and maintains the < 60 seconds utility without data complexity.
- **2026-04-27**: Implementing real-time "Open Now" computation for listings.
- **2026-04-27**: Adding transient `IsCurrentlyOpen bool` field to `Listing` struct in `internal/domain/listing.go`.
- **2026-04-27**: Creating lightweight Regex-based hours parser `ComputeIsOpen` in `internal/service/hours_parser.go` with table-driven tests covering at least 5 unstructured formats.
- **2026-04-27**: Updating UI listing card with pulsing green dot and "Open Now" or gray "Closed" badge.
- **2026-04-27**: Adding Quality Proxy (Rating & Review Count) to the Listing Domain to help users assess quality quickly. Proposing to use `text-yellow-400` (Tailwind gold) for the star icon.

# Execution Plan

## Phase 2: TDD & Implementation
- [ ] Add `Rating float64` and `ReviewCount int` to `Listing` struct in `internal/domain/listing.go`
- [ ] Create SQLite migration for `rating` and `review_count` columns in `internal/repository/sqlite/migrations/`
- [ ] Update `sqlite_listing_crud.go`, `sqlite_listing_read.go`, and `queries.go` to support new columns
- [ ] Update default sorting logic in listing retrieval queries (higher rating first for same heat/distance)
- [ ] Write CRUD tests in `internal/repository/sqlite/` verifying persistence and sorting

## Phase 3: UI & Verification
- [ ] Update `ui/templates/partials/listing_card.html` with Star Icon and Gold Text for Rating (if ReviewCount > 0)
- [ ] Verify UI using `go run ./cmd/verify browser` or browser subagent
