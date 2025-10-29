# 📋 Freshease Backend - Complete Test Case Specification

## Test Case Inventory

This document provides a comprehensive list of all automated test cases in the Freshease backend test suite.

**Total Test Files**: 40  
**Total Test Cases**: 200+ individual test scenarios  
**Modules Covered**: 13  

---

## 🔐 Authentication Module (`auth/authoidc`)

### Controller Tests (`controller_test.go`)

#### TestController_Start
- **TC-AUTH-001**: `error - unknown provider`
  - **Input**: Provider = "unknown"
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests handling of unsupported OIDC providers

- **TC-AUTH-002**: `error - empty provider`
  - **Input**: Provider = ""
  - **Expected**: HTTP 404 Not Found
  - **Description**: Tests handling of empty provider parameter

#### TestController_Callback
- **TC-AUTH-003**: `error - missing state`
  - **Input**: State = "", Code = "test-code"
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests OAuth callback with missing state parameter

- **TC-AUTH-004**: `error - state mismatch`
  - **Input**: State = "different-state", Cookie = "test-state"
  - **Expected**: HTTP 401 Unauthorized
  - **Description**: Tests OAuth state validation failure

- **TC-AUTH-005**: `error - unknown provider`
  - **Input**: Provider = "unknown", Valid state/code
  - **Expected**: HTTP 307 Temporary Redirect
  - **Description**: Tests callback with unsupported provider

#### TestController_Integration
- **TC-AUTH-006**: `test error cases`
  - **Description**: Integration test for OAuth flow error scenarios

### Service Tests (`service_test.go`)

#### TestService_AuthCodeURL
- **TC-AUTH-007**: `success - google provider`
  - **Input**: Provider = "google"
  - **Expected**: Valid OAuth URL
  - **Description**: Tests Google OAuth URL generation

- **TC-AUTH-008**: `success - line provider`
  - **Input**: Provider = "line"
  - **Expected**: Valid OAuth URL
  - **Description**: Tests LINE OAuth URL generation

- **TC-AUTH-009**: `error - unknown provider`
  - **Input**: Provider = "unknown"
  - **Expected**: Error
  - **Description**: Tests error handling for unknown providers

#### TestService_ExchangeAndLogin
- **TC-AUTH-010**: `error - unknown provider`
  - **Input**: Provider = "unknown", Valid code
  - **Expected**: Error
  - **Description**: Tests token exchange error handling

#### TestService_IssueJWT
- **TC-AUTH-011**: `success - issue JWT`
  - **Input**: Valid user ID and email
  - **Expected**: Valid JWT token
  - **Description**: Tests JWT token generation

---

## 🛒 Cart Items Module (`cart_items`)

### Controller Tests (`controller_test.go`)

#### TestController_ListCart_items
- **TC-CART-001**: `success - returns cart items list`
  - **Input**: Valid request
  - **Expected**: HTTP 200, Cart items array
  - **Description**: Tests successful cart items listing

- **TC-CART-002**: `error - service returns error`
  - **Input**: Service error
  - **Expected**: HTTP 500 Internal Server Error
  - **Description**: Tests service error handling

#### TestController_GetCart_item
- **TC-CART-003**: `success - returns cart item by ID`
  - **Input**: Valid UUID
  - **Expected**: HTTP 200, Cart item data
  - **Description**: Tests successful cart item retrieval

- **TC-CART-004**: `error - invalid UUID`
  - **Input**: Invalid UUID string
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests invalid UUID handling

- **TC-CART-005**: `error - cart item not found`
  - **Input**: Valid UUID, non-existent item
  - **Expected**: HTTP 404 Not Found
  - **Description**: Tests not found scenario

#### TestController_CreateCart_item
- **TC-CART-006**: `success - creates new cart item`
  - **Input**: Valid cart item data
  - **Expected**: HTTP 201 Created
  - **Description**: Tests successful cart item creation

