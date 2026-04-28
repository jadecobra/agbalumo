# Decision Log

- **2026-04-27**: Stripping non-food noise from Food listing UI. ContactEmail is removed for Food listings in `modal_detail.html` (Note: not found in `listing_card.html`). Description is removed/truncated for Food listings in `listing_card.html` to focalize `TopDish` and `RegionalSpecialty`.
- **Product Interrogation**: Removing the description entirely for Food listings in the grid card emphasizes the signature dish and specialty. However, keeping a 1-line clamp provides consistent vertical sizing across mixed-type grids. Decision: Completely remove the description block for Food listings in `listing_card.html` to maximize focus on food specs, ensuring zero residual padding.
- **Observability Strategy**: Track click-through rates on "View Menu" for Food vs Job detail views to confirm increased focus on primary food CTAs.

# Execution Plan

## Phase 2: Autonomous Execution Loop (TDD)
- [ ] Implement `ContactEmail` wrapping in `ui/templates/partials/modal_detail.html`
- [ ] Implement `Description` truncation/removal and focalize `TopDish` / `RegionalSpecialty` in `ui/templates/partials/listing_card.html`
- [ ] TDD: Create or modify browser test suites (Playwright/headless) via `.agents/skills/browser-verify/SKILL.md`
- [ ] Assert Email field is ABSENT from the DOM for Food listings
- [ ] Assert Email field is PRESENT for Job listings

## Phase 3: Audit & Resilience
- [ ] Run `go run ./cmd/verify precommit`
- [ ] Run `go run ./cmd/verify browser` for automated viewport regression tests
- [ ] Confirm no orphaned borders or empty spacing remains around modified sections
