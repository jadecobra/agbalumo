# Authentication Flow

This document summarizes the overall authentication and authorization flow in Agbalumo.

## 1. User Authentication (Google OAuth 2.0)
Agbalumo uses Google OAuth exclusively for user registration and login.
There are no traditional username/password combinations.

### Flow
1. **Initiation**: User clicks "Login with Google", which hits `/auth/google/login`.
2. **Redirect**: The server uses `golang.org/x/oauth2` to generate a Google Login URL and redirects the user.
3. **Google Consent**: The user grants access (profile and email scopes) on Google's domain.
4. **Callback**: Google redirects to `/auth/google/callback` with an auth `code`.
5. **Token Exchange**: The server exchanges the `code` for an access token directly with Google.
6. **User Fetch**: The server uses the token to hit `https://www.googleapis.com/oauth2/v2/userinfo` and retrieves the user's `GoogleID`, `Email`, `Name`, and `Picture`.
7. **Find or Create**: The server checks the SQLite database for a user with the given `GoogleID`.
   - If missing, a new `domain.User` record is created.
   - If found, the avatar and name are updated if they changed.
8. **Session**: A local secure cookie session is created (using `gorilla/sessions`), storing the internal `User.ID`.
9. **Final Redirect**: User is redirected to `/`.

### Development Login Check
There is a `/auth/dev-login?email=xxx` route for local development that simulates Google's behavior. This is strictly disabled in production.

## 2. Session Middleware (`AuthMiddleware` & `OptionalAuthMiddleware`)
- Every HTTP request passes through middleware that attempts to read the `user_id` from the secure session cookie.
- **OptionalAuthMiddleware**: If the user exists in the DB, they are attached to `c.Get("User")`. If not, no error is thrown (used for public pages).
- **AuthMiddleware**: Inherits from optional auth. If `c.Get("User")` is nil, redirects to the Google Login prompt.

## 3. Admin Authorization
Agbalumo uses a simple "Access Code" promotion system for Admin rights, instead of a separate admin table.

1. **Prerequisite**: The user must be logged in via Google (Standard User).
2. **Access Page**: The user visits `/admin/login` and enters the secret Admin Code.
3. **Promotion**: If the code matches the server's `ADMIN_CODE` environment variable, the user's `Role` in the database is updated from `UserRoleUser` to `UserRoleAdmin`.
4. **Admin Routes**: Routes under `/admin` are protected by `AdminMiddleware`, which checks that `user.Role == domain.UserRoleAdmin`. If not, they are redirected back to `/admin/login` or the main feed.