- **TC-CART-007**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests creation error handling

#### TestController_UpdateCart_item
- **TC-CART-008**: `success - updates cart item`
  - **Input**: Valid UUID and update data
  - **Expected**: HTTP 201 Created
  - **Description**: Tests successful cart item update

- **TC-CART-009**: `error - invalid UUID`
  - **Input**: Invalid UUID string
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests invalid UUID in update

- **TC-CART-010**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests update error handling

#### TestController_DeleteCart_item
- **TC-CART-011**: `success - deletes cart item`
  - **Input**: Valid UUID
  - **Expected**: HTTP 202 Accepted
  - **Description**: Tests successful cart item deletion

- **TC-CART-012**: `error - invalid UUID`
  - **Input**: Invalid UUID string
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests invalid UUID in deletion

- **TC-CART-013**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests deletion error handling

#### TestController_EdgeCases
- **TC-CART-014**: `empty cart items list`
  - **Input**: Empty cart
  - **Expected**: HTTP 200, Empty array
  - **Description**: Tests empty cart scenario

### Service Tests (`service_test.go`)

#### TestService_List
- **TC-CART-015**: `success - returns cart items list`
- **TC-CART-016**: `error - repository returns error`

#### TestService_Get
- **TC-CART-017**: `success - returns cart item by ID`
- **TC-CART-018**: `error - cart item not found`

#### TestService_Create
- **TC-CART-019**: `success - creates new cart item`
- **TC-CART-020**: `error - repository returns error`

#### TestService_Update
- **TC-CART-021**: `success - updates cart item`
- **TC-CART-022**: `error - repository returns error`

#### TestService_Delete
- **TC-CART-023**: `success - deletes cart item`
- **TC-CART-024**: `error - repository returns error`

---

## 🤖 GenAI Module (`genai`)

### Controller Tests (`controller_test.go`)

#### TestController_GenerateWeeklyMeals
- **TC-GENAI-001**: `success - generates weekly meals`
  - **Input**: Complete user profile data
  - **Expected**: HTTP 200, Weekly meal plan
  - **Description**: Tests successful weekly meal generation

- **TC-GENAI-002**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 502 Bad Gateway
  - **Description**: Tests meal generation error handling

#### TestController_GenerateDailyMeals
- **TC-GENAI-003**: `success - generates daily meals`
  - **Input**: Complete user profile data
  - **Expected**: HTTP 200, Daily meal plan
  - **Description**: Tests successful daily meal generation

- **TC-GENAI-004**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 502 Bad Gateway
  - **Description**: Tests daily meal generation error handling

#### TestController_RequestParsing
- **TC-GENAI-005**: `handles malformed JSON`
  - **Input**: Invalid JSON payload
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests malformed JSON handling

- **TC-GENAI-006**: `handles empty request body`
  - **Input**: Empty request body
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests empty request handling

#### TestController_EdgeCases
- **TC-GENAI-007**: `weekly meals with minimal data`
  - **Input**: Minimal user profile
  - **Expected**: HTTP 200, Basic meal plan
  - **Description**: Tests meal generation with minimal data

- **TC-GENAI-008**: `daily meals with extreme values`
  - **Input**: Extreme user profile values
  - **Expected**: HTTP 200, Adapted meal plan
  - **Description**: Tests meal generation with extreme values

#### TestController_ErrorHandling
- **TC-GENAI-009**: `service returns nil result`
  - **Input**: Service returns nil
  - **Expected**: HTTP 200, Empty result
  - **Description**: Tests nil result handling

- **TC-GENAI-010**: `service returns timeout error`
  - **Input**: Service timeout
  - **Expected**: HTTP 502 Bad Gateway
  - **Description**: Tests timeout error handling

### Service Tests (`service_test.go`)

