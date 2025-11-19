package carts

import (
	"context"
	"errors"

	"freshease/backend/ent"
	"freshease/backend/ent/cart"
	"freshease/backend/ent/cart_item"
	"freshease/backend/ent/product"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]*GetCartDTO, error)
	Get(ctx context.Context, id uuid.UUID) (*GetCartDTO, error)
	Create(ctx context.Context, dto CreateCartDTO) (*GetCartDTO, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateCartDTO) (*GetCartDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// New methods for cart operations
	GetCurrentCart(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error)
	AddItemToCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*GetCartDTO, error)
	UpdateCartItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID, quantity int) (*GetCartDTO, error)
	RemoveCartItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID) (*GetCartDTO, error)
	ApplyPromoCode(ctx context.Context, userID uuid.UUID, promoCode string) (*GetCartDTO, error)
	RemovePromoCode(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error)
	ClearCart(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error)
}

type service struct {
	repo     Repository
	entClient *ent.Client
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func NewServiceWithClient(r Repository, client *ent.Client) Service {
	return &service{repo: r, entClient: client}
}

func (s *service) List(ctx context.Context) ([]*GetCartDTO, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*GetCartDTO, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, dto CreateCartDTO) (*GetCartDTO, error) {
	return s.repo.Create(ctx, &dto)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, dto UpdateCartDTO) (*GetCartDTO, error) {
	dto.ID = id
	return s.repo.Update(ctx, &dto)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetCurrentCart(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error) {
	cart, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.calculateCartTotals(cart), nil
}

func (s *service) AddItemToCart(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	// Get or create cart
	cartDTO, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get product
	prod, err := s.entClient.Product.Get(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check if item already exists in cart
	existingItem, err := s.entClient.Cart_item.Query().
		Where(
			cart_item.HasCartWith(cart.ID(cartDTO.ID)),
			cart_item.HasProductWith(product.ID(productID)),
		).
		Only(ctx)

	if err == nil {
		// Update existing item
		newQty := existingItem.Qty + quantity
		newLineTotal := prod.Price * float64(newQty)
		
		_, err = s.entClient.Cart_item.UpdateOneID(existingItem.ID).
			SetQty(newQty).
			SetUnitPrice(prod.Price).
			SetLineTotal(newLineTotal).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// Create new item
		cartEntity, err := s.entClient.Cart.Get(ctx, cartDTO.ID)
		if err != nil {
			return nil, err
		}

		lineTotal := prod.Price * float64(quantity)
		_, err = s.entClient.Cart_item.Create().
			SetID(uuid.New()).
			SetQty(quantity).
			SetUnitPrice(prod.Price).
			SetLineTotal(lineTotal).
			SetCart(cartEntity).
			SetProduct(prod).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Recalculate cart totals
	return s.recalculateCart(ctx, cartDTO.ID)
}

func (s *service) UpdateCartItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID, quantity int) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	// Verify cart item belongs to user's cart
	cartDTO, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	item, err := s.entClient.Cart_item.Query().
		Where(
			cart_item.ID(cartItemID),
			cart_item.HasCartWith(cart.ID(cartDTO.ID)),
		).
		WithProduct().
		Only(ctx)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	if quantity <= 0 {
		// Remove item
		return s.RemoveCartItem(ctx, userID, cartItemID)
	}

	// Update quantity
	lineTotal := item.UnitPrice * float64(quantity)
	_, err = s.entClient.Cart_item.UpdateOneID(cartItemID).
		SetQty(quantity).
		SetLineTotal(lineTotal).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.recalculateCart(ctx, cartDTO.ID)
}

func (s *service) RemoveCartItem(ctx context.Context, userID uuid.UUID, cartItemID uuid.UUID) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	// Verify cart item belongs to user's cart
	cartDTO, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.entClient.Cart_item.Query().
		Where(
			cart_item.ID(cartItemID),
			cart_item.HasCartWith(cart.ID(cartDTO.ID)),
		).
		Only(ctx)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	err = s.entClient.Cart_item.DeleteOneID(cartItemID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return s.recalculateCart(ctx, cartDTO.ID)
}

func (s *service) ApplyPromoCode(ctx context.Context, userID uuid.UUID, promoCode string) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	// For now, we'll store promo code in a simple way
	// In a real system, you'd have a promo codes table
	cart, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate discount based on promo code
	discount := 0.0
	if promoCode == "FRESH10" {
		discount = cart.Subtotal * 0.10
	} else if promoCode == "FREESHIP" {
		// Free shipping - discount equals shipping cost
		shipping := s.calculateShipping(cart.Subtotal)
		discount = shipping
	}

	cart.PromoCode = &promoCode
	cart.PromoDiscount = discount
	cart.Discount = discount

	// Update cart discount
	_, err = s.entClient.Cart.UpdateOneID(cart.ID).
		SetDiscount(discount).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.calculateCartTotals(cart), nil
}

func (s *service) RemovePromoCode(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	cart, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	cart.PromoCode = nil
	cart.PromoDiscount = 0.0
	cart.Discount = 0.0

	// Update cart discount
	_, err = s.entClient.Cart.UpdateOneID(cart.ID).
		SetDiscount(0.0).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.calculateCartTotals(cart), nil
}

func (s *service) ClearCart(ctx context.Context, userID uuid.UUID) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	cartDTO, err := s.repo.GetOrCreateCartForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Delete all cart items
	_, err = s.entClient.Cart_item.Delete().
		Where(cart_item.HasCartWith(cart.ID(cartDTO.ID))).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Reset cart totals
	_, err = s.entClient.Cart.UpdateOneID(cartDTO.ID).
		SetSubtotal(0.0).
		SetDiscount(0.0).
		SetTotal(0.0).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.recalculateCart(ctx, cartDTO.ID)
}

// Helper functions
func (s *service) recalculateCart(ctx context.Context, cartID uuid.UUID) (*GetCartDTO, error) {
	if s.entClient == nil {
		return nil, errors.New("ent client not initialized")
	}

	// Reload cart with items
	cartEntity, err := s.entClient.Cart.Query().
		Where(cart.ID(cartID)).
		WithItems(func(q *ent.CartItemQuery) {
			q.WithProduct()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate subtotal from items
	subtotal := 0.0
	for _, item := range cartEntity.Edges.Items {
		subtotal += item.LineTotal
	}

	// Update cart subtotal
	_, err = s.entClient.Cart.UpdateOneID(cartID).
		SetSubtotal(subtotal).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	repo := s.repo.(*EntRepo)
	cart := repo.cartToDTO(cartEntity)

	return s.calculateCartTotals(cart), nil
}

func (s *service) calculateCartTotals(cart *GetCartDTO) *GetCartDTO {
	// Calculate subtotal from items if not already set
	subtotal := cart.Subtotal
	if len(cart.Items) > 0 {
		subtotal = 0.0
		for _, item := range cart.Items {
			subtotal += item.LineTotal
		}
	}

	// Calculate shipping (free if subtotal >= 200)
	shipping := s.calculateShipping(subtotal)

	// Calculate tax (7% VAT on subtotal after discount)
	tax := (subtotal - cart.PromoDiscount) * 0.07

	// Calculate total
	total := (subtotal - cart.PromoDiscount) + shipping + tax

	cart.Subtotal = subtotal
	cart.Shipping = shipping
	cart.Tax = tax
	cart.Total = total

	return cart
}

func (s *service) calculateShipping(subtotal float64) float64 {
	if subtotal >= 200 {
		return 0.0
	}
	return 20.0
}
