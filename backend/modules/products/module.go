package products

import (
	"github.com/gofiber/fiber/v2"
	"freshease/backend/ent"
	"freshease/backend/modules/product_categories"
	"freshease/backend/modules/uploads"
)

// RegisterModuleWithEnt wires Ent repo -> service -> controller and mounts routes.
func RegisterModuleWithEnt(api fiber.Router, client *ent.Client, uploadsSvc uploads.Service) {
	repo := NewEntRepo(client)
	// Create product categories service for category support
	productCategoryRepo := product_categories.NewEntRepo(client)
	productCategorySvc := product_categories.NewService(productCategoryRepo)
	svc := NewServiceWithProductCategories(repo, uploadsSvc, productCategorySvc)
	ctl := NewController(svc)
	Routes(api, ctl)
}
