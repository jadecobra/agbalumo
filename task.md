# Decision Log

- **2026-04-27**: Implementing dynamic website logo fallback for listings lacking an ImageURL. This enhances visual utility without introducing significant latency.
- **2026-04-27**: Using Google's S2 Favicon service (`https://s2.googleusercontent.com/s2/favicons?domain=<hostname>&sz=256`) to fetch higher resolution icons dynamically via the template function.
- **2026-04-27**: Implementing Delivery Platform Badges for Food Listings. Storing normalized platform names (UberEats, DoorDash, Grubhub) as a JSON array of strings in the `Listing` struct and SQLite DB.
- **2026-04-27**: Decision made to use **Passive Badges** (non-clickable) instead of buttons, as we are only storing platform names and not specific deep links/URLs. This avoids dead ends and maintains the < 60 seconds utility without data complexity.
- **2026-04-27**: Implementing real-time "Open Now" computation for listings.
- **2026-04-27**: Adding transient `IsCurrentlyOpen bool` field to `Listing` struct in `internal/domain/listing.go`.
- **2026-04-27**: Creating lightweight Regex-based hours parser `ComputeIsOpen` in `internal/service/hours_parser.go` with table-driven tests covering at least 5 unstructured formats.
- **2026-04-27**: Updating UI listing card with pulsing green dot and "Open Now" or gray "Closed" badge.

# Execution Plan

## Phase 2: TDD & Implementation
- [x] Add `IsCurrentlyOpen bool` transient field to `Listing` in `internal/domain/listing.go`
- [x] Implement `ComputeIsOpen(hoursText string, currentTime time.Time) bool` in `internal/service/hours_parser.go`
- [x] Write robust table-driven tests for `ComputeIsOpen` with 5+ formats in `internal/service/hours_parser_test.go`
- [x] Update `internal/module/listing/listing.go` to iterate and compute `IsCurrentlyOpen` using `time.Now()`

## Phase 3: UI & Verification
- [x] Update `ui/templates/partials/listing_card.html` with "Open Now" / "Closed" badges
- [x] Verify UI using `go run ./cmd/verify browser` or browser subagent



