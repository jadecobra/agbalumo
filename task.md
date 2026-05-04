# Decision Log

- **Decision**: Remove the "Open Now / Closed" status indicator from listing cards and detail modals.
- **Date**: 2026-05-04
- **Reason**: The signal is unreliable due to missing time zone context and data source fragility. Providing incorrect information violates the "Zero-Cognitive-Load Curation" and "Truth over Feature" principles.
- **Tradeoff**: Loss of "Open Now" feature in favor of system reliability and user trust.
- **Constraint**: Any UI element that depends on server-side time without location-aware precision must be removed until localized time support is implemented.

# Execution Plan

## Phase 2: Autonomous Execution Loop (TDD)
- [x] **TDD (Red)**: Update `listing_card_test.go` and `listing_read_test.go` to assert absence of status badges.
- [x] **Refactor Backend**: Remove `IsCurrentlyOpen` from `domain.Listing` and its calculation in `ListingHandler`.
- [x] **Update UI (Green)**: Remove status badges from `listing_card.html` and `modal_detail.html`.
- [x] **Codify Lesson**: Add `[TRIGGER: TIME_SENSITIVE_SIGNAL]` to `coding-standards.md`.

## Phase 3: Audit & Resilience
- [x] **ADR**: Create `docs/adr/2026-05-04-remove-unreliable-open-status.md`.
- [x] **Verification**: Run `go test ./...` and `go run ./cmd/verify visual-audit`.
- [x] **Final CI**: Run `go run ./cmd/verify precommit` and `git push`.
