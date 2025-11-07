# 📋 Freshease Backend - Comprehensive Test Cases

## Test Suite Overview

| Metric | Value |
|--------|-------|
| **Total Test Files** | 40+ |
| **Total Test Cases** | 200+ |
| **Modules Covered** | 13 |
| **Test Status** | ✅ ALL TESTS PASSING |
| **Overall Coverage** | 62.0% |
| **Last Updated** | October 29, 2025 |

---

## 🔐 Authentication Module (`auth/authoidc`) - Detailed Test Scenarios

### Test Scenario 1: OAuth Start Flow - Success Cases

#### TC-AUTH-001: Google OAuth Start - Success
**Test Scenario**: User initiates Google OAuth authentication successfully
**Test Steps**:
1. **Setup**: Initialize mock service with Google provider configured
2. **Action**: Send GET request to `/auth/google/start`
3. **Verify**: 
   - HTTP status code is 307 (Temporary Redirect)
   - Location header contains Google OAuth URL
   - URL includes state parameter
   - URL includes nonce parameter
   - Response sets `oidc_state` cookie (HttpOnly, SameSite=Lax)
   - Response sets `oidc_nonce` cookie (HttpOnly, SameSite=Lax)
   - State is stored in StateStore with 10-minute expiration
4. **Expected Result**: Successful redirect to Google OAuth with proper state management

#### TC-AUTH-002: LINE OAuth Start - Success
**Test Scenario**: User initiates LINE OAuth authentication successfully
**Test Steps**:
1. **Setup**: Initialize mock service with LINE provider configured
2. **Action**: Send GET request to `/auth/line/start`
3. **Verify**:
   - HTTP status code is 307 (Temporary Redirect)
   - Location header contains LINE OAuth URL
   - URL includes state parameter
   - URL includes nonce parameter
   - Response sets `oidc_state` cookie (HttpOnly, SameSite=Lax)
   - Response sets `oidc_nonce` cookie (HttpOnly, SameSite=Lax)
   - State is stored in StateStore with 10-minute expiration
4. **Expected Result**: Successful redirect to LINE OAuth with proper state management

### Test Scenario 2: OAuth Start Flow - Error Cases

#### TC-AUTH-003: Unknown Provider - Error
**Test Scenario**: User attempts OAuth with unsupported provider
**Test Steps**:
1. **Setup**: Initialize mock service with no providers configured
2. **Action**: Send GET request to `/auth/unknown/start`
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message: "unknown provider"
   - No cookies are set
   - No redirect occurs
4. **Expected Result**: Proper error handling for unsupported providers

#### TC-AUTH-004: Empty Provider - Error
**Test Scenario**: User attempts OAuth with empty provider parameter
**Test Steps**:
1. **Setup**: Initialize mock service
2. **Action**: Send GET request to `/auth//start` (empty provider)
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - No cookies are set
   - No redirect occurs
4. **Expected Result**: Proper handling of malformed URLs

#### TC-AUTH-005: Service Error - Error
**Test Scenario**: OAuth service returns an error during URL generation
**Test Steps**:
1. **Setup**: Initialize mock service that returns error for AuthCodeURL
2. **Action**: Send GET request to `/auth/google/start`
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message from service
   - No cookies are set
   - No redirect occurs
4. **Expected Result**: Proper error propagation from service layer

### Test Scenario 3: OAuth Callback Flow - Success Cases

#### TC-AUTH-006: Valid State and Code - Success
**Test Scenario**: OAuth provider returns valid authorization code and state
**Test Steps**:
1. **Setup**: 
   - Initialize mock service with Google provider
   - Store state "test-state" in StateStore
2. **Action**: Send GET request to `/auth/google/callback?state=test-state&code=test-code`
3. **Verify**:
   - HTTP status code is 307 (Temporary Redirect)
   - Location header contains `freshease://callback` scheme
   - URL includes `code=test-code` parameter
   - URL includes `state=test-state` parameter
   - State is validated against StateStore
4. **Expected Result**: Successful redirect to mobile app with authorization code

#### TC-AUTH-007: State from Cookie Fallback - Success
**Test Scenario**: State validation falls back to cookie when StateStore is empty
**Test Steps**:
1. **Setup**:
   - Initialize mock service with Google provider
   - Do NOT store state in StateStore
2. **Action**: 
   - Send GET request to `/auth/google/callback?state=test-state&code=test-code`
   - Include cookie: `oidc_state=test-state`
3. **Verify**:
   - HTTP status code is 307 (Temporary Redirect)
   - Location header contains `freshease://callback` scheme
   - State validation succeeds using cookie fallback
4. **Expected Result**: Successful authentication using cookie-based state validation

### Test Scenario 4: OAuth Callback Flow - Error Cases

#### TC-AUTH-008: Missing State Parameter - Error
**Test Scenario**: OAuth callback missing required state parameter
**Test Steps**:
1. **Setup**: Initialize mock service with Google provider
2. **Action**: Send GET request to `/auth/google/callback?code=test-code` (no state)
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message: "missing code or state"
   - No redirect occurs
