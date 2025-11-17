# Freshease Test Cases

## Test Case Format
- **Test Case ID**: Unique identifier (e.g., TC-BE-001)
- **Title**: Brief description
- **Description**: Detailed test scenario
- **Pre-conditions**: Required setup
- **Steps**: Test execution steps
- **Expected Result**: Expected outcome
- **Priority**: P0 (Critical), P1 (High), P2 (Medium), P3 (Low)

---

## Backend Test Cases

### Authentication & Authorization

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-001 | User Registration - Valid Data | Test user registration with valid input | None | 1. POST /api/auth/password/register<br>2. Send valid user data (email, password, name) | Status 201, user created with UUID, email verified | P0 |
| TC-BE-002 | User Registration - Invalid Email | Test registration with invalid email format | None | 1. POST /api/auth/password/register<br>2. Send invalid email (e.g., "invalid-email") | Status 400, validation error message | P0 |
| TC-BE-003 | User Registration - Weak Password | Test registration with password < 8 characters | None | 1. POST /api/auth/password/register<br>2. Send password with 7 characters | Status 400, password validation error | P0 |
| TC-BE-004 | User Login - Valid Credentials | Test login with correct email and password | User exists in database | 1. POST /api/auth/password/login<br>2. Send valid credentials | Status 200, JWT token returned | P0 |
| TC-BE-005 | User Login - Invalid Credentials | Test login with wrong password | User exists in database | 1. POST /api/auth/password/login<br>2. Send wrong password | Status 401, unauthorized error | P0 |
| TC-BE-006 | OIDC Authentication - Google | Test Google OAuth flow | OIDC provider configured | 1. GET /api/auth/oidc/google<br>2. Complete OAuth flow | Redirect to callback with code, user authenticated | P1 |
| TC-BE-007 | Protected Endpoint - Valid Token | Test accessing protected endpoint with valid JWT | User logged in, token obtained | 1. GET /api/whoami<br>2. Include Authorization: Bearer {token} | Status 200, user data returned | P0 |
| TC-BE-008 | Protected Endpoint - Invalid Token | Test accessing protected endpoint with invalid token | None | 1. GET /api/whoami<br>2. Include invalid token | Status 401, unauthorized error | P0 |
| TC-BE-009 | Protected Endpoint - Missing Token | Test accessing protected endpoint without token | None | 1. GET /api/whoami<br>2. No Authorization header | Status 401, unauthorized error | P0 |

### User Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-010 | Get User - Valid ID | Retrieve user by ID | User exists | 1. GET /api/users/{id} | Status 200, user data returned | P0 |
| TC-BE-011 | Get User - Invalid ID | Retrieve user with non-existent ID | None | 1. GET /api/users/{invalid-uuid} | Status 404, user not found | P1 |
| TC-BE-012 | List Users | Retrieve all users | Multiple users exist | 1. GET /api/users | Status 200, array of users returned | P1 |
| TC-BE-013 | Update User - Valid Data | Update user information | User exists, authenticated | 1. PUT /api/users/{id}<br>2. Send update data | Status 200, updated user data | P0 |
| TC-BE-014 | Update User - Invalid Email | Update user with invalid email | User exists | 1. PUT /api/users/{id}<br>2. Send invalid email | Status 400, validation error | P1 |
| TC-BE-015 | Delete User - Valid ID | Delete user | User exists, authenticated | 1. DELETE /api/users/{id} | Status 200/204, user deleted | P1 |
| TC-BE-016 | Upload User Avatar | Upload user profile picture | User exists, authenticated | 1. POST /api/users/{id}/avatar<br>2. Upload image file | Status 200, image URL returned | P1 |
| TC-BE-017 | Upload User Avatar - Invalid File | Upload non-image file as avatar | User exists | 1. POST /api/users/{id}/avatar<br>2. Upload text file | Status 400, invalid file type error | P2 |

