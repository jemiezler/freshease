# Admin Authentication Setup

## Overview
The admin panel uses Google OAuth for authentication. The authentication flow is handled through the backend API.

## Frontend Implementation

### Files Created
1. **`lib/auth.ts`** - Authentication utilities (token management, API calls)
2. **`lib/auth-context.tsx`** - React context for authentication state
3. **`app/login/page.tsx`** - Login page with Google OAuth button
4. **`app/auth/callback/page.tsx`** - OAuth callback handler
5. **Updated `app/layout.tsx`** - Added AuthProvider
6. **Updated `app/layout-content.tsx`** - Added authentication checks
7. **Updated `components/ui/Topbar.tsx`** - Added user info and logout button

### How It Works
1. User clicks "Continue with Google" on the login page
2. User is redirected to backend `/api/auth/google/start`
3. Backend redirects to Google OAuth
4. Google redirects back to backend callback
5. Backend should redirect to frontend callback page with code and state
6. Frontend exchanges code for JWT token
7. Token is stored in localStorage as `admin_token`
8. User is redirected to dashboard

## Backend Configuration Required

### Current Issue
The backend callback (`/api/auth/{provider}/callback`) currently redirects to a mobile app deep link:
```
freshease://callback?code=...&state=...
```

### Required Change
The backend needs to be updated to support web redirects. The callback should redirect to:
```
http://localhost:3000/auth/callback?code=...&state=...&provider=google
```
(or the production admin URL)

### Backend Modification Needed
In `backend/modules/auth/authoidc/controller.go`, the `Callback` function needs to:
1. Check if the request is from a web client (e.g., check User-Agent or a query parameter)
2. Redirect to the web callback URL instead of the mobile deep link

Example modification:
```go
// Check if web client (you can add a query param like ?web=true when starting OAuth)
isWeb := c.Query("web") == "true" || c.Get("User-Agent") contains "Mozilla"
if isWeb {
    adminURL := os.Getenv("ADMIN_CALLBACK_URL") // e.g., http://localhost:3000/auth/callback
    redirectURL := fmt.Sprintf("%s?code=%s&state=%s&provider=%s", adminURL, code, state, provider)
    return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
}
// Mobile app deep link (existing)
redirectURL := fmt.Sprintf("freshease://callback?code=%s&state=%s", code, state)
return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
```

### Alternative: Environment-Based Redirect
Alternatively, you can use an environment variable to determine the redirect URL:
```go
redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
if redirectURL == "" {
    redirectURL = "freshease://callback" // Default to mobile
}
redirectURL = fmt.Sprintf("%s?code=%s&state=%s&provider=%s", redirectURL, code, state, provider)
```

## Environment Variables

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
```

### Backend (.env)
```env
ADMIN_CALLBACK_URL=http://localhost:3000/auth/callback
# or for production
ADMIN_CALLBACK_URL=https://admin.yourdomain.com/auth/callback
```

## Testing

1. Start the backend server
2. Start the admin frontend: `npm run dev` or `yarn dev`
3. Navigate to `http://localhost:3000/login`
4. Click "Continue with Google"
5. Complete Google OAuth flow
6. Should redirect back to admin panel and be authenticated

## Security Notes

- Tokens are stored in localStorage (consider httpOnly cookies for production)
- Token is automatically included in API requests via `ApiClient`
- Unauthenticated users are automatically redirected to login
- Token is validated on each page load via `/api/whoami` endpoint