4. **Expected Result**: Proper validation of required parameters

#### TC-AUTH-009: Missing Code Parameter - Error
**Test Scenario**: OAuth callback missing required code parameter
**Test Steps**:
1. **Setup**: Initialize mock service with Google provider
2. **Action**: Send GET request to `/auth/google/callback?state=test-state` (no code)
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message: "missing code or state"
   - No redirect occurs
4. **Expected Result**: Proper validation of required parameters

#### TC-AUTH-010: State Mismatch - Error
**Test Scenario**: OAuth callback with state that doesn't match stored state
**Test Steps**:
1. **Setup**:
   - Initialize mock service with Google provider
   - Store state "stored-state" in StateStore
2. **Action**: Send GET request to `/auth/google/callback?state=different-state&code=test-code`
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains JSON error message: "invalid state"
   - No redirect occurs
4. **Expected Result**: Proper CSRF protection through state validation

#### TC-AUTH-011: No State Validation - Error
**Test Scenario**: OAuth callback with no valid state (neither in store nor cookie)
**Test Steps**:
1. **Setup**:
   - Initialize mock service with Google provider
   - Do NOT store state in StateStore
2. **Action**: Send GET request to `/auth/google/callback?state=test-state&code=test-code` (no cookie)
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains JSON error message: "invalid state"
   - No redirect occurs
4. **Expected Result**: Proper security validation preventing unauthorized access

### Test Scenario 5: OAuth Exchange Flow - Success Cases

#### TC-AUTH-012: Valid Exchange Request - Success
**Test Scenario**: Mobile app exchanges authorization code for JWT token
**Test Steps**:
1. **Setup**:
   - Initialize mock service with Google provider
   - Store state "test-state" in StateStore
   - Mock service returns JWT token
2. **Action**: Send POST request to `/auth/google/exchange` with JSON body:
   ```json
   {
     "code": "test-code",
     "state": "test-state"
   }
   ```
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains JSON with access token
   - Response includes success message: "Authentication successful"
   - State is removed from StateStore
   - Cookies are cleared (oidc_state and oidc_nonce set to expire)
4. **Expected Result**: Successful token exchange and cleanup

#### TC-AUTH-013: Exchange with Cookie State Fallback - Success
**Test Scenario**: Exchange succeeds using cookie-based state validation
**Test Steps**:
1. **Setup**:
   - Initialize mock service with Google provider
   - Do NOT store state in StateStore
2. **Action**: 
   - Send POST request to `/auth/google/exchange` with JSON body
   - Include cookie: `oidc_state=test-state`
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains JSON with access token
   - State validation succeeds using cookie fallback
4. **Expected Result**: Successful authentication using cookie-based validation

### Test Scenario 6: OAuth Exchange Flow - Error Cases

#### TC-AUTH-014: Invalid JSON Body - Error
**Test Scenario**: Exchange request with malformed JSON
**Test Steps**:
1. **Setup**: Initialize mock service with Google provider
2. **Action**: Send POST request to `/auth/google/exchange` with invalid JSON body
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message: "Invalid request body"
4. **Expected Result**: Proper JSON validation

#### TC-AUTH-015: Missing Code in Exchange - Error
**Test Scenario**: Exchange request missing code field
**Test Steps**:
1. **Setup**: Initialize mock service with Google provider
2. **Action**: Send POST request to `/auth/google/exchange` with JSON body:
   ```json
   {
     "state": "test-state"
   }
   ```
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message: "Missing code or state"
4. **Expected Result**: Proper validation of required fields

#### TC-AUTH-016: Missing State in Exchange - Error
**Test Scenario**: Exchange request missing state field
**Test Steps**:
1. **Setup**: Initialize mock service with Google provider
2. **Action**: Send POST request to `/auth/google/exchange` with JSON body:
   ```json
   {
     "code": "test-code"
   }
   ```
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON error message: "Missing code or state"
4. **Expected Result**: Proper validation of required fields

#### TC-AUTH-017: Invalid State in Exchange - Error
**Test Scenario**: Exchange request with invalid state
**Test Steps**:
1. **Setup**:
   - Initialize mock service with Google provider
   - Store state "stored-state" in StateStore
2. **Action**: Send POST request to `/auth/google/exchange` with JSON body:
   ```json
   {
     "code": "test-code",
     "state": "invalid-state"
   }
   ```
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains JSON error message: "Invalid state"
4. **Expected Result**: Proper state validation for security

#### TC-AUTH-018: Service Exchange Error - Error
**Test Scenario**: OAuth service returns error during token exchange
**Test Steps**:
1. **Setup**:
   - Initialize mock service that returns error for ExchangeAndLogin
   - Store state "test-state" in StateStore
