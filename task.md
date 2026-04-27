# Decision Log

- **2026-04-27**: Implementing dynamic website logo fallback for listings lacking an ImageURL. This enhances visual utility without introducing significant latency.
- **2026-04-27**: Using Google's S2 Favicon service (`https://s2.googleusercontent.com/s2/favicons?domain=<hostname>&sz=256`) to fetch higher resolution icons dynamically via the template function.

# Execution Plan

- [x] **Phase 2: TDD & Implementation**
  - [x] Write failing tests in `internal/ui/renderer_funcs_test.go` for `fallbackImageURL`
  - [x] Create `fallbackImageURL` logic in `internal/ui/renderer_funcs.go`
  - [x] Register `fallbackImageURL` in `BuildGlobalFuncMap()` in `internal/ui/renderer.go`
  - [x] Make tests pass (Green phase)
- [x] **Phase 3: UI Updates & Verification**
  - [x] Update `ui/templates/partials/modal_detail.html`
  - [x] Update `ui/templates/partials/featured_card.html`
  - [x] Update `ui/templates/partials/listing_card.html`
  - [x] Increment cache buster in `head_meta.html`
  - [x] Verify using browser subagent/Playwright

