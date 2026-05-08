# User Module Context

Manages user profiles, saved listings, and personal preferences.

## Privacy & Security
*   **Email Protection**: Never expose a user's email address in public UI or JSON responses.
*   **Ownership**: Ensure users can only modify their own profiles and saved lists.

## Features
*   **Saved Listings**: Logic for bookmarking restaurants. Ensure the "Saved" state is reflected accurately in the `listing` module's UI components.
*   **Profile Consistency**: Maintain synchronization between OAuth provider data (Google Name/Picture) and the internal user record.