2. **Action**: Send POST request to `/auth/google/exchange` with valid JSON body
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains JSON error message from service
4. **Expected Result**: Proper error handling from service layer

### Test Scenario 7: StateStore Functionality

#### TC-AUTH-019: State Expiration - Security Test
**Test Scenario**: Verify that expired states are properly rejected
**Test Steps**:
1. **Setup**:
   - Initialize StateStore
   - Store state with immediate expiration
2. **Action**: Attempt to validate expired state
3. **Verify**:
   - State validation returns false
   - Expired state is not accepted
4. **Expected Result**: Proper expiration handling for security

#### TC-AUTH-020: State Cleanup - Performance Test
**Test Scenario**: Verify automatic cleanup of expired states
**Test Steps**:
1. **Setup**:
   - Initialize StateStore
   - Store multiple states with different expiration times
2. **Action**: Wait for cleanup cycle (5 minutes)
3. **Verify**:
   - Expired states are removed from store
   - Active states remain in store
4. **Expected Result**: Efficient memory management

### Test Scenario 8: Integration Tests

#### TC-AUTH-021: Complete OAuth Flow - Integration
**Test Scenario**: End-to-end OAuth authentication flow
**Test Steps**:
1. **Setup**: Initialize complete OAuth system with Google provider
2. **Step 1**: Send GET request to `/auth/google/start`
   - Verify redirect to Google OAuth
   - Verify state and nonce cookies are set
3. **Step 2**: Simulate OAuth provider callback with valid state and code
   - Send GET request to `/auth/google/callback?state=<stored-state>&code=<auth-code>`
   - Verify redirect to mobile app scheme
4. **Step 3**: Mobile app exchanges code for token
   - Send POST request to `/auth/google/exchange` with code and state
   - Verify JWT token is returned
   - Verify cleanup of state and cookies
5. **Expected Result**: Complete successful authentication flow

#### TC-AUTH-022: Multi-Provider Support - Integration
**Test Scenario**: Verify support for multiple OAuth providers
**Test Steps**:
1. **Setup**: Initialize OAuth system with Google and LINE providers
2. **Test Google Flow**:
   - Send GET request to `/auth/google/start`
   - Verify Google-specific OAuth URL
3. **Test LINE Flow**:
   - Send GET request to `/auth/line/start`
   - Verify LINE-specific OAuth URL
4. **Verify**: Both providers work independently
5. **Expected Result**: Multi-provider OAuth support

### Test Scenario 9: Security Tests

#### TC-AUTH-023: CSRF Protection - Security Test
**Test Scenario**: Verify CSRF protection through state validation
**Test Steps**:
1. **Setup**: Initialize OAuth system
2. **Action**: Attempt callback with state from different session
3. **Verify**:
   - Request is rejected
   - No authentication occurs
4. **Expected Result**: Proper CSRF protection

#### TC-AUTH-024: State Replay Attack - Security Test
**Test Scenario**: Verify protection against state replay attacks
**Test Steps**:
1. **Setup**: Initialize OAuth system with stored state
2. **Action**: Attempt to reuse same state multiple times
3. **Verify**:
   - First use succeeds
   - Subsequent uses fail
   - State is properly consumed
4. **Expected Result**: Protection against replay attacks

#### TC-AUTH-025: Cookie Security - Security Test
**Test Scenario**: Verify secure cookie configuration
**Test Steps**:
1. **Setup**: Initialize OAuth system
2. **Action**: Send OAuth start request
3. **Verify**:
   - Cookies are HttpOnly
   - Cookies have SameSite=Lax
   - Cookies have secure path
4. **Expected Result**: Proper cookie security configuration

### Test Scenario 10: Error Handling and Edge Cases

#### TC-AUTH-026: Malformed URLs - Edge Case
**Test Scenario**: Handle malformed OAuth URLs gracefully
**Test Steps**:
1. **Setup**: Initialize OAuth system
2. **Action**: Send requests with various malformed URLs
3. **Verify**: Proper error responses for malformed requests
4. **Expected Result**: Graceful handling of malformed input

#### TC-AUTH-027: Concurrent Requests - Edge Case
**Test Scenario**: Handle concurrent OAuth requests
**Test Steps**:
1. **Setup**: Initialize OAuth system
2. **Action**: Send multiple concurrent OAuth start requests
3. **Verify**: All requests are handled correctly
4. **Expected Result**: Thread-safe OAuth handling

#### TC-AUTH-028: Large State Values - Edge Case
**Test Scenario**: Handle unusually large state values
**Test Steps**:
1. **Setup**: Initialize OAuth system
2. **Action**: Send OAuth request with very large state value
3. **Verify**: System handles large values gracefully
4. **Expected Result**: Robust handling of edge case inputs

---

## 👥 Users Module (`modules/users`) - Test Scenarios

### Test Scenario 11: User Management - CRUD Operations

#### TC-USER-001: List Users - Success
**Test Scenario**: Retrieve list of all users successfully
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test users
2. **Action**: Send GET request to `/api/users`
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains JSON array of users
   - Each user object contains required fields (id, email, name, etc.)
   - Response includes proper pagination metadata
