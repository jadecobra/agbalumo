# Task: Decommission 'Signature Dish' Feature

The "Signature Dish" feature is providing unreliable data in production and needs to be removed to maintain data integrity and user trust.

## Phase 1: Planning & ADR
- [x] Identify all occurrences of `TopDish` and "Signature Dish"
- [x] Create ADR-20260505 documenting the decommissioning
- [x] Get approval for the decommissioning plan

## Phase 2: Implementation

### Domain & Logic
- [x] Remove `TopDish` from `internal/domain/listing.go`
- [x] Remove `FieldTopDish` from `internal/domain/constants.go`
- [x] Remove `TopDish` from `AdaSignals` in `internal/service/scraper.go` and disable scraping logic
- [x] Remove `TopDish` mapping in `internal/service/scraper_job.go`
- [x] Remove `TopDish` scan/write in `internal/repository/sqlite/`

### UI & Templates
- [x] Remove "Signature Dish" rendering from `ui/templates/partials/listing_card.html`
- [x] Remove "Signature Dish" rendering from `ui/templates/partials/modal_detail.html`
- [x] Remove `top_dish` input from `ui/templates/components/listing_form_common_fields.html`

### Cleanup
- [x] Update `internal/module/listing/listing_form.go`
- [x] Update all affected tests (scraper, repository, renderer)
- [x] Update documentation in `docs/cli/verify.md`

## Phase 3: Verification
- [x] Run `go run ./cmd/verify preflight`
- [x] Run all Go tests `go test ./...`
- [x] Run `go run ./cmd/verify precommit`
- [x] Perform browser verification of listings and forms
