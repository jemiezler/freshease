package payments

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

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
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

	// Create test payments
	payment1, err := client.Payment.Create().
		SetID(uuid.New()).
		SetProvider("stripe").
		SetStatus("pending").
		SetAmount(110.0).
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	payment2, err := client.Payment.Create().
		SetID(uuid.New()).
		SetProvider("paypal").
		SetStatus("completed").
		SetAmount(110.0).
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Test List
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify results
	foundIDs := make(map[uuid.UUID]bool)
	for _, payment := range result {
		foundIDs[payment.ID] = true
		assert.NotEmpty(t, payment.Provider)
		assert.NotEmpty(t, payment.Status)
		assert.Greater(t, payment.Amount, 0.0)
		assert.Equal(t, order.ID, payment.OrderID)
	}

	assert.True(t, foundIDs[payment1.ID])
	assert.True(t, foundIDs[payment2.ID])
}

func TestEntRepo_FindByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
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

	// Create test payment
	providerRef := "pay_123456"
	createdPayment, err := client.Payment.Create().
		SetID(uuid.New()).
		SetProvider("stripe").
		SetNillableProviderRef(&providerRef).
		SetStatus("pending").
		SetAmount(110.0).
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - success
	result, err := repo.FindByID(ctx, createdPayment.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdPayment.ID, result.ID)
	assert.Equal(t, "stripe", result.Provider)
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, 110.0, result.Amount)
	assert.NotNil(t, result.ProviderRef)
	assert.Equal(t, providerRef, *result.ProviderRef)
	assert.Equal(t, order.ID, result.OrderID)

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

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
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

	providerRef := "pay_789012"
	paidAt := time.Now()
	dto := &CreatePaymentDTO{
		ID:          uuid.New(),
		Provider:    "stripe",
		ProviderRef: &providerRef,
		Status:      "completed",
		Amount:      110.0,
		PaidAt:      &paidAt,
		OrderID:     order.ID,
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dto.ID, result.ID)
	assert.Equal(t, dto.Provider, result.Provider)
	assert.Equal(t, dto.Status, result.Status)
	assert.Equal(t, dto.Amount, result.Amount)
	assert.NotNil(t, result.ProviderRef)
	assert.Equal(t, providerRef, *result.ProviderRef)
	assert.NotNil(t, result.PaidAt)
	assert.Equal(t, order.ID, result.OrderID)
}

func TestEntRepo_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
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

	// Create test payment
	createdPayment, err := client.Payment.Create().
		SetID(uuid.New()).
		SetProvider("stripe").
		SetStatus("pending").
		SetAmount(110.0).
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Update payment
	newStatus := "completed"
	newProviderRef := "pay_updated"
	paidAt := time.Now()
	dto := &UpdatePaymentDTO{
		ID:          createdPayment.ID,
		Status:      &newStatus,
		ProviderRef: &newProviderRef,
		PaidAt:      &paidAt,
	}

	result, err := repo.Update(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdPayment.ID, result.ID)
	assert.Equal(t, "completed", result.Status)
	assert.NotNil(t, result.ProviderRef)
	assert.Equal(t, newProviderRef, *result.ProviderRef)
	assert.NotNil(t, result.PaidAt)
}

func TestEntRepo_Delete(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create test user and order
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	order, err := client.Order.Create().
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

	// Create test payment
	createdPayment, err := client.Payment.Create().
		SetID(uuid.New()).
		SetProvider("stripe").
		SetStatus("pending").
		SetAmount(110.0).
		AddOrder(order).
		Save(ctx)
	require.NoError(t, err)

	// Delete payment
	err = repo.Delete(ctx, createdPayment.ID)
	require.NoError(t, err)

	// Verify payment is deleted
	_, err = repo.FindByID(ctx, createdPayment.ID)
	assert.Error(t, err)
}