4. **Expected Result**: Successful retrieval of user list

#### TC-USER-002: Get User by ID - Success
**Test Scenario**: Retrieve specific user by valid UUID
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test user
2. **Action**: Send GET request to `/api/users/{valid-uuid}`
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains single user object
   - User object matches requested UUID
   - All user fields are present and valid
4. **Expected Result**: Successful retrieval of specific user

#### TC-USER-003: Create User - Success
**Test Scenario**: Create new user with valid data
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/users` with JSON body:
   ```json
   {
     "email": "test@example.com",
     "name": "Test User",
     "avatar": "https://example.com/avatar.jpg"
   }
   ```
3. **Verify**:
   - HTTP status code is 201 (Created)
   - Response body contains created user object
   - User has generated UUID
   - Timestamps are set correctly
4. **Expected Result**: Successful user creation

#### TC-USER-004: Update User - Success
**Test Scenario**: Update existing user with valid data
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test user
2. **Action**: Send PUT request to `/api/users/{valid-uuid}` with JSON body:
   ```json
   {
     "name": "Updated Name",
     "avatar": "https://example.com/new-avatar.jpg"
   }
   ```
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains updated user object
   - Updated fields reflect new values
   - Timestamps are updated correctly
4. **Expected Result**: Successful user update

#### TC-USER-005: Delete User - Success
**Test Scenario**: Delete existing user by valid UUID
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test user
2. **Action**: Send DELETE request to `/api/users/{valid-uuid}`
3. **Verify**:
   - HTTP status code is 204 (No Content)
   - Response body is empty
   - User is removed from repository
4. **Expected Result**: Successful user deletion

### Test Scenario 12: User Management - Error Cases

#### TC-USER-006: Get User - Invalid UUID
**Test Scenario**: Attempt to retrieve user with invalid UUID format
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send GET request to `/api/users/invalid-uuid`
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper UUID validation

#### TC-USER-007: Get User - Not Found
**Test Scenario**: Attempt to retrieve non-existent user
**Test Steps**:
1. **Setup**: Initialize service with empty mock repository
2. **Action**: Send GET request to `/api/users/{non-existent-uuid}`
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains "user not found" message
4. **Expected Result**: Proper handling of missing resources

#### TC-USER-008: Create User - Duplicate Email
**Test Scenario**: Attempt to create user with existing email
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing user with email
2. **Action**: Send POST request to `/api/users` with duplicate email
3. **Verify**:
   - HTTP status code is 409 (Conflict)
   - Response body contains "email already exists" message
4. **Expected Result**: Proper duplicate email handling

#### TC-USER-009: Update User - Not Found
**Test Scenario**: Attempt to update non-existent user
**Test Steps**:
1. **Setup**: Initialize service with empty mock repository
2. **Action**: Send PUT request to `/api/users/{non-existent-uuid}` with valid data
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains "user not found" message
4. **Expected Result**: Proper handling of missing resources

#### TC-USER-010: Delete User - Not Found
**Test Scenario**: Attempt to delete non-existent user
**Test Steps**:
1. **Setup**: Initialize service with empty mock repository
2. **Action**: Send DELETE request to `/api/users/{non-existent-uuid}`
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains "user not found" message
4. **Expected Result**: Proper handling of missing resources

---

## 🛍️ Products Module (`modules/products`) - Test Scenarios

### Test Scenario 13: Product Management - CRUD Operations

#### TC-PROD-001: List Products - Success
**Test Scenario**: Retrieve list of all products successfully
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test products
2. **Action**: Send GET request to `/api/products`
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains JSON array of products
   - Each product object contains required fields
   - Response includes proper pagination metadata
4. **Expected Result**: Successful retrieval of product list

#### TC-PROD-002: Get Product by ID - Success
**Test Scenario**: Retrieve specific product by valid UUID
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test product
2. **Action**: Send GET request to `/api/products/{valid-uuid}`
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains single product object
   - Product object matches requested UUID
   - All product fields are present and valid
4. **Expected Result**: Successful retrieval of specific product

#### TC-PROD-003: Create Product - Success
**Test Scenario**: Create new product with valid data
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/products` with JSON body:
   ```json
   {
     "name": "Test Product",
     "description": "Test product description",
     "price": 9.99,
     "category_id": "valid-category-uuid",
     "vendor_id": "valid-vendor-uuid"
   }
   ```
3. **Verify**:
   - HTTP status code is 201 (Created)
   - Response body contains created product object
   - Product has generated UUID
   - Timestamps are set correctly
4. **Expected Result**: Successful product creation

#### TC-PROD-004: Update Product - Success
**Test Scenario**: Update existing product with valid data
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test product
2. **Action**: Send PUT request to `/api/products/{valid-uuid}` with JSON body:
   ```json
   {
     "name": "Updated Product Name",
     "price": 19.99
   }
   ```
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains updated product object
   - Updated fields reflect new values
   - Timestamps are updated correctly
