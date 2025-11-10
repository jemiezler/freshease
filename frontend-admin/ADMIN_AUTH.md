# Admin Authentication with Password Login

## Overview
The admin panel now supports both OAuth (Google) and password-based authentication. Admin users can be initialized with an admin role and login with email/password.

## Backend Implementation

### New Module: `backend/modules/auth/password`
- **Service** (`service.go`): Handles password-based authentication and admin initialization
- **Controller** (`controller.go`): HTTP handlers for login and init-admin endpoints
- **DTOs** (`dto.go`): Request/response types
- **Module** (`module.go`): Route registration

### Endpoints

#### POST `/api/auth/login`
Authenticate with email and password.

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "data": {
    "accessToken": "jwt_token_here",
    "user": {
      "id": "uuid",
      "email": "admin@example.com",
      "name": "Admin Name",
      "role": "admin"
    }
  },
  "message": "Login successful"
}
```

#### POST `/api/auth/init-admin`
Initialize the first admin user. Can only be called if no admin user exists.

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "password123",
  "name": "Admin Name"
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "admin@example.com",
      "name": "Admin Name",
      "role": "admin"
    }
  },
  "message": "Admin user created successfully"
}
```

**Error Responses:**
- `409 Conflict`: Admin user already exists
- `400 Bad Request`: Validation errors or other issues

### Features
- **Password Hashing**: Uses bcrypt for secure password storage
- **Admin Role Creation**: Automatically creates "admin" role if it doesn't exist
- **Role Assignment**: Assigns admin role to the initialized user
- **JWT Token**: Returns JWT token compatible with existing auth middleware
- **Single Admin Initialization**: Prevents multiple admins from being created via init-admin endpoint

## Frontend Implementation

### New Pages

#### `/login`
- Supports both Google OAuth and email/password login
- Toggle between OAuth and password login forms
- Link to admin initialization page
- Error handling and validation

#### `/init-admin`
- Admin initialization form
- Password confirmation
- Validation (password length, matching passwords)
- Success/error feedback
- Redirects to login after successful initialization

### Updated Files

#### `lib/auth.ts`
- Added `loginWithPassword()` function
- Added `initAdmin()` function
- Updated `AuthResponse` interface to include user data

#### `lib/auth-context.tsx`
- Updated `login()` to be async
- Improved user refresh logic
- Added support for init-admin page (no auth required)

#### `app/layout-content.tsx`
- Added support for init-admin page (no sidebar/topbar)
- Improved route protection

### Authentication Flow

1. **First Time Setup:**
   - Admin visits `/init-admin`
   - Fills in name, email, password
   - Admin user is created with admin role
   - Redirected to `/login`

2. **Login:**
   - Admin visits `/login`
   - Chooses OAuth or email/password
   - For email/password: enters email and password
   - Receives JWT token
   - Token stored in localStorage
   - Redirected to dashboard

3. **Protected Routes:**
   - All routes except `/login` and `/init-admin` require authentication
   - Unauthenticated users are redirected to `/login`
   - Token is validated on each API request

## Security Notes

- Passwords are hashed with bcrypt (cost: 10)
- JWT tokens are signed with HS256
- Tokens include user ID and email
- Password minimum length: 8 characters
- Admin initialization can only be done once (prevents multiple admins via this endpoint)
- Token is stored in localStorage (consider httpOnly cookies for production)

## Usage

### Initial Setup
1. Start the backend server
2. Navigate to `http://localhost:3000/init-admin`
3. Fill in admin details and submit
4. Admin user is created with admin role

### Login
1. Navigate to `http://localhost:3000/login`
2. Choose login method:
   - **Google OAuth**: Click "Continue with Google"
   - **Email/Password**: Click "Continue with Email", enter credentials
3. Upon successful authentication, redirected to dashboard

### Creating Additional Admins
After the first admin is created, additional admin users can be created through:
- The admin panel (users management)
- Direct database operations
- API calls with proper authentication

The `/init-admin` endpoint will return 409 Conflict if an admin already exists.

## Testing

### Test Admin Initialization
```bash
curl -X POST http://localhost:8080/api/auth/init-admin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin1234",
    "name": "Admin User"
  }'
```

### Test Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin1234"
  }'
```

## Environment Variables

No new environment variables required. Uses existing:
- `JWT_SECRET`: For signing JWT tokens
- `JWT_ACCESS_TTL_MIN`: Token expiration time (default: 15 minutes)

## Notes

- The init-admin endpoint should ideally be protected or limited to development environments
- Consider adding rate limiting for login attempts
- Consider adding password strength requirements
- Consider adding email verification for admin users
- The admin role is automatically created if it doesn't exist

