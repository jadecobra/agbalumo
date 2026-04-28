# Decision Log

- **2026-04-27**: Stripping non-food noise from Food listing UI. ContactEmail is removed for Food listings in `modal_detail.html`.
- **2026-04-28**: Enhancing Favicon Aesthetics for Listings.
  - **The Problem**: Centered favicons are too small (`w-24`) and blurry when upscaled by the browser.
  - **Upscaling Strategy**: Use `image-rendering: pixelated` in CSS to make low-res icons sharp instead of blurry. Continue using the Google Favicon service (`size=256`) to avoid latency-heavy scraping.
  - **Resizing Strategy**: Replace fixed sizing (`w-20 h-20` / `w-24 h-24`) with percentage-based sizing (`w-4/5 h-4/5`) combined with `object-contain` to fill the bounding box without cutting off.
  - **Complexity Kill-Switch**: Rejecting a dynamic icon scraper. The latency budget allows < 100ms; scraping external sites for high-res icons would violate this.

# Execution Plan

## Phase 2: Autonomous Execution Loop (TDD)
- [ ] Modify `ui/templates/partials/listing_card.html` to apply `image-rendering: pixelated` to the fallback image.
- [ ] Update `ui/templates/partials/listing_card.html` to change image classes from `w-20 h-20 md:w-24 md:h-24` to `w-3/4 h-3/4 max-w-[128px] max-h-[128px]`.
- [ ] Verify changes visually using `browser_subagent` across viewports.

## Phase 3: Audit & Resilience
- [ ] Run `go run ./cmd/verify browser` for regression testing.