4. **Expected Result**: Successful product update

#### TC-PROD-005: Delete Product - Success
**Test Scenario**: Delete existing product by valid UUID
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test product
2. **Action**: Send DELETE request to `/api/products/{valid-uuid}`
3. **Verify**:
   - HTTP status code is 204 (No Content)
   - Response body is empty
   - Product is removed from repository
4. **Expected Result**: Successful product deletion

### Test Scenario 14: Product Management - Error Cases

#### TC-PROD-006: Get Product - Invalid UUID
**Test Scenario**: Attempt to retrieve product with invalid UUID format
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send GET request to `/api/products/invalid-uuid`
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper UUID validation

#### TC-PROD-007: Get Product - Not Found
**Test Scenario**: Attempt to retrieve non-existent product
**Test Steps**:
1. **Setup**: Initialize service with empty mock repository
2. **Action**: Send GET request to `/api/products/{non-existent-uuid}`
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains "product not found" message
4. **Expected Result**: Proper handling of missing resources

#### TC-PROD-008: Create Product - Invalid Category
**Test Scenario**: Attempt to create product with non-existent category
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/products` with invalid category_id
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper foreign key validation

#### TC-PROD-009: Create Product - Invalid Vendor
**Test Scenario**: Attempt to create product with non-existent vendor
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/products` with invalid vendor_id
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper foreign key validation

#### TC-PROD-010: Update Product - Not Found
**Test Scenario**: Attempt to update non-existent product
**Test Steps**:
1. **Setup**: Initialize service with empty mock repository
2. **Action**: Send PUT request to `/api/products/{non-existent-uuid}` with valid data
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains "product not found" message
4. **Expected Result**: Proper handling of missing resources

---

## 🛒 Cart Items Module (`modules/cart_items`) - Test Scenarios

### Test Scenario 15: Cart Management - CRUD Operations

#### TC-CART-001: List Cart Items - Success
**Test Scenario**: Retrieve list of cart items for user
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test cart items
2. **Action**: Send GET request to `/api/cart-items`
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains JSON array of cart items
   - Each cart item contains required fields
   - Response includes proper pagination metadata
4. **Expected Result**: Successful retrieval of cart items

#### TC-CART-002: Add Item to Cart - Success
**Test Scenario**: Add product to user's cart
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/cart-items` with JSON body:
   ```json
   {
     "product_id": "valid-product-uuid",
     "quantity": 2
   }
   ```
3. **Verify**:
   - HTTP status code is 201 (Created)
   - Response body contains created cart item object
   - Cart item has generated UUID
   - Quantity is set correctly
4. **Expected Result**: Successful cart item creation

#### TC-CART-003: Update Cart Item - Success
**Test Scenario**: Update quantity of existing cart item
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test cart item
2. **Action**: Send PUT request to `/api/cart-items/{valid-uuid}` with JSON body:
   ```json
   {
     "quantity": 5
   }
   ```
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains updated cart item object
   - Quantity is updated correctly
4. **Expected Result**: Successful cart item update

#### TC-CART-004: Remove Cart Item - Success
**Test Scenario**: Remove item from user's cart
**Test Steps**:
1. **Setup**: Initialize service with mock repository containing test cart item
2. **Action**: Send DELETE request to `/api/cart-items/{valid-uuid}`
3. **Verify**:
   - HTTP status code is 204 (No Content)
   - Response body is empty
   - Cart item is removed from repository
4. **Expected Result**: Successful cart item removal

### Test Scenario 16: Cart Management - Error Cases

#### TC-CART-005: Add Item - Invalid Product
**Test Scenario**: Attempt to add non-existent product to cart
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/cart-items` with invalid product_id
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper product validation

#### TC-CART-006: Add Item - Invalid Quantity
**Test Scenario**: Attempt to add item with invalid quantity
**Test Steps**:
1. **Setup**: Initialize service with mock repository
2. **Action**: Send POST request to `/api/cart-items` with quantity: 0
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper quantity validation

#### TC-CART-007: Update Item - Not Found
**Test Scenario**: Attempt to update non-existent cart item
**Test Steps**:
1. **Setup**: Initialize service with empty mock repository
2. **Action**: Send PUT request to `/api/cart-items/{non-existent-uuid}` with valid data
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains "cart item not found" message
4. **Expected Result**: Proper handling of missing resources

---

## 🤖 GenAI Module (`modules/genai`) - Test Scenarios

### Test Scenario 17: AI Meal Generation - Success Cases

#### TC-GENAI-001: Generate Weekly Meals - Success
**Test Scenario**: Generate weekly meal plan using AI
**Test Steps**:
1. **Setup**: Initialize service with mock AI client
2. **Action**: Send POST request to `/api/genai/weekly-meals` with JSON body:
   ```json
   {
     "dietary_preferences": ["vegetarian"],
     "allergies": ["nuts"],
     "budget": 100,
     "servings": 4
   }
   ```
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains weekly meal plan
   - Plan includes 7 days of meals
   - Each meal includes ingredients and instructions
