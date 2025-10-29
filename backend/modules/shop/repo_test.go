package shop

import (
	"context"
	"testing"
	"time"

	"freshease/backend/ent/enttest"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntRepo_GetActiveProducts(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	category := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits").
		SetSlug("fruits").
		SaveX(context.Background())

	vendor := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetEmail("vendor@test.com").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	// Create inventory first
	inventory := client.Inventory.Create().
		SetID(uuid.New()).
		SetQuantity(100).
		SetRestockAmount(50).
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	// Create product with inventory
	product := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetPrice(2.99).
		SetDescription("Fresh red apple").
		SetImageURL("https://example.com/apple.jpg").
		SetUnitLabel("kg").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		AddCatagory(category).
		AddVendor(vendor).
		SetInventory(inventory).
		SaveX(context.Background())

	repo := NewEntRepo(client)

	tests := []struct {
		name     string
		filters  ShopSearchFilters
		expected int
	}{
		{
			name: "get all active products",
			filters: ShopSearchFilters{
				Limit:  10,
				Offset: 0,
			},
			expected: 1,
		},
		{
			name: "filter by category",
			filters: ShopSearchFilters{
				CategoryID: &category.ID,
				Limit:      10,
				Offset:     0,
			},
			expected: 1,
		},
		{
			name: "filter by vendor",
			filters: ShopSearchFilters{
				VendorID: &vendor.ID,
				Limit:    10,
				Offset:   0,
			},
			expected: 1,
		},
		{
			name: "filter by price range",
			filters: ShopSearchFilters{
				MinPrice: func() *float64 { p := 1.0; return &p }(),
				MaxPrice: func() *float64 { p := 5.0; return &p }(),
				Limit:    10,
				Offset:   0,
			},
			expected: 1,
		},
		{
			name: "search by name",
			filters: ShopSearchFilters{
				SearchTerm: func() *string { s := "apple"; return &s }(),
				Limit:      10,
				Offset:     0,
			},
			expected: 1,
		},
		{
			name: "filter in stock only",
			filters: ShopSearchFilters{
				InStock: func() *bool { b := true; return &b }(),
				Limit:   10,
				Offset:  0,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, total, err := repo.GetActiveProducts(context.Background(), tt.filters)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, products, tt.expected)

			if len(products) > 0 {
				p := products[0]
				assert.Equal(t, product.ID, p.ID)
				assert.Equal(t, product.Name, p.Name)
				assert.Equal(t, product.Price, p.Price)
				assert.Equal(t, vendor.ID, p.VendorID)
				assert.Equal(t, *vendor.Name, p.VendorName)
				assert.Equal(t, category.ID, p.CategoryID)
				assert.Equal(t, category.Name, p.CategoryName)
				assert.Equal(t, inventory.Quantity, p.StockQuantity)
				assert.True(t, p.IsInStock)
			}
		})
	}
}

func TestEntRepo_GetProductByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	category := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits").
		SetSlug("fruits").
		SaveX(context.Background())

	vendor := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetEmail("vendor@test.com").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	// Create inventory first
	inventory := client.Inventory.Create().
		SetID(uuid.New()).
		SetQuantity(100).
		SetRestockAmount(50).
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	// Create product with inventory
	product := client.Product.Create().
		SetID(uuid.New()).
		SetName("Apple").
		SetPrice(2.99).
		SetDescription("Fresh red apple").
		SetImageURL("https://example.com/apple.jpg").
		SetUnitLabel("kg").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		AddCatagory(category).
		AddVendor(vendor).
		SetInventory(inventory).
		SaveX(context.Background())

	repo := NewEntRepo(client)

	t.Run("get existing product", func(t *testing.T) {
		result, err := repo.GetProductByID(context.Background(), product.ID)
		require.NoError(t, err)
		assert.Equal(t, product.ID, result.ID)
		assert.Equal(t, product.Name, result.Name)
		assert.Equal(t, product.Price, result.Price)
		assert.Equal(t, vendor.ID, result.VendorID)
		assert.Equal(t, *vendor.Name, result.VendorName)
		assert.Equal(t, category.ID, result.CategoryID)
		assert.Equal(t, category.Name, result.CategoryName)
		assert.Equal(t, inventory.Quantity, result.StockQuantity)
		assert.True(t, result.IsInStock)
	})

	t.Run("get non-existing product", func(t *testing.T) {
		nonExistentID := uuid.New()
		result, err := repo.GetProductByID(context.Background(), nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestEntRepo_GetActiveCategories(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits").
		SetSlug("fruits").
		SaveX(context.Background())

	client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Vegetables").
		SetDescription("Fresh vegetables").
		SetSlug("vegetables").
		SaveX(context.Background())

	repo := NewEntRepo(client)

	t.Run("get all categories", func(t *testing.T) {
		categories, err := repo.GetActiveCategories(context.Background())
		require.NoError(t, err)
		assert.Len(t, categories, 2)

		// Check that categories are sorted by name
		assert.Equal(t, "Fruits", categories[0].Name)
		assert.Equal(t, "Vegetables", categories[1].Name)
	})
}

func TestEntRepo_GetActiveVendors(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Vendor A").
		SetEmail("vendorA@test.com").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Vendor B").
		SetEmail("vendorB@test.com").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	// Create inactive vendor (should not be returned)
	client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Inactive Vendor").
		SetEmail("inactive@test.com").
		SetIsActive("inactive").
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	repo := NewEntRepo(client)

	t.Run("get active vendors only", func(t *testing.T) {
		vendors, err := repo.GetActiveVendors(context.Background())
		require.NoError(t, err)
		assert.Len(t, vendors, 2)

		// Check that vendors are sorted by name
		assert.Equal(t, "Vendor A", vendors[0].Name)
		assert.Equal(t, "Vendor B", vendors[1].Name)
	})
}

func TestEntRepo_GetCategoryByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	category := client.Product_category.Create().
		SetID(uuid.New()).
		SetName("Fruits").
		SetDescription("Fresh fruits").
		SetSlug("fruits").
		SaveX(context.Background())

	repo := NewEntRepo(client)

	t.Run("get existing category", func(t *testing.T) {
		result, err := repo.GetCategoryByID(context.Background(), category.ID)
		require.NoError(t, err)
		assert.Equal(t, category.ID, result.ID)
		assert.Equal(t, category.Name, result.Name)
		assert.Equal(t, category.Description, result.Description)
	})

	t.Run("get non-existing category", func(t *testing.T) {
		nonExistentID := uuid.New()
		result, err := repo.GetCategoryByID(context.Background(), nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestEntRepo_GetVendorByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// Create test data
	vendor := client.Vendor.Create().
		SetID(uuid.New()).
		SetName("Test Vendor").
		SetEmail("vendor@test.com").
		SetIsActive("active").
		SetUpdatedAt(time.Now()).
		SaveX(context.Background())

	repo := NewEntRepo(client)

	t.Run("get existing active vendor", func(t *testing.T) {
		result, err := repo.GetVendorByID(context.Background(), vendor.ID)
		require.NoError(t, err)
		assert.Equal(t, vendor.ID, result.ID)
		assert.Equal(t, *vendor.Name, result.Name)
		assert.Equal(t, *vendor.Email, result.Email)
		assert.Equal(t, vendor.IsActive, result.IsActive)
	})

	t.Run("get non-existing vendor", func(t *testing.T) {
		nonExistentID := uuid.New()
		result, err := repo.GetVendorByID(context.Background(), nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
