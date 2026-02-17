# Roadmap - Pending Features & Debt

## Status Verification
I have audited the codebase and identified the status of all uncompleted TODO items.

| Feature | Status | Notes |
| :--- | :--- | :--- |
| **Test Suite** | ⚠️ Partial | Coverage is 73%. Key handlers covered, but `main.go` and `ui/renderer` need work. |
| **Server Restart** | ✅ Done | `scripts/verify_restart.sh` works well. |
| **Code Review** | ❌ Missing | No `CONTRIBUTING.md` or PR templates. |
| **Bulk Upload** | ✅ Done | Admin feature to upload CSVs implemented and verified. |
| **Admin Dashboard** | ✅ Improved | Lists users, pending items/counts. Missing charts/analytics. |
| **Claim Ownership** | ✅ Done | MVP implemented. Users can claim unowned listings. |
| **UI/Brand** | ✅ Stable | Using Tailwind CLI. Design system defined. |
| **Auth** | ⚠️ Partial | Google Auth works. Apple/Facebook missing. |
| **About Page** | ✅ Done | Static "About Us" page with carousel implemented. |
| **Admin Config** | ⚠️ Partial | Users list added. Settings still missing. |

## Proposed Phase 2 Plan (Prioritized)

### 1. Foundation & Stability (Recommended First)
- [x] **Fix UI/Brand**: Move from Tailwind CDN to a build step (Tailwind CLI) for production stability and offline dev.
- [ ] **Test Coverage**: Boost coverage to 80% (Current: ~74%). Focus on `renderer` and `auth` edge cases.

### 2. Core Features (High Value)
- [x] **About Page**: Create a static "About Us" page with the requested carousel.
- [x] **Claim Ownership**: Critical for business adoption. Allow users to claim "seeded" listings.

### 3. Advanced Features
- [x] **Bulk Upload**: Admin feature to upload CSVs.
- [ ] **Auth Expansion**: Add Apple/Facebook (requires developer accounts).
- [/] **Enhanced Admin**: Users view added. Charts and better moderation tools pending.

## Immediate Next Step
I recommend we start with **Foundation & Stability** (Moving to local Tailwind) or **Core Features** (About Page).

**Which would you like to tackle next?**
