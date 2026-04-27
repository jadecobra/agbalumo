# Decision Log

- **2026-04-27**: Implementing dynamic website logo fallback for listings lacking an ImageURL. This enhances visual utility without introducing significant latency.
- **2026-04-27**: Using Google's S2 Favicon service (`https://s2.googleusercontent.com/s2/favicons?domain=<hostname>&sz=256`) to fetch higher resolution icons dynamically via the template function.
- **2026-04-27**: Implementing Delivery Platform Badges for Food Listings. Storing normalized platform names (UberEats, DoorDash, Grubhub) as a JSON array of strings in the `Listing` struct and SQLite DB.
- **2026-04-27**: Decision made to use **Passive Badges** (non-clickable) instead of buttons, as we are only storing platform names and not specific deep links/URLs. This avoids dead ends and maintains the < 60 seconds utility without data complexity.

# Execution Plan

- [x] **Phase 2: TDD & Implementation**
  - [x] Add `DeliveryPlatforms` field to `Listing` struct in `internal/domain/listing.go`
  - [x] Update SQLite schema migrations (`internal/repository/sqlite/migrations/008_delivery_platforms.sql`)
  - [x] Update CRUD repository in `internal/repository/sqlite/` to save/load `DeliveryPlatforms`
  - [x] Write repository test for `DeliveryPlatforms`
  - [x] Implement `hasDelivery` helper in `internal/ui/renderer_funcs.go`
  - [x] Register `hasDelivery` in `BuildGlobalFuncMap`
  - [x] Write unit test for `hasDelivery`
- [x] **Phase 3: UI Updates & Verification**
  - [x] Update `ui/templates/partials/modal_detail.html` with "Order Delivery" section (Passive Badges)
  - [x] Verify using browser subagent or Playwright

