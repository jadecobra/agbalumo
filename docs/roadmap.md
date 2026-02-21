# Roadmap - Pending Features & Debt

## Status Verification
Audited Feb 19, 2026. Cross-referenced against `TODO`, `spec.md`, `task.md`, and CI pipeline.

| Feature | Status | Notes |
| :--- | :--- | :--- |
| **Test Suite** | ✅ Good | Coverage threshold: 87.5%. Auth coverage boosted to >80%. |
| **Server Restart** | ✅ Done | `scripts/verify_restart.sh` works well. |
| **Code Review** | ✅ Done | `CONTRIBUTING.md` and PR templates created. |
| **Bulk Upload** | ✅ Done | Admin feature to upload CSVs implemented and verified. |
| **Admin Dashboard** | ✅ Done | Lists users, pending items/counts, and growth charts. |
| **Claim Ownership** | ✅ Done | MVP implemented. Users can claim unowned listings. |
| **UI/Brand** | ✅ Stable | Using Tailwind CLI. Design system defined. |
| **Auth** | ⚠️ Partial | Google Auth works. Apple/Facebook missing. Coverage >80%. |
| **About Page** | ✅ Done | Static "About Us" page with carousel implemented. |
| **Admin Config** | ⚠️ Partial | Users list added. Settings still missing. |
| **CI/CD** | ✅ Done | `govulncheck`, `docker scout`, deploy to Fly.io all integrated. |
| **Security** | ✅ Done | Secret scanner in pre-commit, history scrubbed, ImageService extracted. |
| **Architecture** | ✅ Done | DB indexes, pagination, and auth middleware decoupling implemented. |

## Phase 2 — Completed

### 1. Foundation & Stability
- [x] **Fix UI/Brand**: Move from Tailwind CDN to a build step (Tailwind CLI) for production stability and offline dev.
- [x] **Refactor Image Upload**: Extracted image handling to dedicated `ImageService` for modularity and testability.
- [x] **Test Coverage**: Boosted to 84.9% threshold. Weak spots remain in `auth.go`.
- [x] **Scrub History**: Removed tracked test artifacts from git history.
- [x] **Prevent Future Leaks**: Added Secret Scanner to pre-commit hook.
- [x] **Code Review Standards**: Created `CONTRIBUTING.md` and PR templates.

### 2. Core Features (High Value)
- [x] **About Page**: Create a static "About Us" page with the requested carousel.
- [x] **Claim Ownership**: Critical for business adoption. Allow users to claim "seeded" listings.
- [x] **Bulk Upload**: Admin feature to upload CSVs.
- [x] **Feedback Button**: Users can report issues and suggest improvements.
- [x] **Hours of Operation**: Added for Business/Food listings.
- [x] **CLI**: Command-line interface for application management.
- [x] **User Journeys**: Mapped journeys for each user type.

### 3. Advanced Features & Architecture
- [x] **Enhanced Admin**: Users view added. Charts and better moderation tools implemented.
- [x] **CI Integration**: `govulncheck` and `docker scout` added to CI pipeline.
- [x] **Database Optimization**: Added DB indexes for filters and pagination limit/offset for listings to prevent table scans and memory bloat.
- [x] **Middleware Refactor**: Decoupled `AuthMiddleware` from `AuthHandler` to improve modularity.
- [x] **Modularity**: Extracted date parsing and improved validation in domain package and simplified template rendering.

## Phase 3 — Pending

### 1. Auth & Security
- [ ] **Auth Expansion**: Add Apple/Facebook login (requires developer accounts).
- [ ] **DevLogin Hardening**: Tighten `DevLogin` to simulate a generic user, forcing the claim flow even in dev (see `spec.md`).

### 2. Admin & Platform
- [ ] **Admin Customization**: Allow admin to change colors/fonts from dashboard (requested in TODO).
- [ ] **Expiration Ticker**: Background service to deactivate listings past `Deadline` or `EventEnd` (specified in `spec.md` §4.2, never built).

### 3. Testing & Quality
- [ ] **Browser Test Expansion**: Expand browser subagent tests to cover the "Create Listing" flow (see `spec.md`).
- [ ] **Security Monitoring**: Monitor `CVE-2025-60876` (busybox) in `alpine:latest`. Update base image when fixed.

## Immediate Next Steps
1. **Admin Customization** — Allow admin to configure colors/fonts from the dashboard.
2. **Fat Handler Refactoring** — Use structural tags and a validation library to replace manual parsing inside `ListingHandler.bindAndMapListing`.
3. **Auth Expansion** — Begin Apple/Facebook login if developer accounts are available.
