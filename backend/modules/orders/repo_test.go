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

	// Create test addresses
	shippingAddr, err := client.Address.Create().
		SetLine1("123 Shipping St").
		SetCity("Shipping City").
		SetProvince("SC").
		SetCountry("USA").
		SetPostalCode("12345").
		SetUser(user).
		Save(ctx)
	require.NoError(t, err)

	billingAddr, err := client.Address.Create().
		SetLine1("456 Billing Ave").
		SetCity("Billing City").
		SetProvince("BC").
		SetCountry("USA").
		SetPostalCode("67890").
		SetUser(user).
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

	// Test Update - basic fields
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

	// Test Update - with shipping and billing addresses
	// Create a new order for this test to ensure clean state
	orderForAddressTest, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-002").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	shippingAddrID := shippingAddr.ID
	billingAddrID := billingAddr.ID
	newStatusForAddress := "processing"
	dtoWithAddresses := &UpdateOrderDTO{
		ID:                orderForAddressTest.ID,
		Status:            &newStatusForAddress, // Include a field update so mutation has fields
		ShippingAddressID: &shippingAddrID,
		BillingAddressID:  &billingAddrID,
	}

	result2, err := repo.Update(ctx, dtoWithAddresses)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, orderForAddressTest.ID, result2.ID)
	assert.Equal(t, newStatusForAddress, result2.Status) // Status was updated
	assert.NotNil(t, result2.ShippingAddressID)
	assert.NotNil(t, result2.BillingAddressID)
	assert.Equal(t, shippingAddrID, *result2.ShippingAddressID)
	assert.Equal(t, billingAddrID, *result2.BillingAddressID)

	// Test Update - no fields to update (should return error)
	// Create a new order for this test to avoid conflicts
	orderForEmptyTest, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-003").
		SetStatus("pending").
		SetSubtotal(100.0).
		SetShippingFee(10.0).
		SetDiscount(0.0).
		SetTotal(110.0).
		AddUser(user).
		Save(ctx)
	require.NoError(t, err)

	emptyDTO := &UpdateOrderDTO{
		ID: orderForEmptyTest.ID,
		// No fields set
	}
	_, err = repo.Update(ctx, emptyDTO)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no fields to update")
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

func TestEntRepo_List_EdgeCases(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Test empty list
	items, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Empty(t, items)

	// Create user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create order without addresses
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

	// Create address
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

	// Create order with addresses
	order2, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-002").
		SetStatus("completed").
		SetSubtotal(200.0).
		SetShippingFee(20.0).
		SetDiscount(10.0).
		SetTotal(210.0).
		AddUser(user).
		AddShippingAddress(address).
		AddBillingAddress(address).
		Save(ctx)
	require.NoError(t, err)

	// Test List with mixed orders
	result, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify order without addresses
	var order1DTO *GetOrderDTO
	for _, o := range result {
		if o.ID == order1.ID {
			order1DTO = o
			break
		}
	}
	require.NotNil(t, order1DTO)
	assert.Nil(t, order1DTO.ShippingAddressID)
	assert.Nil(t, order1DTO.BillingAddressID)

	// Verify order with addresses
	var order2DTO *GetOrderDTO
	for _, o := range result {
		if o.ID == order2.ID {
			order2DTO = o
			break
		}
	}
	require.NotNil(t, order2DTO)
	assert.NotNil(t, order2DTO.ShippingAddressID)
	assert.NotNil(t, order2DTO.BillingAddressID)
	assert.Equal(t, address.ID, *order2DTO.ShippingAddressID)
	assert.Equal(t, address.ID, *order2DTO.BillingAddressID)
}

func TestEntRepo_FindByID_EdgeCases(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Test FindByID - order without addresses
	orderWithoutAddr, err := client.Order.Create().
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

	result, err := repo.FindByID(ctx, orderWithoutAddr.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.ShippingAddressID)
	assert.Nil(t, result.BillingAddressID)

	// Test FindByID - order with addresses
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

	orderWithAddr, err := client.Order.Create().
		SetID(uuid.New()).
		SetOrderNo("ORD-002").
		SetStatus("completed").
		SetSubtotal(200.0).
		SetShippingFee(20.0).
		SetDiscount(10.0).
		SetTotal(210.0).
		AddUser(user).
		AddShippingAddress(address).
		AddBillingAddress(address).
		Save(ctx)
	require.NoError(t, err)

	result2, err := repo.FindByID(ctx, orderWithAddr.ID)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.NotNil(t, result2.ShippingAddressID)
	assert.NotNil(t, result2.BillingAddressID)
	assert.Equal(t, address.ID, *result2.ShippingAddressID)
	assert.Equal(t, address.ID, *result2.BillingAddressID)
}

