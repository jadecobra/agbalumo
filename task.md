# Task: Detail Modal Element Restructuring

## Decision Log

- **Decision 1**: Relocate `.Listing.Type` (category) and `rating_stars` (ratings) out of the header image box into a dedicated, clean horizontal metadata row.
- **Decision 2**: Keep the metadata row **static during scrolling** (fixed immediately below the image header and above the scrollable content). This will be achieved by inserting a new flex container in the modal's main column layout sibling to the image box and the scrollable `.overflow-y-auto` container.
- **Decision 3**: Keep the close button (`btn_close`) overlaid absolute in the top-right corner of the header image box.
- **Decision 4**: Keep the save button (`save_button`) as an absolute overlay inside the image box, but line it up side-by-side with the close button. Ensure they are the exact same size (`w-11 h-11`) and use matching dark translucent backgrounds (`bg-black/40 backdrop-blur-sm rounded-full`) for visual balance and premium contrast.
- **Decision 5**: Center the title (`.Listing.Title`) horizontally at the bottom of the image box (centered text overlay on top of the bottom gradient scrim).

---

## Execution Plan

### Phase 1: Planning and Approvals
- [x] Define execution strategy and gain user approval on task.md and clarifying questions.

### Phase 2: Autonomous Execution Loop (TDD)
- [ ] Run standard Playwright E2E UI tests locally (`go run ./cmd/verify browser`) to capture the baseline state (RED).
- [ ] Restructure `ui/templates/partials/modal_detail.html` (GREEN):
  - Group close and save buttons inside a top-right absolute-positioned flex row container.
  - Apply matching classes to close and save buttons to ensure identical sizing and styling.
  - Center the title text overlay horizontally at the bottom of the image box.
  - Extract category and ratings out of the image box.
  - Insert a new pinned metadata row below the image box, styled with responsive spacing and clear borders.
- [ ] Verify accessibility, mobile responsiveness, and layout integrity across target viewports.

### Phase 3: Audit & Resilience
- [ ] Run `go run ./cmd/verify design` to enforce brand compliance and check design tokens.
- [ ] Run `go run ./cmd/verify browser` to execute the full E2E UI verification suite and ensure accessibility (Axe) and surface-parity tests pass.
- [ ] Generate linux visual regression snapshots if Playwright baselines are impacted.
- [ ] Run local CI audit checks (`go run ./cmd/verify precommit`).
- [ ] Commit all changes atomically with strict Conventional Commits format (`style(ui): restructure detail modal elements and center title`).
- [ ] Run `./scripts/pushw.sh` or `gh run watch` to ensure the remote CI build passes.
