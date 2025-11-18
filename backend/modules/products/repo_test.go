package products

import (
	"context"
	"testing"
	"time"

	"freshease/backend/ent/enttest"
	"freshease/backend/internal/common/errs"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Test empty list
	products, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, products)

	// Create required entities first
	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	category, err := client.Category.Create().
		SetName("Test Category").
		SetSlug("test-category").
		Save(ctx)
	require.NoError(t, err)

	// Create test products first (inventories need products)
	product1 := &CreateProductDTO{
		ID:          uuid.New(),
		Name:        "Apple",
		SKU:         "APPLE-001",
		Price:       1.99,
		Description: stringPtr("Fresh red apples"),
		UnitLabel:   "kg",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	product2 := &CreateProductDTO{
		ID:          uuid.New(),
		Name:        "Banana",
		SKU:         "BANANA-001",
		Price:       0.99,
		Description: stringPtr("Yellow bananas"),
		UnitLabel:   "kg",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create products with relationships
	prod1, err := client.Product.Create().
		SetID(product1.ID).
		SetName(product1.Name).
		SetSku(product1.SKU).
		SetPrice(product1.Price).
		SetUnitLabel(product1.UnitLabel).
		SetIsActive(product1.IsActive).
		SetCreatedAt(product1.CreatedAt).
		SetUpdatedAt(product1.UpdatedAt).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	if product1.Description != nil {
		prod1, err = prod1.Update().SetDescription(*product1.Description).Save(ctx)
		require.NoError(t, err)
	}

	// Create inventory1 with product and vendor
	_, err = client.Inventory.Create().
		SetQuantity(100).
		SetReorderLevel(50).
		SetProduct(prod1).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create product_category join
	_, err = client.Product_category.Create().
		SetProduct(prod1).
		SetCategory(category).
		Save(ctx)
	require.NoError(t, err)

	prod2, err := client.Product.Create().
		SetID(product2.ID).
		SetName(product2.Name).
		SetSku(product2.SKU).
		SetPrice(product2.Price).
		SetUnitLabel(product2.UnitLabel).
		SetIsActive(product2.IsActive).
		SetCreatedAt(product2.CreatedAt).
		SetUpdatedAt(product2.UpdatedAt).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	if product2.Description != nil {
		prod2, err = prod2.Update().SetDescription(*product2.Description).Save(ctx)
		require.NoError(t, err)
	}

	// Create inventory2 with product and vendor
	_, err = client.Inventory.Create().
		SetQuantity(200).
		SetReorderLevel(100).
		SetProduct(prod2).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create product_category join
	_, err = client.Product_category.Create().
		SetProduct(prod2).
		SetCategory(category).
		Save(ctx)
	require.NoError(t, err)

	// Test populated list
	products, err = repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, products, 2)

	// Verify products are returned
	productMap := make(map[uuid.UUID]*GetProductDTO)
	for _, product := range products {
		productMap[product.ID] = product
	}

	assert.Contains(t, productMap, product1.ID)
	assert.Contains(t, productMap, product2.ID)
}

func TestRepository_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Test product not found
	nonExistentID := uuid.New()
	_, err := repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)

	// Create required entities first
	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	category, err := client.Category.Create().
		SetName("Test Category").
		SetSlug("test-category").
		Save(ctx)
	require.NoError(t, err)

	// Create test product
	createDTO := &CreateProductDTO{
		ID:          uuid.New(),
		Name:        "Orange",
		SKU:         "ORANGE-001",
		Price:       2.49,
		Description: stringPtr("Fresh oranges"),
		UnitLabel:   "kg",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create product with relationships
	prod, err := client.Product.Create().
		SetID(createDTO.ID).
		SetName(createDTO.Name).
		SetSku(createDTO.SKU).
		SetPrice(createDTO.Price).
		SetUnitLabel(createDTO.UnitLabel).
		SetIsActive(createDTO.IsActive).
		SetCreatedAt(createDTO.CreatedAt).
		SetUpdatedAt(createDTO.UpdatedAt).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	if createDTO.Description != nil {
		prod, err = prod.Update().SetDescription(*createDTO.Description).Save(ctx)
		require.NoError(t, err)
	}

	// Create inventory with product and vendor
	_, err = client.Inventory.Create().
		SetQuantity(100).
		SetReorderLevel(50).
		SetProduct(prod).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Create product_category join
	_, err = client.Product_category.Create().
		SetProduct(prod).
		SetCategory(category).
		Save(ctx)
	require.NoError(t, err)

	// Test finding existing product
	foundProduct, err := repo.FindByID(ctx, createDTO.ID)
	require.NoError(t, err)
	assert.Equal(t, createDTO.ID, foundProduct.ID)
	assert.Equal(t, createDTO.Name, foundProduct.Name)
	assert.Equal(t, createDTO.SKU, foundProduct.SKU)
	assert.Equal(t, createDTO.Price, foundProduct.Price)
	if createDTO.Description != nil {
		assert.Equal(t, *createDTO.Description, *foundProduct.Description)
	}
	assert.Equal(t, createDTO.UnitLabel, foundProduct.UnitLabel)
	assert.Equal(t, createDTO.IsActive, foundProduct.IsActive)
}

