# Decision Log

- **2026-04-27**: Implementing dynamic website logo fallback for listings lacking an ImageURL. This enhances visual utility without introducing significant latency.
- **2026-04-27**: Using Google's S2 Favicon service (`https://s2.googleusercontent.com/s2/favicons?domain=<hostname>&sz=256`) to fetch higher resolution icons dynamically via the template function.

# Execution Plan

- [ ] **Phase 2: TDD & Implementation**
  - [ ] Write failing tests in `internal/ui/renderer_funcs_test.go` for `fallbackImageURL`
  - [ ] Create `fallbackImageURL` logic in `internal/ui/renderer_funcs.go`
  - [ ] Register `fallbackImageURL` in `BuildGlobalFuncMap()` in `internal/ui/renderer.go`
  - [ ] Make tests pass (Green phase)
- [ ] **Phase 3: UI Updates & Verification**
  - [ ] Update `ui/templates/partials/modal_detail.html`
  - [ ] Update `ui/templates/partials/featured_card.html`
  - [ ] Update `ui/templates/partials/listing_card.html`
  - [ ] Increment cache buster in `head_meta.html`
  - [ ] Verify using browser subagent/Playwright