### Product Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-018 | Create Product - Valid Data | Create new product | Vendor exists | 1. POST /api/products<br>2. Send valid product data | Status 201, product created with UUID | P0 |
| TC-BE-019 | Create Product - Missing Required Fields | Create product without required fields | None | 1. POST /api/products<br>2. Send data without name | Status 400, validation error | P0 |
| TC-BE-020 | Create Product - Negative Price | Create product with negative price | None | 1. POST /api/products<br>2. Send price = -10 | Status 400, price validation error | P0 |
| TC-BE-021 | Get Product - Valid ID | Retrieve product by ID | Product exists | 1. GET /api/products/{id} | Status 200, product data returned | P0 |
| TC-BE-022 | List Products | Retrieve all products | Multiple products exist | 1. GET /api/products | Status 200, array of products | P0 |
| TC-BE-023 | Update Product - Valid Data | Update product information | Product exists | 1. PUT /api/products/{id}<br>2. Send update data | Status 200, updated product | P0 |
| TC-BE-024 | Delete Product | Delete product | Product exists | 1. DELETE /api/products/{id} | Status 200/204, product deleted | P1 |
| TC-BE-025 | Search Products - By Category | Search products filtered by category | Products with categories exist | 1. GET /api/shop/products?category_id={id} | Status 200, filtered products | P0 |
| TC-BE-026 | Search Products - By Price Range | Search products within price range | Products exist | 1. GET /api/shop/products?min_price=10&max_price=50 | Status 200, products in range | P1 |
| TC-BE-027 | Search Products - By Vendor | Search products by vendor | Products from vendor exist | 1. GET /api/shop/products?vendor_id={id} | Status 200, vendor products | P1 |

### Cart Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-028 | Create Cart | Create new shopping cart | User exists | 1. POST /api/carts<br>2. Send cart data | Status 201, cart created | P0 |
| TC-BE-029 | Get Cart - Valid ID | Retrieve cart by ID | Cart exists | 1. GET /api/carts/{id} | Status 200, cart data with items | P0 |
| TC-BE-030 | Add Item to Cart | Add product to cart | Cart and product exist | 1. POST /api/cart-items<br>2. Send cart_id, product_id, quantity | Status 201, cart item created | P0 |
| TC-BE-031 | Add Item to Cart - Invalid Quantity | Add item with quantity 0 | Cart and product exist | 1. POST /api/cart-items<br>2. Send quantity = 0 | Status 400, validation error | P0 |
| TC-BE-032 | Update Cart Item Quantity | Update quantity of cart item | Cart item exists | 1. PUT /api/cart-items/{id}<br>2. Send new quantity | Status 200, quantity updated | P0 |
| TC-BE-033 | Remove Item from Cart | Delete cart item | Cart item exists | 1. DELETE /api/cart-items/{id} | Status 200/204, item removed | P0 |
| TC-BE-034 | Get Cart Total | Calculate cart total | Cart with items exists | 1. GET /api/carts/{id} | Status 200, total calculated correctly | P0 |
| TC-BE-035 | Clear Cart | Remove all items from cart | Cart with items exists | 1. Delete all cart items | All items removed, cart empty | P1 |

### Order Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-036 | Create Order - Valid Data | Create order from cart | Cart with items exists, user authenticated | 1. POST /api/orders<br>2. Send order data | Status 201, order created | P0 |
| TC-BE-037 | Create Order - Empty Cart | Create order with no items | Empty cart exists | 1. POST /api/orders<br>2. Send order data | Status 400, validation error | P0 |
| TC-BE-038 | Get Order - Valid ID | Retrieve order by ID | Order exists | 1. GET /api/orders/{id} | Status 200, order with items returned | P0 |
| TC-BE-039 | List Orders - By User | Get all orders for user | User has orders | 1. GET /api/orders?user_id={id} | Status 200, user orders returned | P0 |
| TC-BE-040 | Update Order Status | Update order status | Order exists | 1. PUT /api/orders/{id}<br>2. Send new status | Status 200, status updated | P0 |
| TC-BE-041 | Cancel Order | Cancel pending order | Order with status "pending" exists | 1. PUT /api/orders/{id}<br>2. Set status to "cancelled" | Status 200, order cancelled | P1 |
| TC-BE-042 | Calculate Order Total | Verify order total calculation | Order with items exists | 1. GET /api/orders/{id} | Total = subtotal + shipping - discount | P0 |

### Payment Processing

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-043 | Create Payment - Valid Data | Create payment for order | Order exists | 1. POST /api/payments<br>2. Send payment data | Status 201, payment created | P0 |
| TC-BE-044 | Get Payment - Valid ID | Retrieve payment by ID | Payment exists | 1. GET /api/payments/{id} | Status 200, payment data | P1 |
| TC-BE-045 | Update Payment Status | Update payment status | Payment exists | 1. PUT /api/payments/{id}<br>2. Send new status | Status 200, status updated | P1 |