4. **Expected Result**: Successful AI meal generation

#### TC-GENAI-002: Generate Daily Meals - Success
**Test Scenario**: Generate daily meal plan using AI
**Test Steps**:
1. **Setup**: Initialize service with mock AI client
2. **Action**: Send POST request to `/api/genai/daily-meals` with JSON body:
   ```json
   {
     "meal_type": "dinner",
     "dietary_preferences": ["keto"],
     "allergies": ["dairy"],
     "budget": 50
   }
   ```
3. **Verify**:
   - HTTP status code is 200 (OK)
   - Response body contains daily meal plan
   - Plan includes single meal with ingredients
4. **Expected Result**: Successful AI meal generation

### Test Scenario 18: AI Meal Generation - Error Cases

#### TC-GENAI-003: AI Service Error - Error
**Test Scenario**: AI service returns error during meal generation
**Test Steps**:
1. **Setup**: Initialize service with mock AI client that returns error
2. **Action**: Send POST request to `/api/genai/weekly-meals` with valid data
3. **Verify**:
   - HTTP status code is 500 (Internal Server Error)
   - Response body contains error message
4. **Expected Result**: Proper error handling from AI service

#### TC-GENAI-004: Invalid Request Data - Error
**Test Scenario**: Send invalid data to meal generation endpoint
**Test Steps**:
1. **Setup**: Initialize service with mock AI client
2. **Action**: Send POST request to `/api/genai/weekly-meals` with invalid JSON
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper request validation

---

## 🔧 Middleware Module (`internal/common/middleware`) - Test Scenarios

### Test Scenario 19: Authentication Middleware

#### TC-MID-001: Valid JWT Token - Success
**Test Scenario**: Request with valid JWT token passes authentication
**Test Steps**:
1. **Setup**: Initialize middleware with valid JWT secret
2. **Action**: Send request with valid JWT token in Authorization header
3. **Verify**:
   - Request passes through middleware
   - User context is set correctly
   - No authentication errors occur
4. **Expected Result**: Successful authentication

#### TC-MID-002: Invalid JWT Token - Error
**Test Scenario**: Request with invalid JWT token fails authentication
**Test Steps**:
1. **Setup**: Initialize middleware with JWT secret
2. **Action**: Send request with invalid JWT token in Authorization header
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains authentication error message
4. **Expected Result**: Proper authentication failure handling

#### TC-MID-003: Missing Authorization Header - Error
**Test Scenario**: Request without Authorization header fails authentication
**Test Steps**:
1. **Setup**: Initialize middleware with JWT secret
2. **Action**: Send request without Authorization header
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains "missing authorization header" message
4. **Expected Result**: Proper handling of missing authentication

#### TC-MID-004: Expired JWT Token - Error
**Test Scenario**: Request with expired JWT token fails authentication
**Test Steps**:
1. **Setup**: Initialize middleware with JWT secret
2. **Action**: Send request with expired JWT token in Authorization header
3. **Verify**:
   - HTTP status code is 401 (Unauthorized)
   - Response body contains "token expired" message
4. **Expected Result**: Proper handling of expired tokens

### Test Scenario 20: Request Validation Middleware

#### TC-MID-005: Valid Request Body - Success
**Test Scenario**: Request with valid JSON body passes validation
**Test Steps**:
1. **Setup**: Initialize validation middleware
2. **Action**: Send request with valid JSON body matching schema
3. **Verify**:
   - Request passes through middleware
   - Body is parsed correctly
   - No validation errors occur
4. **Expected Result**: Successful request validation

#### TC-MID-006: Invalid Request Body - Error
**Test Scenario**: Request with invalid JSON body fails validation
**Test Steps**:
1. **Setup**: Initialize validation middleware
2. **Action**: Send request with invalid JSON body
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains validation error message
4. **Expected Result**: Proper request validation

#### TC-MID-007: Malformed JSON - Error
**Test Scenario**: Request with malformed JSON fails validation
**Test Steps**:
1. **Setup**: Initialize validation middleware
2. **Action**: Send request with malformed JSON body
3. **Verify**:
   - HTTP status code is 400 (Bad Request)
   - Response body contains JSON parsing error message
4. **Expected Result**: Proper JSON validation

---

## ⚙️ Configuration Module (`internal/common/config`) - Test Scenarios

### Test Scenario 21: Configuration Loading

#### TC-CONFIG-001: Load Default Configuration - Success
**Test Scenario**: Load configuration with default values when no environment variables are set
**Test Steps**:
1. **Setup**: Clear all environment variables
2. **Action**: Initialize configuration loader
3. **Verify**:
   - Configuration loads successfully
   - Default values are used for all settings
   - No errors occur during loading
