# API Route Verification via UI

This plan outlines the steps to verify that all API routes specified in `docs/api.md` and `docs/openapi.yaml` are accessible and functional through the UI

## Verification Plan

We will use the UI to walk through the user and admin journeys, covering all documented API endpoints.

### 1. Public Routes
- [ ] **Home**: `GET /` - Verify homepage loads and shows featured listings.
- [ ] **About**: `GET /about` - Verify about page loads.
- [ ] **Search/Filter**: `GET /listings/fragment` - Use the search bar and category filters to verify dynamic listing loading.
- [ ] **Listing Detail**: `GET /listings/:id` - Click a listing and verify the detail page/modal loads.
- [ ] **Google OAuth**: `GET /auth/google/login` - Verify OAuth flow initiates.

### 2. User Routes (Authentication Required)
- [ ] **Dev Login**: `GET /auth/dev?email=dev@agbalumo.com` - Verify session is established and redirected to home.
- [ ] **User Profile**: `GET /profile` - Verify profile page shows user info and their listings.
- [ ] **Create Listing**: `POST /listings` - Open the "New Listing" modal, fill the form, and submit. Verify redirect and appearance in feed.
- [ ] **Edit Listing**: `GET /listings/:id/edit` - Open edit form for the new listing.
- [ ] **Update Listing**: `PUT/POST /listings/:id` - Change listing details and submit. Verify changes are saved.
- [ ] **Claim Listing**: `POST /listings/:id/claim` - Find an unclaimed listing and attempt to claim it.
- [ ] **Feedback**: `GET /feedback/modal` and `POST /feedback` - Open feedback modal and submit feedback.
- [ ] **Delete Listing**: `DELETE /listings/:id` - Delete the created listing and verify it's removed.
- [ ] **Logout**: `GET /auth/logout` - Verify session is cleared.

### 3. Admin Routes (Admin Authentication Required)
- [ ] **Admin Login**: `GET /admin/login` - Verify login page loads (should be accessible without auth).
- [ ] **Admin Login Action**: `POST /admin/login` with `code=agbalumo2024` - Verify admin session is established.
- [ ] **Admin Dashboard**: `GET /admin` - Verify statistics and pending listings are visible.
- [ ] **User Management**: `GET /admin/users` - Verify list of all users.
- [ ] **Listing Management**: `GET /admin/listings` - Verify all listings with filters/sorting.
- [ ] **Approve/Reject**: `POST /admin/listings/:id/approve` and `:id/reject` - Moderate a pending listing.
- [ ] **Toggle Featured**: `POST /admin/listings/:id/featured` - Toggle promotion.
- [ ] **Bulk Actions**: `POST /admin/listings/bulk` - Select multiple listings and approve/reject them.
- [ ] **Hard Delete**: `GET /admin/listings/delete-confirm` and `POST /admin/listings/delete`.
- [ ] **Bulk Upload**: `POST /admin/upload` - Upload `sample_upload.csv` and verify listings are created.

### Execution Steps
1. Start the server: `go run main.go serve` (Port 8443 by default).
2. Invoke browser subagent with detailed instructions for each journey.
3. Capture screenshots and record videos for each verification step.

---

## Notes

### Bug Fix Applied
- **Admin login routes** (`GET /admin/login` and `POST /admin/login`) were incorrectly requiring authentication. Fixed to allow unauthenticated access so users can actually reach the login page.
