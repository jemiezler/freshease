package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"freshease/backend/ent/enttest"
	"freshease/backend/internal/common/config"
	httpserver "freshease/backend/internal/common/http"
	"freshease/backend/modules/uploads"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUploadsService for integration tests
type MockUploadsService struct {
	mock.Mock
}

func (m *MockUploadsService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	args := m.Called(ctx, file, folder)
	return args.String(0), args.Error(1)
}

func (m *MockUploadsService) DeleteImage(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

func (m *MockUploadsService) GetImageURL(ctx context.Context, objectName string) (string, error) {
	args := m.Called(ctx, objectName)
	return args.String(0), args.Error(1)
}

func (m *MockUploadsService) GetImage(ctx context.Context, objectName string) (io.ReadCloser, *minio.ObjectInfo, error) {
	args := m.Called(ctx, objectName)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(io.ReadCloser), nil, args.Error(2)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(*minio.ObjectInfo), args.Error(2)
}

// TestFullSystem_OrderFlow tests the complete order flow from product creation to payment
func TestFullSystem_OrderFlow(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	// Create test app with all modules
	app := fiber.New()
	cfg := config.Config{
		HTTPPort: ":3000",
		DatabaseURL: "sqlite3://:memory:",
	}
	mockUploads := new(MockUploadsService)
	apiGroup := app.Group("/api")
	httpserver.RegisterRoutes(apiGroup, app, client, cfg)

	ctx := context.Background()

	t.Run("Complete order flow", func(t *testing.T) {
		// Step 1: Create a user
		user, err := client.User.Create().
			SetID(uuid.New()).
			SetEmail("testuser@example.com").
			SetName("Test User").
			SetPassword("password123").
			Save(ctx)
		require.NoError(t, err)

		// Step 2: Create a vendor
		vendor, err := client.Vendor.Create().
			SetID(uuid.New()).
			SetName("Test Vendor").
			SetContact("vendor@example.com").
			Save(ctx)
		require.NoError(t, err)

		// Step 3: Create a category
		category, err := client.Category.Create().
			SetID(uuid.New()).
			SetName("Test Category").
			SetSlug("test-category").
			Save(ctx)
		require.NoError(t, err)

		// Step 4: Create a product
		product, err := client.Product.Create().
			SetID(uuid.New()).
			SetName("Test Product").
			SetSku("TEST-001").
			SetPrice(99.99).
			SetDescription("Test product description").
			SetUnitLabel("kg").
			SetIsActive(true).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetVendor(vendor).
			Save(ctx)
		require.NoError(t, err)

		// Step 5: Create product category relationship
		_, err = client.Product_category.Create().
			SetProduct(product).
			SetCategory(category).
			Save(ctx)
		require.NoError(t, err)

		// Step 6: Create inventory
		_, err = client.Inventory.Create().
			SetQuantity(100).
			SetReorderLevel(50).
			SetUpdatedAt(time.Now()).
			SetProduct(product).
			SetVendor(vendor).
			Save(ctx)
		require.NoError(t, err)

		// Step 7: Create an address
		address, err := client.Address.Create().
			SetID(uuid.New()).
			SetStreet("123 Main St").
			SetCity("Test City").
			SetState("Test State").
			SetZipCode("12345").
			SetCountry("Test Country").
			Save(ctx)
		require.NoError(t, err)

		// Step 8: Create an order via API
		orderData := map[string]interface{}{
			"id":          uuid.New().String(),
			"order_no":    "ORD-001",
			"status":      "pending",
			"subtotal":    99.99,
			"shipping_fee": 10.0,
			"discount":    0.0,
			"total":       109.99,
			"user_id":     user.ID.String(),
			"shipping_address_id": address.ID.String(),
			"billing_address_id":  address.ID.String(),
		}

		jsonBody, err := json.Marshal(orderData)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var orderResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orderResponse)
		require.NoError(t, err)
		assert.Equal(t, "Order Created Successfully", orderResponse["message"])

		// Step 9: Create order item
		orderItemData := map[string]interface{}{
			"id":         uuid.New().String(),
			"order_id":   orderData["id"],
			"product_id": product.ID.String(),
			"quantity":   2,
			"price":      99.99,
			"subtotal":   199.98,
		}

		jsonBody, err = json.Marshal(orderItemData)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/order_items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Step 10: Create payment
		paymentData := map[string]interface{}{
			"id":       uuid.New().String(),
			"provider": "stripe",
			"status":   "pending",
			"amount":   109.99,
			"order_id": orderData["id"],
		}

		jsonBody, err = json.Marshal(paymentData)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/payments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Step 11: Create delivery
		deliveryData := map[string]interface{}{
			"id":       uuid.New().String(),
			"provider": "fedex",
			"status":   "pending",
			"order_id": orderData["id"],
		}

		jsonBody, err = json.Marshal(deliveryData)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/deliveries", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Step 12: Verify order retrieval
		req = httptest.NewRequest(http.MethodGet, "/api/orders/"+orderData["id"].(string), nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var getOrderResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&getOrderResponse)
		require.NoError(t, err)
		assert.Equal(t, "Order Retrieved Successfully", getOrderResponse["message"])
	})
}