func TestRepository_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test product DTO
	createDTO := &CreateProductDTO{
		ID:          uuid.New(),
		Name:        "Test Product",
		SKU:         "TEST-001",
		Price:       19.99,
		Description: stringPtr("Test product description"),
		UnitLabel:   "kg",
		ImageURL:    stringPtr("images/test.jpg"),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Quantity:    100,
		ReorderLevel: 50,
	}

	// Test creating product
	createdProduct, err := repo.Create(ctx, createDTO)
	require.NoError(t, err)
	assert.Equal(t, createDTO.ID, createdProduct.ID)
	assert.Equal(t, createDTO.Name, createdProduct.Name)
	assert.Equal(t, createDTO.SKU, createdProduct.SKU)
	assert.Equal(t, createDTO.Price, createdProduct.Price)
	assert.Equal(t, createDTO.UnitLabel, createdProduct.UnitLabel)
	assert.Equal(t, createDTO.IsActive, createdProduct.IsActive)
	if createDTO.Description != nil {
		assert.Equal(t, *createDTO.Description, *createdProduct.Description)
	}
	if createDTO.ImageURL != nil {
		assert.Equal(t, *createDTO.ImageURL, *createdProduct.ImageURL)
	}
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	// Create initial product
	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Original Product").
		SetSku("ORIG-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Update product
	newName := "Updated Product"
	newPrice := 15.99
	newDescription := "Updated description"
	updateDTO := &UpdateProductDTO{
		ID:          product.ID,
		Name:        &newName,
		Price:       &newPrice,
		Description: &newDescription,
	}

	updatedProduct, err := repo.Update(ctx, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, product.ID, updatedProduct.ID)
	assert.Equal(t, newName, updatedProduct.Name)
	assert.Equal(t, newPrice, updatedProduct.Price)
	assert.Equal(t, newDescription, *updatedProduct.Description)
	assert.Equal(t, product.Sku, updatedProduct.SKU) // SKU not updated

	// Test Update - SKU
	newSKU := "UPD-001"
	updateDTO2 := &UpdateProductDTO{
		ID:  product.ID,
		SKU: &newSKU,
	}
	updatedProduct2, err := repo.Update(ctx, updateDTO2)
	require.NoError(t, err)
	assert.Equal(t, newSKU, updatedProduct2.SKU)

	// Test Update - UnitLabel
	newUnitLabel := "lb"
	updateDTO3 := &UpdateProductDTO{
		ID:        product.ID,
		UnitLabel: &newUnitLabel,
	}
	updatedProduct3, err := repo.Update(ctx, updateDTO3)
	require.NoError(t, err)
	assert.Equal(t, newUnitLabel, updatedProduct3.UnitLabel)

	// Test Update - ImageURL
	newImageURL := "https://example.com/new-image.jpg"
	updateDTO4 := &UpdateProductDTO{
		ID:       product.ID,
		ImageURL: &newImageURL,
	}
	updatedProduct4, err := repo.Update(ctx, updateDTO4)
	require.NoError(t, err)
	assert.NotNil(t, updatedProduct4.ImageURL)
	assert.Equal(t, newImageURL, *updatedProduct4.ImageURL)

	// Test Update - IsActive
	newIsActive := false
	updateDTO5 := &UpdateProductDTO{
		ID:       product.ID,
		IsActive: &newIsActive,
	}
	updatedProduct5, err := repo.Update(ctx, updateDTO5)
	require.NoError(t, err)
	assert.Equal(t, newIsActive, updatedProduct5.IsActive)

	// Test Update - clear ImageURL (set to empty string)
	emptyImageURL := ""
	updateDTO6 := &UpdateProductDTO{
		ID:       product.ID,
		ImageURL: &emptyImageURL,
	}
	updatedProduct6, err := repo.Update(ctx, updateDTO6)
	require.NoError(t, err)
	if updatedProduct6.ImageURL != nil {
		assert.Equal(t, emptyImageURL, *updatedProduct6.ImageURL)
	}

	// Test Update - no fields to update
	noUpdateDTO := &UpdateProductDTO{ID: product.ID}
	_, err = repo.Update(ctx, noUpdateDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errs.NoFieldsToUpdate.Error())
}

func TestRepository_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create required entities
	vendor, err := client.Vendor.Create().
		SetName("Test Vendor").
		SetContact("vendor@example.com").
		Save(ctx)
	require.NoError(t, err)

	// Create product
	product, err := client.Product.Create().
		SetID(uuid.New()).
		SetName("Test Product").
		SetSku("TEST-001").
		SetPrice(10.99).
		SetUnitLabel("kg").
		SetIsActive(true).
		SetVendor(vendor).
		Save(ctx)
	require.NoError(t, err)

	// Test deleting product
	err = repo.Delete(ctx, product.ID)
	require.NoError(t, err)

	// Verify product is deleted
	_, err = repo.FindByID(ctx, product.ID)
	assert.Error(t, err)
}