4. **Expected Result**: Successful default configuration loading

#### TC-CONFIG-002: Load Environment Configuration - Success
**Test Scenario**: Load configuration with environment variable values
**Test Steps**:
1. **Setup**: Set environment variables for configuration
2. **Action**: Initialize configuration loader
3. **Verify**:
   - Configuration loads successfully
   - Environment variable values are used
   - All settings are properly configured
4. **Expected Result**: Successful environment-based configuration loading

#### TC-CONFIG-003: Mixed Configuration - Success
**Test Scenario**: Load configuration with mix of environment variables and defaults
**Test Steps**:
1. **Setup**: Set some environment variables, leave others unset
2. **Action**: Initialize configuration loader
3. **Verify**:
   - Configuration loads successfully
   - Set environment variables are used
   - Unset variables use default values
4. **Expected Result**: Successful mixed configuration loading

---

## 🛠️ Helpers Module (`internal/common/helpers`) - Test Scenarios

### Test Scenario 22: Utility Functions

#### TC-HELP-001: Pointer Helper - Success
**Test Scenario**: Convert value to pointer correctly
**Test Steps**:
1. **Setup**: Initialize helper function
2. **Action**: Call PtrIfNotNil with string value
3. **Verify**:
   - Function returns pointer to string
   - Pointer value matches input value
   - No errors occur
4. **Expected Result**: Successful pointer conversion

#### TC-HELP-002: Time Helper - Success
**Test Scenario**: Convert time to ISO string correctly
**Test Steps**:
1. **Setup**: Initialize helper function with time value
2. **Action**: Call TimeToISOString with time value
3. **Verify**:
   - Function returns ISO string format
   - String represents correct time
   - Format is valid ISO 8601
4. **Expected Result**: Successful time conversion

#### TC-HELP-003: Nil Input Handling - Edge Case
**Test Scenario**: Handle nil input values gracefully
**Test Steps**:
1. **Setup**: Initialize helper functions
2. **Action**: Call helper functions with nil input
3. **Verify**:
   - Functions handle nil gracefully
   - No panics occur
   - Appropriate nil values are returned
4. **Expected Result**: Proper nil handling

---

## 🗄️ Database Module (`internal/common/db`) - Test Scenarios

### Test Scenario 23: Database Connection

#### TC-DB-001: Invalid DSN - Error
**Test Scenario**: Attempt to connect with invalid database DSN
**Test Steps**:
1. **Setup**: Prepare invalid database DSN
2. **Action**: Attempt to create database client
3. **Verify**:
   - Connection fails with appropriate error
   - Error message indicates DSN issue
   - No connection is established
4. **Expected Result**: Proper error handling for invalid DSN

#### TC-DB-002: Valid DSN - Success
**Test Scenario**: Connect to database with valid DSN
**Test Steps**:
1. **Setup**: Prepare valid database DSN
2. **Action**: Create database client
3. **Verify**:
   - Connection succeeds
   - Client is properly initialized
   - Database operations can be performed
4. **Expected Result**: Successful database connection

#### TC-DB-003: Connection Timeout - Error
**Test Scenario**: Handle database connection timeout
**Test Steps**:
1. **Setup**: Configure short connection timeout
2. **Action**: Attempt to connect to slow database
3. **Verify**:
   - Connection times out appropriately
   - Error message indicates timeout
   - No connection is established
4. **Expected Result**: Proper timeout handling

---

## 🌐 HTTP Router Module (`internal/common/http`) - Test Scenarios

### Test Scenario 24: Route Registration

#### TC-HTTP-001: Register Routes - Success
**Test Scenario**: Register all module routes successfully
**Test Steps**:
1. **Setup**: Initialize HTTP router with all modules
2. **Action**: Register all module routes
3. **Verify**:
   - All routes are registered without errors
   - Routes are accessible via HTTP
   - Route paths are correct
4. **Expected Result**: Successful route registration

#### TC-HTTP-002: API Group Registration - Success
**Test Scenario**: Register API routes under /api group
**Test Steps**:
1. **Setup**: Initialize HTTP router
2. **Action**: Register routes under /api group
3. **Verify**:
   - Routes are accessible under /api prefix
   - Group middleware is applied
   - Routes function correctly
4. **Expected Result**: Successful API group registration

#### TC-HTTP-003: Route Not Found - Error
**Test Scenario**: Handle requests to non-existent routes
**Test Steps**:
1. **Setup**: Initialize HTTP router with limited routes
2. **Action**: Send request to non-existent route
3. **Verify**:
   - HTTP status code is 404 (Not Found)
   - Response body contains appropriate error message
4. **Expected Result**: Proper 404 handling

---

## 🔗 Integration Tests (`internal/common/testutils`) - Test Scenarios

### Test Scenario 25: End-to-End API Integration