// TestFullSystem_ProductManagementFlow tests product creation, update, and deletion
func TestFullSystem_ProductManagementFlow(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	app := fiber.New()
	cfg := config.Config{
		HTTPPort:    ":3000",
		DatabaseURL: "sqlite3://:memory:",
	}
	mockUploads := new(MockUploadsService)
	mockUploads.On("UploadImage", mock.Anything, mock.Anything, mock.Anything).Return("test-image-url", nil)
	mockUploads.On("GetImageURL", mock.Anything, mock.Anything).Return("https://example.com/image.jpg", nil)

	apiGroup := app.Group("/api")
	httpserver.RegisterRoutes(apiGroup, app, client, cfg)

	ctx := context.Background()

	t.Run("Product management flow", func(t *testing.T) {
		// Create vendor
		vendor, err := client.Vendor.Create().
			SetID(uuid.New()).
			SetName("Test Vendor").
			SetContact("vendor@example.com").
			Save(ctx)
		require.NoError(t, err)

		// Create category
		category, err := client.Category.Create().
			SetID(uuid.New()).
			SetName("Test Category").
			SetSlug("test-category").
			Save(ctx)
		require.NoError(t, err)

		// Create product via API
		productData := map[string]interface{}{
			"id":          uuid.New().String(),
			"name":        "Test Product",
			"sku":         "TEST-001",
			"price":       99.99,
			"description": "Test product description",
			"unit_label":  "kg",
			"is_active":   true,
			"vendor_id":   vendor.ID.String(),
		}

		jsonBody, err := json.Marshal(productData)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var productResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&productResponse)
		require.NoError(t, err)
		assert.Equal(t, "Product Created Successfully", productResponse["message"])

		// Update product
		productID := productResponse["data"].(map[string]interface{})["id"].(string)
		updateData := map[string]interface{}{
			"name":  "Updated Product Name",
			"price": 149.99,
		}

		jsonBody, err = json.Marshal(updateData)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPatch, "/api/products/"+productID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Verify update
		req = httptest.NewRequest(http.MethodGet, "/api/products/"+productID, nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var getProductResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&getProductResponse)
		require.NoError(t, err)
		data := getProductResponse["data"].(map[string]interface{})
		assert.Equal(t, "Updated Product Name", data["name"])
	})
}

// TestFullSystem_UserFlow tests user creation, update, and cart operations
func TestFullSystem_UserFlow(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	app := fiber.New()
	cfg := config.Config{
		HTTPPort:    ":3000",
		DatabaseURL: "sqlite3://:memory:",
	}
	mockUploads := new(MockUploadsService)
	apiGroup := app.Group("/api")
	httpserver.RegisterRoutes(apiGroup, app, client, cfg)

	ctx := context.Background()

	t.Run("User and cart flow", func(t *testing.T) {
		// Create user via API
		userData := map[string]interface{}{
			"id":       uuid.New().String(),
			"email":    "cartuser@example.com",
			"password": "password123",
			"name":     "Cart User",
		}

		jsonBody, err := json.Marshal(userData)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Create vendor and product
		vendor, err := client.Vendor.Create().
			SetID(uuid.New()).
			SetName("Test Vendor").
			SetContact("vendor@example.com").
			Save(ctx)
		require.NoError(t, err)

		product, err := client.Product.Create().
			SetID(uuid.New()).
			SetName("Cart Product").
			SetSku("CART-001").
			SetPrice(49.99).
			SetDescription("Product for cart").
			SetUnitLabel("piece").
			SetIsActive(true).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetVendor(vendor).
			Save(ctx)
		require.NoError(t, err)

		// Create cart
		userID := userData["id"].(string)
		cartData := map[string]interface{}{
			"id":      uuid.New().String(),
			"user_id": userID,
		}

		jsonBody, err = json.Marshal(cartData)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/carts", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Add item to cart
		cartResponse := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&cartResponse)
		cartID := cartResponse["data"].(map[string]interface{})["id"].(string)

		cartItemData := map[string]interface{}{
			"id":         uuid.New().String(),
			"cart_id":    cartID,
			"product_id": product.ID.String(),
			"quantity":   3,
			"price":      49.99,
		}

		jsonBody, err = json.Marshal(cartItemData)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/api/cart_items", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Verify cart items
		req = httptest.NewRequest(http.MethodGet, "/api/cart_items", nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestFullSystem_ErrorHandling tests error handling across the system
func TestFullSystem_ErrorHandling(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	app := fiber.New()
	cfg := config.Config{
		HTTPPort:    ":3000",
		DatabaseURL: "sqlite3://:memory:",
	}
	mockUploads := new(MockUploadsService)
	apiGroup := app.Group("/api")
	httpserver.RegisterRoutes(apiGroup, app, client, cfg)

	t.Run("Invalid UUID handling", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products/invalid-uuid", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "invalid uuid", response["message"])
	})

	t.Run("Non-existent resource", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		req := httptest.NewRequest(http.MethodGet, "/api/products/"+nonExistentID, nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestFullSystem_HealthCheck tests health check endpoint
func TestFullSystem_HealthCheck(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	app := fiber.New()
	cfg := config.Config{
		HTTPPort:    ":3000",
		DatabaseURL: "sqlite3://:memory:",
	}
	mockUploads := new(MockUploadsService)
	apiGroup := app.Group("/api")
	httpserver.RegisterRoutes(apiGroup, app, client, cfg)

	t.Run("Health check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
		assert.Equal(t, "freshease-api", response["service"])
	})
}

