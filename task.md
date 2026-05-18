# Task: Detail Card Grid Restructuring

## Decision Log

- **Decision 1**: Restructure the detail modal scrollable content from a single vertical list into a responsive, premium two-column CSS grid on desktop viewports (`md:grid-cols-2`) and a clean single-column layout on mobile viewports.
- **Decision 2**: Design the two columns with clear content segmentation:
  - **Left Column**: Main info (Location/Map link, Specialties/Verified Origin flag, Job/Event specific details, and the "About" Description block).
  - **Right Column**: Interactive CTAs, hours, and actions (View Menu, Order Delivery platforms, and Contact Info grid).
- **Decision 3**: Transform the Contact Section from a vertical link list into a highly premium **2x2 Action Button Grid** on desktop (WhatsApp, Call, Email, Website) using high-contrast brutalist outlines, flat shadow offsets, and vibrant hover micro-animations to improve aesthetic density and click rates.
- **Decision 4**: Expand the modal container's maximum width from the default `md:max-w-sm` (384px) to `md:max-w-4xl` (896px) on desktop to give the new grid layout sufficient breathing room and premium editorial feel.
- **Decision 5**: Ensure WCAG 2 AA compliance for all new action grid elements, verifying contrast ratios against background surfaces (specifically `earth-cream`, `white`, and `earth-dark`).

---

## Execution Plan

### Phase 1: Planning and Approvals
- [ ] Present clarifying questions and obtain user approval on the grid layout and modal width adjustments.
- [ ] Get sign-off on `task.md`.

### Phase 2: Autonomous Execution Loop (TDD)
- [ ] Run preflight checks and execute existing Playwright E2E UI tests locally (`go run ./cmd/verify browser`) to capture the baseline state.
- [ ] Restructure `ui/templates/partials/modal_detail.html`:
  - Update `modal_base` invocation to pass `MaxWidthClass="md:max-w-4xl"` for desktop viewports.
  - Group scrollable content inside a responsive grid: `<div class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8">`.
  - Group address, origin flags, specialties, job/event cards, and description in the left column.
  - Group CTA buttons, delivery block, hours of operation, and the new 2x2 contact action grid in the right column.
- [ ] Modernize the contact section into a 2x2 grid (`grid grid-cols-2 gap-3`) with uniform touch target heights (minimum 44px) and matching brutalist button styles.
- [ ] Verify accessibility, mobile responsiveness, and layout integrity across target viewports.

### Phase 3: Audit & Resilience
- [ ] Run `go run ./cmd/verify design` to enforce brand compliance and check design tokens.
- [ ] Run `go run ./cmd/verify browser` to execute the full E2E UI verification suite and ensure accessibility (Axe) and surface-parity tests pass.
- [ ] Generate Linux visual regression snapshots if Playwright baselines are impacted.
- [ ] Run local CI audit checks (`go run ./cmd/verify precommit`).
- [ ] Commit all changes atomically with strict Conventional Commits format (`feat(ui): restructure detail modal content into 2-column desktop grid and 2x2 contact grid`).
- [ ] Run `./scripts/pushw.sh` or `gh run watch` to ensure the remote CI build passes.