#### TestService_GenerateWeeklyMeals
- **TC-GENAI-011**: `success - generates weekly meals with complete profile`
- **TC-GENAI-012**: `success - generates weekly meals with incomplete profile, loads from repo`
- **TC-GENAI-013**: `success - generates weekly meals without user ID (no repo calls)`
- **TC-GENAI-014**: `error - repository GetUserProfile fails`

#### TestService_GenerateDailyMeals
- **TC-GENAI-015**: `success - generates daily meals with complete profile`
- **TC-GENAI-016**: `success - generates daily meals with incomplete profile, loads from repo`
- **TC-GENAI-017**: `success - generates daily meals without user ID (no repo calls)`
- **TC-GENAI-018**: `error - repository GetUserProfile fails`

#### TestService_ProfileLoading
- **TC-GENAI-019**: `profile loading fills missing fields correctly`

### Repository Tests (`repo_test.go`)

#### TestRepository_GetUserProfile
- **TC-GENAI-020**: `success - returns existing profile`
- **TC-GENAI-021**: `error - profile not found`

#### TestRepository_SaveGeneratedPlan
- **TC-GENAI-022**: `success - saves plan`
- **TC-GENAI-023**: `success - saves empty plan`
- **TC-GENAI-024**: `success - saves nil plan`

#### TestRepository_Integration
- **TC-GENAI-025**: `complete workflow - save profile and plan`
- **TC-GENAI-026**: `multiple users workflow`

---

## 👥 Users Module (`users`)

### Controller Tests (`controller_test.go`)

#### TestController_ListUsers
- **TC-USER-001**: `success - returns users list`
  - **Input**: Valid request
  - **Expected**: HTTP 200, Users array
  - **Description**: Tests successful users listing

- **TC-USER-002**: `error - service returns error`
  - **Input**: Service error
  - **Expected**: HTTP 500 Internal Server Error
  - **Description**: Tests service error handling

#### TestController_GetUser
- **TC-USER-003**: `success - returns user by ID`
  - **Input**: Valid UUID
  - **Expected**: HTTP 200, User data
  - **Description**: Tests successful user retrieval

- **TC-USER-004**: `error - invalid UUID`
  - **Input**: Invalid UUID string
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests invalid UUID handling

- **TC-USER-005**: `error - user not found`
  - **Input**: Valid UUID, non-existent user
  - **Expected**: HTTP 404 Not Found
  - **Description**: Tests not found scenario

#### TestController_CreateUser
- **TC-USER-006**: `success - creates new user`
  - **Input**: Valid user data
  - **Expected**: HTTP 201 Created
  - **Description**: Tests successful user creation

- **TC-USER-007**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests creation error handling

#### TestController_UpdateUser
- **TC-USER-008**: `success - updates user`
  - **Input**: Valid UUID and update data
  - **Expected**: HTTP 201 Created
  - **Description**: Tests successful user update

- **TC-USER-009**: `error - invalid UUID`
  - **Input**: Invalid UUID string
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests invalid UUID in update

- **TC-USER-010**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests update error handling

#### TestController_DeleteUser
- **TC-USER-011**: `success - deletes user`
  - **Input**: Valid UUID
  - **Expected**: HTTP 202 Accepted
  - **Description**: Tests successful user deletion

- **TC-USER-012**: `error - invalid UUID`
  - **Input**: Invalid UUID string
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests invalid UUID in deletion

- **TC-USER-013**: `error - service returns error`
  - **Input**: Service error scenario
  - **Expected**: HTTP 400 Bad Request
  - **Description**: Tests deletion error handling

### Repository Tests (`repo_test.go`)

#### TestEntRepo_List
- **TC-USER-014**: `success - returns empty list when no users`
- **TC-USER-015**: `success - returns users list`

#### TestEntRepo_FindByID
- **TC-USER-016**: `success - returns user by ID`
- **TC-USER-017**: `error - user not found`

#### TestEntRepo_Create
- **TC-USER-018**: `success - creates new user`
- **TC-USER-019**: `success - creates user with minimal data`
- **TC-USER-020**: `error - duplicate email`

