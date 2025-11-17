package orders

import (
	"context"
	"testing"
	"time"

	"freshease/backend/ent/enttest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestEntRepo_List(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test orders
	order1, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	order2, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-002").
		SetStatus("completed").
		SetSubtotal(200.0).
		SetShippingFee(20.0).
		SetDiscount(10.0).
		SetTotal(210.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results
	foundIDs := make(map[uuid.UUID]bool)
	for _, order := range result {
		foundIDs[order.ID] = true
		assert.NotEmpty(t, order.OrderNo)
		assert.Greater(t, order.Total, 0.0)
		assert.Equal(t, user.ID, order.UserID)
	}

	assert.True(t, foundIDs[order1.ID])
	assert.True(t, foundIDs[order2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test order
	createdOrder, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdOrder.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdOrder.ID, result.ID)
	assert.Equal(t, "ORD-001", result.OrderNo)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, 110.0, result.Total)
	assert.Equal(t, user.ID, result.UserID)

	// Test FindByID - not found
	nonExistentID := uuid.New()
	_, err = repo.FindByID(ctx, nonExistentID)
	assert.Error(t, err)
}

func TestEntRepo_Create(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test address
	address, err := client.Address.Create().
		SetID(uuid.New()).
		SetLine1("123 Main St").
		SetCity("Test City").
		SetProvince("Test State").
		SetPostalCode("12345").
		SetCountry("Test Country").
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	now := time.Now()
	dto := &CreateOrderDTO{
		ID:                uuid.New(),
		OrderNo:           "ORD-003",
		Status:            "pending",
		Subtotal:          150.0,
		ShippingFee:       15.0,
		Discount:          5.0,
		Total:             160.0,
		PlacedAt:          &now,
		UserID:            user.ID,
		ShippingAddressID: &address.ID,
		BillingAddressID:  &address.ID,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.OrderNo, result.OrderNo)
	assert.Equal(t, dto.Status, result.Status)
	assert.Equal(t, dto.Total, result.Total)
	assert.Equal(t, user.ID, result.UserID)
	assert.NotNil(t, result.ShippingAddressID)
	assert.NotNil(t, result.BillingAddressID)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test order
	createdOrder, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Update order
	newStatus := "completed"
	newTotal := 120.0
	dto := &UpdateOrderDTO{
		ID:     createdOrder.ID,
		Status: &newStatus,
		Total:  &newTotal,
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdOrder.ID, result.ID)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, 120.0, result.Total)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create test order
	createdOrder, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-001").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	// Delete order
	err = repo.Delete(ctx, createdOrder.ID)
	require.NoError(t, err)

	// Verify order is deleted
	_, err = repo.FindByID(ctx, createdOrder.ID)
	assert.Error(t, err)
}

