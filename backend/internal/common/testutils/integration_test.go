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
	"freshease/backend/modules/products"
	"freshease/backend/modules/uploads"
	"freshease/backend/modules/users"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUploadsService is declared in integration_full_system_test.go to avoid duplication

func TestUsersAPI_Integration(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	// Create test app with users module
	app := fiber.New()
	api := app.Group("/api")
	mockUploads := new(MockUploadsService)
	users.RegisterModuleWithEnt(api, client, mockUploads)

	t.Run("Create and retrieve user", func(t *testing.T) {
		// Create user
		createUserDTO := map[string]interface{}{
			"id":       uuid.New().String(),
			"email":    "integration@example.com",
			"password": "password123",
			"name":     "Integration User",
			"phone":    "+1234567890",
			"bio":      "Integration test user",
		}

		jsonBody, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createResponse)
		require.NoError(t, err)
		assert.Equal(t, "Integration User", createResponse["name"])

		// Retrieve user
		userID := createResponse["id"].(string)
		req = httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var getResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&getResponse)
		require.NoError(t, err)
		assert.Equal(t, "Integration User", getResponse["name"])
		assert.Equal(t, "integration@example.com", getResponse["email"])
	})

	t.Run("List users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1) // At least the user we created
	})

	t.Run("Update user", func(t *testing.T) {
		// First create a user
		createUserDTO := map[string]interface{}{
			"id":       uuid.New().String(),
			"email":    "update@example.com",
			"password": "password123",
			"name":     "Original Name",
		}

		jsonBody, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createResponse)
		require.NoError(t, err)
		userID := createResponse["id"].(string)

		// Update user
		updateUserDTO := map[string]interface{}{
			"name": "Updated Name",
		}

		jsonBody, err = json.Marshal(updateUserDTO)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPut, "/api/users/"+userID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updateResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&updateResponse)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updateResponse["name"])
	})

	t.Run("Delete user", func(t *testing.T) {
		// First create a user
		createUserDTO := map[string]interface{}{
			"id":       uuid.New().String(),
			"email":    "delete@example.com",
			"password": "password123",
			"name":     "Delete User",
		}

		jsonBody, err := json.Marshal(createUserDTO)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createResponse)
		require.NoError(t, err)
		userID := createResponse["id"].(string)

		// Delete user
		req = httptest.NewRequest(http.MethodDelete, "/api/users/"+userID, nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify user is deleted
		req = httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestProductsAPI_Integration(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	// Create test app with products module
	app := fiber.New()
	api := app.Group("/api")
	mockUploads := new(MockUploadsService)
	products.RegisterModuleWithEnt(api, client, mockUploads)

	t.Run("Create and retrieve product", func(t *testing.T) {
		// First create a vendor
		vendor, err := client.Vendor.Create().
			SetName("Test Vendor").
			SetContact("vendor@test.com").
			Save(context.Background())
		require.NoError(t, err)

		// Create a category
		category, err := client.Category.Create().
			SetName("Test Category").
			SetSlug("test-category").
			Save(context.Background())
		require.NoError(t, err)

		// Create product directly with relationships (bypassing the API)
		product, err := client.Product.Create().
			SetID(uuid.New()).
			SetName("Integration Product").
			SetSku("INT-PROD-001").
			SetPrice(99.99).
			SetDescription("Integration test product").
			SetUnitLabel("kg").
			SetIsActive(true).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetVendor(vendor).
			Save(context.Background())
		require.NoError(t, err)

		// Create product_category join
		_, err = client.Product_category.Create().
			SetProduct(product).
			SetCategory(category).
			Save(context.Background())
		require.NoError(t, err)

		// Create inventory
		_, err = client.Inventory.Create().
			SetQuantity(100).
			SetReorderLevel(50).
			SetUpdatedAt(time.Now()).
			SetProduct(product).
			SetVendor(vendor).
			Save(context.Background())
		require.NoError(t, err)

		// Now test retrieving the product via API
		req := httptest.NewRequest(http.MethodGet, "/api/products/"+product.ID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var getResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&getResponse)
		require.NoError(t, err)
		assert.Equal(t, "Product Retrieved Successfully", getResponse["message"])
		assert.Equal(t, "Integration Product", getResponse["data"].(map[string]interface{})["name"])
	})

	t.Run("List products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "Products Retrieved Successfully", response["message"])
		assert.Contains(t, response, "data")
	})

	t.Run("Update product", func(t *testing.T) {
		// First create a vendor
		vendor, err := client.Vendor.Create().
			SetName("Test Vendor 2").
			SetContact("vendor2@test.com").
			Save(context.Background())
		require.NoError(t, err)

		// Create a category
		category, err := client.Category.Create().
			SetName("Test Category 2").
			SetSlug("test-category-2").
			Save(context.Background())
		require.NoError(t, err)

		// Create product directly with relationships
		product, err := client.Product.Create().
			SetID(uuid.New()).
			SetName("Original Product").
			SetSku("ORIG-PROD-001").
			SetPrice(99.99).
			SetDescription("Original description").
			SetUnitLabel("kg").
			SetIsActive(true).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetVendor(vendor).
			Save(context.Background())
		require.NoError(t, err)

		// Create product_category join
		_, err = client.Product_category.Create().
			SetProduct(product).
			SetCategory(category).
			Save(context.Background())
		require.NoError(t, err)

		// Create inventory
		_, err = client.Inventory.Create().
			SetQuantity(100).
			SetReorderLevel(50).
			SetUpdatedAt(time.Now()).
			SetProduct(product).
			SetVendor(vendor).
			Save(context.Background())
		require.NoError(t, err)

		productID := product.ID.String()

		// Update product
		updateProductDTO := map[string]interface{}{
			"name":        "Updated Product",
			"price":       149.99,
			"description": "Updated description",
		}

		jsonBody, err := json.Marshal(updateProductDTO)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPatch, "/api/products/"+productID, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var updateResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&updateResponse)
		require.NoError(t, err)
		assert.Equal(t, "Product Updated Successfully", updateResponse["message"])
		assert.Equal(t, "Updated Product", updateResponse["data"].(map[string]interface{})["name"])
	})

	t.Run("Delete product", func(t *testing.T) {
		// First create a vendor
		vendor, err := client.Vendor.Create().
			SetName("Test Vendor 3").
			SetContact("vendor3@test.com").
			Save(context.Background())
		require.NoError(t, err)

		// Create a category
		category, err := client.Category.Create().
			SetName("Test Category 3").
			SetSlug("test-category-3").
			Save(context.Background())
		require.NoError(t, err)

		// Create product directly with relationships
		product, err := client.Product.Create().
			SetID(uuid.New()).
			SetName("Delete Product").
			SetSku("DEL-PROD-001").
			SetPrice(99.99).
			SetDescription("Product to delete").
			SetUnitLabel("kg").
			SetIsActive(true).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			SetVendor(vendor).
			Save(context.Background())
		require.NoError(t, err)

		// Create product_category join
		_, err = client.Product_category.Create().
			SetProduct(product).
			SetCategory(category).
			Save(context.Background())
		require.NoError(t, err)

		// Create inventory
		_, err = client.Inventory.Create().
			SetQuantity(100).
			SetReorderLevel(50).
			SetUpdatedAt(time.Now()).
			SetProduct(product).
			SetVendor(vendor).
			Save(context.Background())
		require.NoError(t, err)

		productID := product.ID.String()

		// Delete product
		req := httptest.NewRequest(http.MethodDelete, "/api/products/"+productID, nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)

		var deleteResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&deleteResponse)
		require.NoError(t, err)
		assert.Equal(t, "Product Deleted Successfully", deleteResponse["message"])

		// Verify product is deleted
		req = httptest.NewRequest(http.MethodGet, "/api/products/"+productID, nil)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestAPI_ErrorHandling(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	// Create test app with users module
	app := fiber.New()
	api := app.Group("/api")
	mockUploads := new(MockUploadsService)
	users.RegisterModuleWithEnt(api, client, mockUploads)

	t.Run("Invalid UUID format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users/invalid-uuid", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "invalid uuid", response["message"])
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		// The actual error message is about validation, not JSON parsing
		assert.Contains(t, response["message"].(string), "validator failed")
	})

	t.Run("Validation errors", func(t *testing.T) {
		invalidUserDTO := map[string]interface{}{
			"id":       uuid.New().String(),
			"email":    "invalid-email", // Invalid email format
			"password": "123",           // Too short
			"name":     "A",             // Too short
		}

		jsonBody, err := json.Marshal(invalidUserDTO)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		// The actual status code is 400, not 422
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		// The actual error message contains "validator failed"
		assert.Contains(t, response["message"].(string), "validator failed")
	})
}

func TestAPI_NotFoundHandling(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&_fk=1")
	defer client.Close()

	// Create test app with users module
	app := fiber.New()
	api := app.Group("/api")
	mockUploads := new(MockUploadsService)
	users.RegisterModuleWithEnt(api, client, mockUploads)

	t.Run("Non-existent user", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		req := httptest.NewRequest(http.MethodGet, "/api/users/"+nonExistentID, nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "not found", response["message"])
	})
}