### GenAI Integration

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-046 | Generate Meal Plan | Generate weekly meal plan | GenAI API configured | 1. POST /api/genai/meal-plan<br>2. Send user preferences | Status 200, meal plan generated | P1 |
| TC-BE-047 | Generate Recipe | Generate recipe from ingredients | GenAI API configured | 1. POST /api/genai/recipe<br>2. Send ingredients list | Status 200, recipe generated | P2 |

### File Upload

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-BE-048 | Upload Image - Valid File | Upload image file | MinIO configured | 1. POST /api/uploads<br>2. Upload image file | Status 200, presigned URL returned | P1 |
| TC-BE-049 | Upload Image - Invalid File Type | Upload non-image file | None | 1. POST /api/uploads<br>2. Upload text file | Status 400, invalid file type | P2 |
| TC-BE-050 | Get Presigned URL | Get presigned URL for image | Image exists in MinIO | 1. GET /api/uploads/{object_name} | Status 200, presigned URL | P1 |

---

## Frontend (Flutter) Test Cases

### Authentication

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-FE-001 | Login - Valid Credentials | User logs in with correct credentials | App installed, user exists | 1. Enter email and password<br>2. Tap login button | User authenticated, navigated to home | P0 |
| TC-FE-002 | Login - Invalid Credentials | User logs in with wrong password | App installed | 1. Enter email and wrong password<br>2. Tap login button | Error message displayed, stay on login page | P0 |
| TC-FE-003 | Sign Up - Valid Data | New user registers | App installed | 1. Fill registration form<br>2. Submit | User created, logged in, navigated to home | P0 |
| TC-FE-004 | Sign Up - Validation Errors | User submits invalid data | App installed | 1. Submit form with invalid email<br>2. Submit | Validation errors displayed | P0 |
| TC-FE-005 | OAuth Login - Google | User logs in with Google | OAuth configured | 1. Tap "Login with Google"<br>2. Complete OAuth flow | User authenticated via OAuth | P1 |

### Product Browsing

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-FE-006 | Browse Products | User views product list | User logged in | 1. Navigate to shop page | Products displayed in list/grid | P0 |
| TC-FE-007 | Search Products | User searches for products | Products exist | 1. Enter search term<br>2. Tap search | Filtered products displayed | P0 |
| TC-FE-008 | Filter Products - Category | Filter products by category | Products with categories exist | 1. Select category filter | Products filtered by category | P1 |
| TC-FE-009 | View Product Details | User views product details | Product exists | 1. Tap on product | Product details page displayed | P0 |
| TC-FE-010 | Add to Cart from Details | Add product to cart from details page | Product exists, user logged in | 1. View product details<br>2. Tap "Add to Cart" | Product added, cart updated | P0 |

### Cart Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-FE-011 | View Cart | User views shopping cart | Cart has items | 1. Navigate to cart page | Cart items displayed with totals | P0 |
| TC-FE-012 | Update Item Quantity | Change quantity in cart | Cart item exists | 1. Tap quantity controls<br>2. Update quantity | Quantity updated, total recalculated | P0 |
| TC-FE-013 | Remove Item from Cart | Remove item from cart | Cart item exists | 1. Tap remove button | Item removed, cart updated | P0 |
| TC-FE-014 | Empty Cart | Clear all items from cart | Cart has items | 1. Tap clear cart | All items removed | P1 |

### Checkout

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-FE-015 | Checkout - Valid Address | Complete checkout with valid address | Cart has items, user logged in | 1. Navigate to checkout<br>2. Select address<br>3. Complete payment | Order created, confirmation shown | P0 |
| TC-FE-016 | Checkout - Add New Address | Add address during checkout | Cart has items | 1. Navigate to checkout<br>2. Add new address<br>3. Complete order | Address saved, order created | P1 |
| TC-FE-017 | Checkout - Empty Cart | Attempt checkout with empty cart | Cart is empty | 1. Navigate to checkout | Error message, redirect to cart | P0 |

### User Profile

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-FE-018 | View Profile | User views own profile | User logged in | 1. Navigate to profile page | Profile information displayed | P0 |
| TC-FE-019 | Edit Profile | Update profile information | User logged in | 1. Tap edit profile<br>2. Update fields<br>3. Save | Profile updated, changes saved | P1 |
| TC-FE-020 | Upload Profile Picture | Change profile avatar | User logged in | 1. Tap avatar<br>2. Select image<br>3. Upload | Avatar updated | P2 |

---

## Frontend-Admin Test Cases

### Authentication

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-ADMIN-001 | Admin Login - Valid Credentials | Admin logs in | Admin account exists | 1. Enter credentials<br>2. Submit | Admin authenticated, dashboard shown | P0 |
| TC-ADMIN-002 | Admin Login - Invalid Credentials | Admin login with wrong password | None | 1. Enter wrong password<br>2. Submit | Error message, stay on login | P0 |

