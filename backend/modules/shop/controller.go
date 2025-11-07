package shop

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct{ svc Service }

func NewController(s Service) *Controller { return &Controller{svc: s} }

func (ctl *Controller) Register(r fiber.Router) {
	// Product endpoints
	r.Get("/products", ctl.SearchProducts)
	r.Get("/products/:id", ctl.GetProduct)

	// Category endpoints
	r.Get("/categories", ctl.GetCategories)
	r.Get("/categories/:id", ctl.GetCategory)

	// Vendor endpoints
	r.Get("/vendors", ctl.GetVendors)
	r.Get("/vendors/:id", ctl.GetVendor)
}

// SearchProducts godoc
// @Summary      Search products
// @Description  Public catalog search with optional filters
// @Tags         shop
// @Produce      json
// @Param        category_id  query string false "Filter by category UUID"
// @Param        vendor_id    query string false "Filter by vendor UUID"
// @Param        min_price    query number false "Minimum price"
// @Param        max_price    query number false "Maximum price"
// @Param        search       query string false "Search term"
// @Param        in_stock     query boolean false "Only items in stock"
// @Param        limit        query integer false "Limit results"
// @Param        offset       query integer false "Offset results"
// @Success      200 {array}  products.GetProductDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /shop/products [get]
func (ctl *Controller) SearchProducts(c *fiber.Ctx) error {
	filters := ShopSearchFilters{}

	// Parse query parameters
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filters.CategoryID = &categoryID
		}
	}

	if vendorIDStr := c.Query("vendor_id"); vendorIDStr != "" {
		if vendorID, err := uuid.Parse(vendorIDStr); err == nil {
			filters.VendorID = &vendorID
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	if searchTerm := c.Query("search"); searchTerm != "" {
		filters.SearchTerm = &searchTerm
	}

	if inStockStr := c.Query("in_stock"); inStockStr != "" {
		if inStock, err := strconv.ParseBool(inStockStr); err == nil {
			filters.InStock = &inStock
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	result, err := ctl.svc.SearchProducts(c.Context(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to search products",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    result,
		"message": "Products retrieved successfully",
	})
}

// GetProduct godoc
// @Summary      Get product (public)
// @Tags         shop
// @Produce      json
// @Param        id   path      string true "Product ID (UUID)"
// @Success      200  {object}  products.GetProductDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /shop/products/{id} [get]
func (ctl *Controller) GetProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
		})
	}

	product, err := ctl.svc.GetProduct(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    product,
		"message": "Product retrieved successfully",
	})
}

// GetCategories godoc
// @Summary      List categories (public)
// @Tags         shop
// @Produce      json
// @Success      200 {array}  product_categories.GetProductCategoryDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /shop/categories [get]
func (ctl *Controller) GetCategories(c *fiber.Ctx) error {
	categories, err := ctl.svc.GetCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve categories",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    categories,
		"message": "Categories retrieved successfully",
	})
}

// GetCategory godoc
// @Summary      Get category (public)
// @Tags         shop
// @Produce      json
// @Param        id   path      string true "Category ID (UUID)"
// @Success      200  {object}  product_categories.GetProductCategoryDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /shop/categories/{id} [get]
func (ctl *Controller) GetCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid category ID",
		})
	}

	category, err := ctl.svc.GetCategory(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    category,
		"message": "Category retrieved successfully",
	})
}

// GetVendors godoc
// @Summary      List vendors (public)
// @Tags         shop
// @Produce      json
// @Success      200 {array}  vendors.GetVendorDTO
// @Failure      500 {object} map[string]interface{}
// @Router       /shop/vendors [get]
func (ctl *Controller) GetVendors(c *fiber.Ctx) error {
	vendors, err := ctl.svc.GetVendors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve vendors",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    vendors,
		"message": "Vendors retrieved successfully",
	})
}

// GetVendor godoc
// @Summary      Get vendor (public)
// @Tags         shop
// @Produce      json
// @Param        id   path      string true "Vendor ID (UUID)"
// @Success      200  {object}  vendors.GetVendorDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /shop/vendors/{id} [get]
func (ctl *Controller) GetVendor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid vendor ID",
		})
	}

	vendor, err := ctl.svc.GetVendor(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Vendor not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    vendor,
		"message": "Vendor retrieved successfully",
	})
}
