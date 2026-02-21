# Agbalumo: User and Admin Journeys Mapping

Based on the current routing structure in `cmd/server.go`, here is the comprehensive mapping of user and admin journeys throughout the application.

## 1. Unauthenticated Visitor Journey

The entry point for a general user who has not logged in. Their access is limited to read-only views of public data.

*   **Landing & Discovery:**
    *   `GET /` - **Home Page:** Arrives at the main application page. Can browse the general feed of public listings.
    *   `GET /listings/fragment` - **Dynamic Loading:** As the visitor scrolls, or applies filters, listings are dynamically loaded into the feed using HTMX.
*   **Deep Dive:**
    *   `GET /listings/:id` - **Listing Details:** Clicks on a specific listing to view its full details.
*   **Information:**
    *   `GET /about` - **About Page:** Reads information regarding the platform's purpose.
*   **Onboarding:**
    *   `GET /auth/google/login` - **Google OAuth:** Initiates the login process to become an Authenticated User.
    *   `GET /auth/dev` - **Dev Login:** (Development environments only) Bypasses OAuth for local testing.
    *   `GET /auth/google/callback` - **Callback:** Completes the authentication loop and establishes the user session.

---

## 2. Authenticated Standard User Journey

Once logged in, the user gains the ability to create, manage, and interact with platform data, specifically their own listings and profile.

*   **Profile Management:**
    *   `GET /profile` - **User Profile:** Accesses their personal dashboard to view saved items, account details, and listings they manage.
*   **Listing Creation & Management:**
    *   `POST /listings` - **Create:** Submits a new listing to the platform.
    *   `GET /listings/:id/edit` - **Edit View:** Opens the form to modify an existing listing they own.
    *   `PUT /listings/:id` / `POST /listings/:id` - **Update Action:** Submits the modifications to their listing (handling both RESTful PUT and standard form POST fallbacks).
    *   `DELETE /listings/:id` - **Delete Action:** Removes a listing they own from the platform.
    *   `POST /listings/:id/claim` - **Claim Action:** Requests ownership/management rights over a currently unmanaged listing.
*   **Platform Interaction:**
    *   `GET /feedback/modal` - **Feedback UI:** Opens a modal to provide thoughts or report issues.
    *   `POST /feedback` - **Submit Feedback:** Sends the feedback to the system.
*   **Session Termination:**
    *   `GET /auth/logout` - **Logout:** Ends the active session and returns the user to an Unauthenticated Visitor state.

---

## 3. Administrator Journey

Administrators are authenticated users with elevated privileges allowing them to moderate content, manage users, and perform bulk operations. 

*   **Elevation (Admin Login):**
    *   `GET /admin/login` - **Admin Portal Login:** An authenticated user attempts to access the admin area and is presented with a secondary login/verification screen (e.g., verifying an admin password or specific role check).
    *   `POST /admin/login` - **Admin Authentication:** Submits credentials to establish an administrative session.
*   **High-Level Overview:**
    *   `GET /admin` - **Admin Dashboard:** The main landing page for administrators, likely showing system metrics, pending approvals, and quick actions.
*   **User Management:**
    *   `GET /admin/users` - **User Directory:** Views a list of all registered users on the platform.
*   **Content Moderation & Management:**
    *   `GET /admin/listings` - **Global Listing View:** Accesses a table/list of all listings across the system regardless of ownership, often used to filter for "pending" status.
    *   `POST /admin/listings/:id/approve` - **Approve Listing:** Marks a pending listing as approved for public view.
    *   `POST /admin/listings/:id/reject` - **Reject Listing:** Declines a submitted listing, keeping it off the public feed.
    *   `GET /admin/listings/delete-confirm` - **Delete Confirmation:** UI prompt before hard-deleting a listing.
    *   `POST /admin/listings/delete` - **Hard Delete Action:** Permanently removes a listing from the database.
    *   `POST /admin/listings/bulk` - **Bulk Operations:** Applies an action (like approve or reject) to multiple selected listings simultaneously.
*   **Data Ingestion:**
    *   `POST /admin/upload` - **Bulk Upload:** Uploads a CSV data file to programmatically create many listings at once.
