# 🔧 Google OAuth Callback Debug Guide

## 🚨 Issue: "User canceled login" - No Callback Received

The error `PlatformException(CANCELED, User canceled login, null, null)` indicates that the OAuth flow is not completing properly. Here's how to debug and fix it:

## 🔍 Root Cause Analysis

The issue is likely one of these:

1. **Missing Environment Variables** - Backend OAuth not configured
2. **Wrong Redirect URI** - Google OAuth redirect mismatch  
3. **Backend Not Running** - OAuth endpoints not available
4. **Deep Link Configuration** - Android manifest issue

## 🛠️ Step-by-Step Debugging

### Step 1: Check Backend OAuth Configuration

The backend needs these environment variables:

```bash
# Required OAuth Environment Variables
OAUTH_BASE_URL=http://localhost:8080
JWT_SECRET=your-super-secret-jwt-key-here
JWT_ACCESS_TTL_MIN=15

# Google OAuth (get from Google Cloud Console)
OIDC_GOOGLE_ISSUER=https://accounts.google.com
OIDC_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
OIDC_GOOGLE_CLIENT_SECRET=your-client-secret
OIDC_GOOGLE_REDIRECT_PATH=/api/auth/google/callback
```

### Step 2: Verify Backend is Running

```bash
# Check if backend is running
curl http://localhost:8080/api/auth/google/start

# Should redirect to Google OAuth, not return error
```

### Step 3: Check Google Cloud Console Settings

In Google Cloud Console, ensure:

1. **Authorized redirect URIs** includes:
   - `http://localhost:8080/api/auth/google/callback` (for development)
   - `https://your-domain.com/api/auth/google/callback` (for production)

2. **OAuth consent screen** is configured
3. **Client ID and Secret** are correct

### Step 4: Test OAuth Flow Manually

1. Open browser to: `http://localhost:8080/api/auth/google/start`
2. Should redirect to Google login
3. Complete login
4. Should redirect to: `freshease://callback?code=...&state=...`

### Step 5: Check Flutter Deep Link

The Android manifest should have:
```xml
<intent-filter android:autoVerify="true">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="freshease" android:host="callback" />
</intent-filter>
```

## 🎯 Quick Fixes

### Fix 1: Set Environment Variables

Create `.env` file in backend directory:
```bash
cd /home/jemiezler/Desktop/freshease/backend
cat > .env << EOF
OAUTH_BASE_URL=http://localhost:8080
JWT_SECRET=your-super-secret-jwt-key-here
JWT_ACCESS_TTL_MIN=15
OIDC_GOOGLE_ISSUER=https://accounts.google.com
OIDC_GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
OIDC_GOOGLE_CLIENT_SECRET=your-client-secret
OIDC_GOOGLE_REDIRECT_PATH=/api/auth/google/callback
DATABASE_URL=postgres://postgres:user1234@localhost:5432/trail-teller_db?sslmode=disable
HTTP_PORT=:8080
ENT_DEBUG=false
EOF
```

### Fix 2: Start Backend

```bash
cd /home/jemiezler/Desktop/freshease/backend
go run main.go
```

### Fix 3: Test OAuth Endpoint

```bash
# Test OAuth start endpoint
curl -v http://localhost:8080/api/auth/google/start

# Should return 302 redirect to Google
```

## 🔍 Debug Logs

Add these to your Flutter app to see what's happening:

```dart
// In AuthApi.signInWithGoogle()
print('🔐 Starting Google Sign-in process...');
print('🌐 Opening OAuth URL: $startUrl');
print('📱 Using callback scheme: $_scheme');

// After FlutterWebAuth2.authenticate()
print('✅ OAuth callback received: $callbackUrl');
```

## 🎯 Expected Success Flow

1. **User taps "Continue with Google"**
2. **Flutter opens**: `http://localhost:8080/api/auth/google/start`
3. **Backend redirects to**: Google OAuth login
4. **User completes login**
5. **Google redirects to**: `http://localhost:8080/api/auth/google/callback`
6. **Backend redirects to**: `freshease://callback?code=...&state=...`
7. **Flutter receives callback**
8. **Flutter exchanges code for token**
9. **Success!** User is logged in

## 🚨 Common Issues

- **"missing env: OAUTH_BASE_URL"** → Set environment variables
- **"User canceled login"** → OAuth flow not completing
- **"Page not found"** → Backend not running or wrong URL
- **"Invalid redirect URI"** → Google Console configuration mismatch

## ✅ Success Indicators

When working correctly, you'll see:
```
🔐 Starting Google Sign-in process...
🌐 Opening OAuth URL: http://localhost:8080/api/auth/google/start
📱 Using callback scheme: freshease
✅ OAuth callback received: freshease://callback?code=...
✅ OAuth code and state extracted successfully
🔄 Exchanging code for tokens...
✅ Access token received successfully
💾 Access token stored in SharedPreferences
👤 Fetching user profile...
✅ User profile fetched successfully
```

The key is ensuring the backend OAuth is properly configured and running! 🎯