#### TestEntRepo_Update
- **TC-USER-021**: `success - updates user`
- **TC-USER-022**: `success - partial update`
- **TC-USER-023**: `error - no fields to update`
- **TC-USER-024**: `error - user not found`

#### TestEntRepo_Delete
- **TC-USER-025**: `success - deletes user`
- **TC-USER-026**: `error - user not found`

### Service Tests (`service_test.go`)

#### TestService_List
- **TC-USER-027**: `success - returns users list`
- **TC-USER-028**: `error - repository returns error`

#### TestService_Get
- **TC-USER-029**: `success - returns user by ID`
- **TC-USER-030**: `error - user not found`

#### TestService_Create
- **TC-USER-031**: `success - creates new user`
- **TC-USER-032**: `error - repository returns error`

#### TestService_Update
- **TC-USER-033**: `success - updates user`
- **TC-USER-034**: `error - repository returns error`

#### TestService_Delete
- **TC-USER-035**: `success - deletes user`
- **TC-USER-036**: `error - repository returns error`

---

## 🛍️ Products Module (`products`)

### Controller Tests (`controller_test.go`)

#### TestController_ListProducts
- **TC-PRODUCT-001**: `success - returns products list`
- **TC-PRODUCT-002**: `error - service returns error`

#### TestController_GetProduct
- **TC-PRODUCT-003**: `success - returns product by ID`
- **TC-PRODUCT-004**: `error - invalid UUID`
- **TC-PRODUCT-005**: `error - product not found`

#### TestController_CreateProduct
- **TC-PRODUCT-006**: `success - creates new product`
- **TC-PRODUCT-007**: `error - service returns error`

#### TestController_UpdateProduct
- **TC-PRODUCT-008**: `success - updates product`
- **TC-PRODUCT-009**: `error - invalid UUID`
- **TC-PRODUCT-010**: `error - service returns error`

#### TestController_DeleteProduct
- **TC-PRODUCT-011**: `success - deletes product`
- **TC-PRODUCT-012**: `error - invalid UUID`
- **TC-PRODUCT-013**: `error - service returns error`

### Service Tests (`service_test.go`)

#### TestService_List
- **TC-PRODUCT-014**: `success - returns products list`
- **TC-PRODUCT-015**: `error - repository returns error`

#### TestService_Get
- **TC-PRODUCT-016**: `success - returns product by ID`
- **TC-PRODUCT-017**: `error - product not found`

#### TestService_Create
- **TC-PRODUCT-018**: `success - creates new product`
- **TC-PRODUCT-019**: `error - repository returns error`

#### TestService_Update
- **TC-PRODUCT-020**: `success - updates product`
- **TC-PRODUCT-021**: `error - repository returns error`

#### TestService_Delete
- **TC-PRODUCT-022**: `success - deletes product`
- **TC-PRODUCT-023**: `error - repository returns error`

---

## 🔧 Middleware Tests (`internal/common/middleware`)

### TestBindAndValidate
- **TC-MIDDLEWARE-001**: `success - valid create DTO`
  - **Input**: Valid JSON payload
  - **Expected**: Successful binding and validation
  - **Description**: Tests middleware binding and validation

### TestRequireAuth
- **TC-MIDDLEWARE-002**: `success - valid token`
  - **Input**: Valid JWT token
  - **Expected**: Authentication success
  - **Description**: Tests valid token authentication

- **TC-MIDDLEWARE-003**: `error - missing authorization header`
  - **Input**: No Authorization header
  - **Expected**: Authentication failure
  - **Description**: Tests missing header handling

- **TC-MIDDLEWARE-004**: `error - invalid bearer prefix`
  - **Input**: Invalid Bearer prefix
  - **Expected**: Authentication failure
  - **Description**: Tests invalid Bearer format

