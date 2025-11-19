package carts

import (
	"freshease/backend/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	r.Get("/",   ctl.ListCarts)
	r.Get("/current", ctl.GetCurrentCart)
	r.Get("/:id", ctl.GetCart)
	r.Post("/",  ctl.CreateCart)
	r.Patch("/:id", ctl.UpdateCart)
	r.Patch("/add-item", ctl.AddItemToCart)
	r.Patch("/update-item", ctl.UpdateCartItem)
	r.Delete("/remove-item/:id", ctl.RemoveCartItem)
	r.Post("/apply-promo", ctl.ApplyPromoCode)
	r.Delete("/remove-promo", ctl.RemovePromoCode)
	r.Delete("/clear", ctl.ClearCart)
	r.Delete("/:id", ctl.DeleteCart)
}

// ListCarts godoc
// @Summary      List carts
// @Description  Get all carts
// @Tags         carts
// @Produce      json
// @Success      200 {array}  GetCartDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /carts [get]
func (ctl *Controller) ListCarts(c *fiber.Ctx) error {
	items, err := ctl.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items, "message": "Carts Retrieved Successfully"})
}

// GetCart godoc
// @Summary      Get cart by ID
// @Tags         carts
// @Produce      json
// @Param        id   path      string true "Cart ID (UUID)"
// @Success      200  {object}  GetCartDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /carts/{id} [get]
func (ctl *Controller) GetCart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	item, err := ctl.svc.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": item, "message": "Cart Retrieved Successfully"})
}

// CreateCart godoc
// @Summary      Create cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        payload body      CreateCartDTO true "Cart payload"
// @Success      201     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts [post]
func (ctl *Controller) CreateCart(c *fiber.Ctx) error {
	var dto CreateCartDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Create(c.Context(), dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Cart Created Successfully"})
}

// UpdateCart godoc
// @Summary      Update cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        id      path      string        true "Cart ID (UUID)"
// @Param        payload body      UpdateCartDTO true "Partial/Full update"
// @Success      201     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts/{id} [patch]
func (ctl *Controller) UpdateCart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	var dto UpdateCartDTO
	if err := middleware.BindAndValidate(c, &dto); err != nil {
		return err
	}
	item, err := ctl.svc.Update(c.Context(), id, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": item, "message": "Cart Updated Successfully"})
}

// DeleteCart godoc
// @Summary      Delete cart
// @Tags         carts
// @Produce      json
// @Param        id   path      string true "Cart ID (UUID)"
// @Success      202  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /carts/{id} [delete]
func (ctl *Controller) DeleteCart(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}
	if err := ctl.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Cart Deleted Successfully"})
}

// GetCurrentCart godoc
// @Summary      Get current user's cart
// @Tags         carts
// @Produce      json
// @Success      200  {object}  GetCartDTO
// @Failure      400  {object}  map[string]interface{}
// @Router       /carts/current [get]
func (ctl *Controller) GetCurrentCart(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	cart, err := ctl.svc.GetCurrentCart(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Cart Retrieved Successfully"})
}

// AddItemToCart godoc
// @Summary      Add item to cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        payload body      AddToCartRequest true "Add to cart request"
// @Success      200     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts/add-item [patch]
func (ctl *Controller) AddItemToCart(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	var req AddToCartRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid product id: " + err.Error()})
	}
	cart, err := ctl.svc.AddItemToCart(c.Context(), userID, productID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Item Added Successfully"})
}

// UpdateCartItem godoc
// @Summary      Update cart item quantity
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        payload body      UpdateCartItemRequest true "Update cart item request"
// @Success      200     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts/update-item [patch]
func (ctl *Controller) UpdateCartItem(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	var req UpdateCartItemRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}
	cartItemID, err := uuid.Parse(req.CartItemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid cart item id: " + err.Error()})
	}
	cart, err := ctl.svc.UpdateCartItem(c.Context(), userID, cartItemID, req.Quantity)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Cart Item Updated Successfully"})
}

// RemoveCartItem godoc
// @Summary      Remove item from cart
// @Tags         carts
// @Produce      json
// @Param        id   path      string true "Cart Item ID (UUID)"
// @Success      200  {object}  GetCartDTO
// @Failure      400  {object}  map[string]interface{}
// @Router       /carts/remove-item/{id} [delete]
func (ctl *Controller) RemoveCartItem(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	idStr := c.Params("id")
	cartItemID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid cart item id"})
	}
	cart, err := ctl.svc.RemoveCartItem(c.Context(), userID, cartItemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Item Removed Successfully"})
}

// ApplyPromoCode godoc
// @Summary      Apply promo code to cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Param        payload body      ApplyPromoRequest true "Promo code request"
// @Success      200     {object}  GetCartDTO
// @Failure      400     {object}  map[string]interface{}
// @Router       /carts/apply-promo [post]
func (ctl *Controller) ApplyPromoCode(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	var req ApplyPromoRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}
	cart, err := ctl.svc.ApplyPromoCode(c.Context(), userID, req.PromoCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Promo Code Applied Successfully"})
}

// RemovePromoCode godoc
// @Summary      Remove promo code from cart
// @Tags         carts
// @Produce      json
// @Success      200  {object}  GetCartDTO
// @Failure      400  {object}  map[string]interface{}
// @Router       /carts/remove-promo [delete]
func (ctl *Controller) RemovePromoCode(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	cart, err := ctl.svc.RemovePromoCode(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Promo Code Removed Successfully"})
}

// ClearCart godoc
// @Summary      Clear all items from cart
// @Tags         carts
// @Produce      json
// @Success      200  {object}  GetCartDTO
// @Failure      400  {object}  map[string]interface{}
// @Router       /carts/clear [delete]
func (ctl *Controller) ClearCart(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not authenticated"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id"})
	}
	cart, err := ctl.svc.ClearCart(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cart, "message": "Cart Cleared Successfully"})
}
