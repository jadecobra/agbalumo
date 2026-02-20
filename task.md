# Task Tracking & Agent Assignment

## Legend
- **SDET**: Software Development Engineer in Test (Tests & Verification)
- **BE**: Backend Engineer (Go, DB, Logic)
- **FE**: Frontend/UI Engineer (HTML, Tailwind, HTMX)
- **ARCH**: Lead Architect (Planning, Standards)

## Phase 2: Foundation & Core Features

| ID | Status | Agent | Task | Drift Check |
| :--- | :--- | :--- | :--- | :--- |
| **50** | `[x]` | **ARCH** | **Plan Phase 2 Strategy**<br>Refine architecture for Test Coverage and Bulk Upload. | `implementation_plan.md` updated? |
| **51** | `[x]` | **SDET** | **Boost Test Coverage (Target: 80%)**<br>Focus: `ui/renderer.go` and `cmd/serve.go`. | Improved to ~74% |
| **52** | `[x]` | **BE** | **Bulk Upload Feature (CSV)**<br>Implement efficient CSV parsing and seeding. | Handles 100+ items? |
| **53** | `[x]` | **FE** | **Bulk Upload UI**<br>Create Admin upload form with progress feedback. | Success/Error states visible? |
| **54** | `[x]` | **FE** | **Enhance Feedback UI**<br>Add submission visual feedback (currently silent). | Toast/Modal appears? |
| **55** | `[x]` | **BE** | **Admin Users View**<br>Implement list view for all users. | Table visible? |
| **56** | `[x]` | **ARCH** | **Service Listing Logic**<br>Ensure "Service" type is featured on homepage. | "Service" always visible? |
| **57** | `[x]` | **ARCH** | **Address TODO Questions**<br>Review and answer outstanding questions in TODO file. | Answers recorded? |
| **59** | `[x]` | **SDET** | **Boost Test Coverage (Target: 80%)**<br>Focus: `cmd` and `internal/handler`. | 81.6% |
| **60** | `[x]` | **BE** | **Refactor Root Command**<br>Make `Execute` testable. | Testable? |
| **61** | `[x]` | **BE** | **Refactor Seed Command**<br>Extract config logic. | Testable? |
| **62** | `[x]` | **SEC** | **Scrub History**<br>Remove tracked test artifacts. | Clean history? |
| **63** | `[x]` | **SEC** | **Prevent Future Leaks**<br>Add Secret Scanner to pre-commit hook. | Catches secrets? |
- [x] **[FE]** Polish Claim Ownership UI (HTMX) <!-- id: 46 -->
- [x] **[BE]** Implement Claim Ownership Logic <!-- id: 43 -->
- [x] Restore Test Coverage (Target: 87.8%) <!-- id: 5 -->
- [x] **[SDET]** Verify Claim Flow & Security <!-- id: 45 -->
| **64** | `[x]` | **ARCH** | **Code Review Standards**<br>Create `CONTRIBUTING.md` and PR templates. | Standards defined? |
| **65** | `[x]` | **FE/BE** | **Enhanced Admin UI**<br>Add charts and better moderation tools. | Charts visible? |