- **TC-MIDDLEWARE-005**: `error - invalid token`
  - **Input**: Invalid JWT token
  - **Expected**: Authentication failure
  - **Description**: Tests invalid token handling

- **TC-MIDDLEWARE-006**: `error - expired token`
  - **Input**: Expired JWT token
  - **Expected**: Authentication failure
  - **Description**: Tests expired token handling

### TestRequestLogger
- **TC-MIDDLEWARE-007**: `Request logging`
  - **Description**: Tests request logging functionality

- **TC-MIDDLEWARE-008**: `Request logging with error`
  - **Description**: Tests error logging functionality

---

## 🔗 Integration Tests (`internal/common/testutils`)

### TestUsersAPI_Integration
- **TC-INTEGRATION-001**: `Create and retrieve user`
  - **Description**: End-to-end user creation and retrieval

- **TC-INTEGRATION-002**: `List users`
  - **Description**: End-to-end user listing

- **TC-INTEGRATION-003**: `Update user`
  - **Description**: End-to-end user update

- **TC-INTEGRATION-004**: `Delete user`
  - **Description**: End-to-end user deletion

### TestProductsAPI_Integration
- **TC-INTEGRATION-005**: `Create and retrieve product`
  - **Description**: End-to-end product creation and retrieval

- **TC-INTEGRATION-006**: `List products`
  - **Description**: End-to-end product listing

- **TC-INTEGRATION-007**: `Update product`
  - **Description**: End-to-end product update

- **TC-INTEGRATION-008**: `Delete product`
  - **Description**: End-to-end product deletion

### TestAPI_ErrorHandling
- **TC-INTEGRATION-009**: `Invalid UUID format`
  - **Description**: Tests API error handling for invalid UUIDs

- **TC-INTEGRATION-010**: `Invalid JSON payload`
  - **Description**: Tests API error handling for malformed JSON

- **TC-INTEGRATION-011**: `Validation errors`
  - **Description**: Tests API validation error handling

### TestAPI_NotFoundHandling
- **TC-INTEGRATION-012**: `Non-existent user`
  - **Description**: Tests API handling of non-existent resources

---

## 📊 Test Coverage Summary

### Module Coverage Breakdown
| Module | Test Cases | Coverage | Status |
|--------|------------|----------|--------|
| **Authentication** | 11 | 26.8% | ✅ PASS |
| **Cart Items** | 24 | 41.0% | ✅ PASS |
| **GenAI** | 26 | 41.2% | ✅ PASS |
| **Users** | 36 | 74.6% | ✅ PASS |
| **Products** | 23 | 55.7% | ✅ PASS |
| **Middleware** | 8 | 91.7% | ✅ PASS |
| **Integration** | 12 | 0.0% | ✅ PASS |
| **Other Modules** | 60+ | 55-82% | ✅ PASS |

### Test Categories
- **Unit Tests**: 180+ test cases
- **Integration Tests**: 12 test cases
- **Edge Case Tests**: 20+ test cases
- **Error Handling Tests**: 50+ test cases

### Test Execution Status
- ✅ **All Tests Passing**: 100% success rate
- ✅ **Coverage Threshold**: 43.5% overall coverage
- ✅ **CI/CD Ready**: Automated test execution
- ✅ **Documentation**: Complete test case documentation

---

## 🎯 Test Case Naming Convention

### Format
```
TC-{MODULE}-{NUMBER}: {test_name}
```

### Module Codes
- **AUTH**: Authentication
- **CART**: Cart Items
- **GENAI**: GenAI
- **USER**: Users
- **PRODUCT**: Products
- **MIDDLEWARE**: Middleware
- **INTEGRATION**: Integration Tests

### Test Types
- **Success Cases**: Normal operation scenarios
- **Error Cases**: Error handling scenarios
- **Edge Cases**: Boundary and special conditions
- **Integration**: End-to-end workflows

---

**Total Test Cases**: 200+  
**Last Updated**: October 29, 2025  
**Status**: ✅ All Tests Passing
