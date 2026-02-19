# Roadmap - Pending Features & Debt

## Status Verification
I have audited the codebase and identified the status of all uncompleted TODO items.

| Feature | Status | Notes |
| :--- | :--- | :--- |
| **Test Suite** | ✅ Good | Coverage is 82.2% (Target >80%). Weak spots: `auth.go` (41.7%), `sqlite.go` (64.7%). |
| **Server Restart** | ✅ Done | `scripts/verify_restart.sh` works well. |
| **Code Review** | ✅ Done | `CONTRIBUTING.md` and PR templates created. |
| **Bulk Upload** | ✅ Done | Admin feature to upload CSVs implemented and verified. |
| **Admin Dashboard** | ✅ Done | Lists users, pending items/counts, and growth charts. |
| **Claim Ownership** | ✅ Done | MVP implemented. Users can claim unowned listings. |
| **UI/Brand** | ✅ Stable | Using Tailwind CLI. Design system defined. |
| **Auth** | ⚠️ Partial | Google Auth works. Apple/Facebook missing. |
| **About Page** | ✅ Done | Static "About Us" page with carousel implemented. |
| **Admin Config** | ⚠️ Partial | Users list added. Settings still missing. |

## Phase 2 Plan (Current)

### 1. Foundation & Stability (Completed)
- [x] **Fix UI/Brand**: Move from Tailwind CDN to a build step (Tailwind CLI) for production stability and offline dev.
- [x] **Refactor Image Upload**: Extracted image handling to dedicated `ImageService` for modularity and testability.
- [x] **Test Coverage**: Boosted coverage to 82.2% (Passed Audit). Weak spots remain in `auth.go`.
- [x] **Scrub History**: Removed tracked test artifacts from git history.
- [x] **Prevent Future Leaks**: Added Secret Scanner to pre-commit hook.

### 2. Core Features (High Value)
- [x] **About Page**: Create a static "About Us" page with the requested carousel.
- [x] **Claim Ownership**: Critical for business adoption. Allow users to claim "seeded" listings.
- [x] **Bulk Upload**: Admin feature to upload CSVs.

### 3. Advanced Features (Pending)
- [ ] **Auth Expansion**: Add Apple/Facebook (requires developer accounts).
- [x] **Enhanced Admin**: Users view added. Charts and better moderation tools implemented.
- [ ] **Admin Customization**: Allow admin to change colors/fonts from dashboard (requested in TODO).

### 4. Audit Findings & Weak Spots
- [ ] **Auth Coverage**: Increase coverage for `internal/handler/auth.go` (specifically `findOrCreateUser`).
- [ ] **Security Monitoring**: Monitor `CVE-2025-60876` (busybox) in `alpine:latest`. Update base image when fixed.
- [ ] **CI Integration**: Add `govulncheck` and `docker scout` to CI pipeline for automated security scans.

## Immediate Next Step
I recommend we focus on **Foundation & Stability**:
1.  **Code Review**: Create `CONTRIBUTING.md` and PR templates to standardize contributions.
2.  **Auth Expansion**: Begin research/implementation for Apple/Facebook login if developer accounts are available.

**Which would you like to tackle next?**