#### TC-INT-001: Complete User Flow - Integration
**Test Scenario**: Test complete user management flow
**Test Steps**:
1. **Setup**: Initialize complete system with database
2. **Step 1**: Create user via POST /api/users
   - Verify user creation
3. **Step 2**: Retrieve user via GET /api/users/{id}
   - Verify user retrieval
4. **Step 3**: Update user via PUT /api/users/{id}
   - Verify user update
5. **Step 4**: Delete user via DELETE /api/users/{id}
   - Verify user deletion
6. **Expected Result**: Complete user management flow works end-to-end

#### TC-INT-002: Complete Product Flow - Integration
**Test Scenario**: Test complete product management flow
**Test Steps**:
1. **Setup**: Initialize complete system with database
2. **Step 1**: Create product via POST /api/products
   - Verify product creation
3. **Step 2**: Retrieve product via GET /api/products/{id}
   - Verify product retrieval
4. **Step 3**: Update product via PUT /api/products/{id}
   - Verify product update
5. **Step 4**: Delete product via DELETE /api/products/{id}
   - Verify product deletion
6. **Expected Result**: Complete product management flow works end-to-end

#### TC-INT-003: Authentication Flow - Integration
**Test Scenario**: Test complete OAuth authentication flow
**Test Steps**:
1. **Setup**: Initialize complete system with OAuth providers
2. **Step 1**: Initiate OAuth via GET /auth/google/start
   - Verify OAuth initiation
3. **Step 2**: Handle callback via GET /auth/google/callback
   - Verify callback handling
4. **Step 3**: Exchange code via POST /auth/google/exchange
   - Verify token exchange
5. **Step 4**: Use token for authenticated requests
   - Verify token validation
6. **Expected Result**: Complete OAuth flow works end-to-end

---

## 📊 Test Execution Summary

### Test Categories

| Category | Count | Percentage |
|----------|-------|------------|
| **Unit Tests** | 150+ | 75% |
| **Integration Tests** | 50+ | 25% |
| **Security Tests** | 20+ | 10% |
| **Performance Tests** | 10+ | 5% |

### Test Status Summary

| Status | Count | Percentage |
|--------|-------|------------|
| **Passing** | 200+ | 100% |
| **Failing** | 0 | 0% |
| **Skipped** | 3 | 1.5% |

### Coverage Summary by Module

| Module | Coverage | Test Cases | Status |
|--------|----------|------------|--------|
| **Config** | 100.0% | 9 | ✅ |
| **Helpers** | 100.0% | 13 | ✅ |
| **HTTP/Server** | 100.0% | 6 | ✅ |
| **Middleware** | 91.7% | 8 | ✅ |
| **Users** | 74.6% | 36 | ✅ |
| **HTTP/Router** | 79.4% | 24 | ✅ |
| **DB** | 61.1% | 14 | ✅ |
| **Products** | 55.7% | 28 | ✅ |
| **Cart Items** | 41.0% | 14 | ✅ |
| **GenAI** | 41.2% | 8 | ✅ |
| **Auth/OIDC** | 26.8% | 28 | ✅ |
| **Integration** | 0.0% | 12 | ✅ |
| **Overall** | **62.0%** | **200+** | ✅ |

---

## 🚀 Test Execution Commands

### Running Tests

```bash
# Run all tests
./run_tests.sh

# Run specific module tests
go test ./modules/auth/authoidc -v
go test ./modules/users -v
go test ./modules/products -v
go test ./internal/common/http -v

# Run with coverage
go test ./... -cover

# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run specific test scenarios
go test ./modules/auth/authoidc -run "TestController_Start" -v
go test ./modules/users -run "TestController_CreateUser" -v
```

### Test Data Setup

```bash
# Set up test environment variables
export TEST_DB_DSN="postgres://test:test@localhost:5432/test_db"
export JWT_SECRET="test-secret-key"
export OAUTH_BASE_URL="http://localhost:3000"

# Run tests with environment
go test ./... -v
```

### Continuous Integration

```bash
# CI/CD pipeline test execution
make test
make test-coverage
make test-integration
make test-security
```

---

## 📝 Test Documentation Standards

### Test Case Format

Each test case should include:
1. **Test ID**: Unique identifier (e.g., TC-AUTH-001)
2. **Test Scenario**: Brief description of what is being tested
3. **Test Steps**: Detailed step-by-step instructions
4. **Expected Result**: Clear description of expected outcome
5. **Setup Requirements**: Any necessary setup or prerequisites
6. **Cleanup**: Any necessary cleanup after test execution

### Test Data Management

- Use consistent test data across all test cases
- Implement proper test data cleanup
- Use factories for generating test data
- Maintain test data isolation between test cases

### Error Testing

- Test all error conditions and edge cases
- Verify proper error messages and status codes
- Test error handling at all layers (controller, service, repository)
- Include security-related error scenarios

---

**Last Updated**: October 29, 2025  
**Test Suite Version**: 2.0.0  
**Go Version**: 1.21+  
**Test Framework**: testify/assert, testify/require