func TestEntRepo_Create_EdgeCases(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Test Create - without addresses
	dto := &CreateOrderDTO{
		ID:          uuid.New(),
		OrderNo:     "ORD-001",
		Status:      "pending",
		Subtotal:    100.0,
		ShippingFee: 10.0,
		Discount:    0.0,
		Total:       110.0,
		UserID:      user.ID,
		// No addresses
	}

	result, err := repo.Create(ctx, dto)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.ShippingAddressID)
	assert.Nil(t, result.BillingAddressID)

	// Test Create - with only shipping address
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

	dto2 := &CreateOrderDTO{
		ID:                uuid.New(),
		OrderNo:           "ORD-002",
		Status:            "pending",
		Subtotal:          150.0,
		ShippingFee:       15.0,
		Discount:          5.0,
		Total:             160.0,
		UserID:            user.ID,
		ShippingAddressID: &address.ID,
		// No billing address
	}

	result2, err := repo.Create(ctx, dto2)
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.NotNil(t, result2.ShippingAddressID)
	assert.Equal(t, address.ID, *result2.ShippingAddressID)
	assert.Nil(t, result2.BillingAddressID)

	// Test Create - with only billing address
	dto3 := &CreateOrderDTO{
		ID:               uuid.New(),
		OrderNo:          "ORD-003",
		Status:           "pending",
		Subtotal:         150.0,
		ShippingFee:      15.0,
		Discount:         5.0,
		Total:            160.0,
		UserID:           user.ID,
		BillingAddressID: &address.ID,
		// No shipping address
	}

	result3, err := repo.Create(ctx, dto3)
	require.NoError(t, err)
	assert.NotNil(t, result3)
	assert.Nil(t, result3.ShippingAddressID)
	assert.NotNil(t, result3.BillingAddressID)
	assert.Equal(t, address.ID, *result3.BillingAddressID)

	// Test Create - error: invalid user ID
	invalidUserID := uuid.New()
	dto4 := &CreateOrderDTO{
		ID:          uuid.New(),
		OrderNo:     "ORD-004",
		Status:      "pending",
		Subtotal:    100.0,
		ShippingFee: 10.0,
		Discount:    0.0,
		Total:       110.0,
		UserID:      invalidUserID,
	}

	_, err = repo.Create(ctx, dto4)
	assert.Error(t, err)

	// Test Create - error: invalid shipping address ID
	invalidAddrID := uuid.New()
	dto5 := &CreateOrderDTO{
		ID:                uuid.New(),
		OrderNo:           "ORD-005",
		Status:            "pending",
		Subtotal:          100.0,
		ShippingFee:       10.0,
		Discount:          0.0,
		Total:             110.0,
		UserID:            user.ID,
		ShippingAddressID: &invalidAddrID,
	}

	_, err = repo.Create(ctx, dto5)
	assert.Error(t, err)

	// Test Create - error: invalid billing address ID
	dto6 := &CreateOrderDTO{
		ID:               uuid.New(),
		OrderNo:          "ORD-006",
		Status:           "pending",
		Subtotal:         100.0,
		ShippingFee:      10.0,
		Discount:         0.0,
		Total:            110.0,
		UserID:           user.ID,
		BillingAddressID: &invalidAddrID,
	}

	_, err = repo.Create(ctx, dto6)
	assert.Error(t, err)
}

func TestEntRepo_Update_EdgeCases(t *testing.T) {
	client := enttest.Open(t, "sqlite3", ":memory:?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := NewEntRepo(client)
	ctx := context.Background()

	// Create user
	user, err := client.User.Create().
		SetID(uuid.New()).
		SetEmail("test@example.com").
		SetName("Test User").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// Create order
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

	// Test Update - update OrderNo
	newOrderNo := "ORD-UPDATED"
	dto1 := &UpdateOrderDTO{
		ID:      order.ID,
		OrderNo: &newOrderNo,
	}
	result1, err := repo.Update(ctx, dto1)
	require.NoError(t, err)
	assert.Equal(t, newOrderNo, result1.OrderNo)

	// Test Update - update Subtotal
	newSubtotal := 150.0
	dto2 := &UpdateOrderDTO{
		ID:       order.ID,
		Subtotal: &newSubtotal,
	}
	result2, err := repo.Update(ctx, dto2)
	require.NoError(t, err)
	assert.Equal(t, newSubtotal, result2.Subtotal)

	// Test Update - update ShippingFee
	newShippingFee := 20.0
	dto3 := &UpdateOrderDTO{
		ID:          order.ID,
		ShippingFee: &newShippingFee,
	}
	result3, err := repo.Update(ctx, dto3)
	require.NoError(t, err)
	assert.Equal(t, newShippingFee, result3.ShippingFee)

	// Test Update - update Discount
	newDiscount := 15.0
	dto4 := &UpdateOrderDTO{
		ID:       order.ID,
		Discount: &newDiscount,
	}
	result4, err := repo.Update(ctx, dto4)
	require.NoError(t, err)
	assert.Equal(t, newDiscount, result4.Discount)

	// Test Update - update PlacedAt
	now := time.Now()
	dto5 := &UpdateOrderDTO{
		ID:       order.ID,
		PlacedAt: &now,
	}
	result5, err := repo.Update(ctx, dto5)
	require.NoError(t, err)
	assert.NotNil(t, result5.PlacedAt)

	// Test Update - error: invalid shipping address ID
	invalidAddrID := uuid.New()
	newStatus := "processing"
	dto6 := &UpdateOrderDTO{
		ID:                order.ID,
		Status:            &newStatus,
		ShippingAddressID: &invalidAddrID,
	}
	_, err = repo.Update(ctx, dto6)
	assert.Error(t, err)

	// Test Update - error: invalid billing address ID
	dto7 := &UpdateOrderDTO{
		ID:               order.ID,
		Status:           &newStatus,
		BillingAddressID: &invalidAddrID,
	}
	_, err = repo.Update(ctx, dto7)
	assert.Error(t, err)
}

