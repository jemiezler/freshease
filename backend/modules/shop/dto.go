package shop

import (
	"time"

	"github.com/google/uuid"
)

// ShopProductDTO represents a product in the shop view
type ShopProductDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	UnitLabel   string    `json:"unit_label"`
	IsActive    string    `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Vendor information
	VendorID   uuid.UUID `json:"vendor_id"`
	VendorName string    `json:"vendor_name"`

	// Category information
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`

	// Inventory information
	StockQuantity int  `json:"stock_quantity"`
	IsInStock     bool `json:"is_in_stock"`
}

// ShopCategoryDTO represents a product category in the shop view
type ShopCategoryDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ShopVendorDTO represents a vendor in the shop view
type ShopVendorDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`
	Website     string    `json:"website"`
	LogoURL     string    `json:"logo_url"`
	Description string    `json:"description"`
	IsActive    string    `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShopSearchFilters represents search and filter parameters
type ShopSearchFilters struct {
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	VendorID   *uuid.UUID `json:"vendor_id,omitempty"`
	MinPrice   *float64   `json:"min_price,omitempty"`
	MaxPrice   *float64   `json:"max_price,omitempty"`
	SearchTerm *string    `json:"search_term,omitempty"`
	InStock    *bool      `json:"in_stock,omitempty"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// ShopSearchResponse represents the response for shop search
type ShopSearchResponse struct {
	Products []*ShopProductDTO `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	HasMore  bool              `json:"has_more"`
}