### User Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-ADMIN-003 | List Users | View all users | Admin logged in | 1. Navigate to users page | User table displayed | P0 |
| TC-ADMIN-004 | Create User | Create new user | Admin logged in | 1. Click "Create User"<br>2. Fill form<br>3. Submit | User created, added to list | P0 |
| TC-ADMIN-005 | Edit User | Update user information | User exists | 1. Click edit on user<br>2. Update fields<br>3. Save | User updated | P0 |
| TC-ADMIN-006 | Delete User | Delete user | User exists | 1. Click delete on user<br>2. Confirm | User deleted, removed from list | P1 |

### Product Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-ADMIN-007 | List Products | View all products | Admin logged in | 1. Navigate to products page | Product table displayed | P0 |
| TC-ADMIN-008 | Create Product | Create new product | Admin logged in, vendor exists | 1. Click "Create Product"<br>2. Fill form<br>3. Submit | Product created | P0 |
| TC-ADMIN-009 | Edit Product | Update product | Product exists | 1. Click edit<br>2. Update fields<br>3. Save | Product updated | P0 |
| TC-ADMIN-010 | Delete Product | Delete product | Product exists | 1. Click delete<br>2. Confirm | Product deleted | P1 |

### Order Management

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-ADMIN-011 | List Orders | View all orders | Admin logged in | 1. Navigate to orders page | Order table displayed | P0 |
| TC-ADMIN-012 | View Order Details | View order information | Order exists | 1. Click on order | Order details displayed | P0 |
| TC-ADMIN-013 | Update Order Status | Change order status | Order exists | 1. Click edit<br>2. Change status<br>3. Save | Status updated | P0 |

### Analytics

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-ADMIN-014 | View Dashboard | View analytics dashboard | Admin logged in | 1. Navigate to dashboard | Analytics data displayed | P1 |
| TC-ADMIN-015 | View Sales Report | View sales analytics | Orders exist | 1. Navigate to analytics<br>2. View sales | Sales data displayed | P2 |

---

## Integration Test Cases

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-INT-001 | End-to-End Order Flow | Complete order from product selection to payment | User logged in, products exist | 1. Browse products<br>2. Add to cart<br>3. Checkout<br>4. Complete payment | Order created, payment processed | P0 |
| TC-INT-002 | User Registration to First Order | New user registers and places order | None | 1. Register user<br>2. Browse products<br>3. Add to cart<br>4. Place order | User created, order placed | P0 |
| TC-INT-003 | Cart Persistence | Cart persists across sessions | User logged in | 1. Add items to cart<br>2. Logout<br>3. Login again | Cart items still present | P1 |
| TC-INT-004 | Image Upload Flow | Upload and retrieve product image | Product exists | 1. Upload image<br>2. Get presigned URL<br>3. Display image | Image uploaded and displayed | P1 |

---

## Performance Test Cases

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-PERF-001 | API Response Time | API endpoints respond within acceptable time | Server running | 1. Make API calls<br>2. Measure response time | Response time < 500ms for 95% of requests | P1 |
| TC-PERF-002 | Product List Loading | Product list loads quickly | Products exist | 1. Load product list<br>2. Measure load time | List loads in < 2 seconds | P1 |
| TC-PERF-003 | Database Query Performance | Database queries are optimized | Database populated | 1. Execute queries<br>2. Measure execution time | Queries complete in < 100ms | P2 |

---

## Security Test Cases

| ID | Title | Description | Pre-conditions | Steps | Expected Result | Priority |
|----|-------|-------------|----------------|-------|-----------------|----------|
| TC-SEC-001 | SQL Injection Prevention | Test SQL injection protection | None | 1. Send SQL injection payload in input | Request rejected, no SQL executed | P0 |
| TC-SEC-002 | XSS Prevention | Test XSS protection | None | 1. Send XSS payload in input | Payload sanitized, not executed | P0 |
| TC-SEC-003 | CSRF Protection | Test CSRF token validation | None | 1. Send request without CSRF token | Request rejected | P1 |
| TC-SEC-004 | Input Validation | Test input validation on all endpoints | None | 1. Send invalid/malicious input | Request rejected with validation error | P0 |
| TC-SEC-005 | Authorization Check | Test unauthorized access prevention | User logged in | 1. Attempt to access other user's data | Access denied, 403 error | P0 |

