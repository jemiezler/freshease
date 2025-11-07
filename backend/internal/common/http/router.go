package http

import (
	"freshease/backend/ent"
	"freshease/backend/ent/user"
	"freshease/backend/internal/common/config"
	"freshease/backend/internal/common/middleware"
	"freshease/backend/modules/addresses"
	"freshease/backend/modules/auth/authoidc"
	"freshease/backend/modules/cart_items"
	"freshease/backend/modules/carts"
	"freshease/backend/modules/genai"
	"freshease/backend/modules/inventories"
	"freshease/backend/modules/permissions"
	"freshease/backend/modules/product_categories"
	"freshease/backend/modules/products"
	"freshease/backend/modules/roles"
	"freshease/backend/modules/shop"
	"freshease/backend/modules/uploads"
	"freshease/backend/modules/users"
	"freshease/backend/modules/vendors"
	"slices"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func RegisterRoutes(app *fiber.App, client *ent.Client, cfg config.Config) {
	api := app.Group("/api")

	log.Debug("[router] registering modules...")

	// 1) Public: OIDC auth (Google/LINE callbacks)
	if err := authoidc.RegisterModule(api, client); err != nil {
		panic(err)
	}
	genai.RegisterModuleWithEnt(api, client)

	// 2) Public: Shop API (no authentication required)
	shop.RegisterModuleWithEnt(api, client)

	// 3) File uploads (public for now, can be secured later)
	if err := uploads.RegisterModule(api, cfg.MinIO); err != nil {
		log.Fatalf("[router] failed to register uploads module: %v", err)
	}

	// mount protected modules on the secured router
	addresses.RegisterModuleWithEnt(api, client)
	cart_items.RegisterModuleWithEnt(api, client)
	carts.RegisterModuleWithEnt(api, client)
	inventories.RegisterModuleWithEnt(api, client)
	permissions.RegisterModuleWithEnt(api, client)
	products.RegisterModuleWithEnt(api, client)
	product_categories.RegisterModuleWithEnt(api, client)
	users.RegisterModuleWithEnt(api, client)
	roles.RegisterModuleWithEnt(api, client)
	vendors.RegisterModuleWithEnt(api, client)

	// 2) Secured area (everything below requires Authorization: Bearer <JWT>)
	secured := api.Group("", middleware.RequireAuth())

	secured.Get("/whoami", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		userEmail := c.Locals("user_email")

		if userID == nil || userEmail == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user not found in token"})
		}

		// Parse user ID from string to UUID
		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user ID"})
		}

		// Get user details from database
		user, err := client.User.Query().Where(user.ID(userUUID)).First(c.Context())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}

		return c.JSON(fiber.Map{
			"id":            user.ID.String(),
			"email":         user.Email,
			"name":          user.Name,
			"phone":         user.Phone,
			"bio":           user.Bio,
			"avatar":        user.Avatar,
			"cover":         user.Cover,
			"date_of_birth": user.DateOfBirth,
			"sex":           user.Sex,
			"goal":          user.Goal,
			"height_cm":     user.HeightCm,
			"weight_kg":     user.WeightKg,
			"status":        user.Status,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		})
	})

	logRegisteredModules(app, "/api")
}

func logRegisteredModules(app *fiber.App, apiPrefix string) {
	type stats struct {
		Routes  int
		Methods map[string]struct{}
	}
	modStats := map[string]*stats{}

	for _, r := range app.GetRoutes() {
		p := r.Path
		if !strings.HasPrefix(p, apiPrefix+"/") {
			continue
		}
		rest := strings.TrimPrefix(p, apiPrefix+"/")
		first := rest
		if i := strings.IndexByte(rest, '/'); i >= 0 {
			first = rest[:i]
		}
		if first == "" {
			continue
		}
		if _, ok := modStats[first]; !ok {
			modStats[first] = &stats{Methods: make(map[string]struct{})}
		}
		modStats[first].Routes++
		modStats[first].Methods[r.Method] = struct{}{}
	}

	if len(modStats) == 0 {
		log.Warn("[router] no modules discovered under ", apiPrefix)
		return
	}

	// Pretty log per module
	for name, s := range modStats {
		// collect method keys
		ms := make([]string, 0, len(s.Methods))
		for m := range s.Methods {
			ms = append(ms, m)
		}
		slices.Sort(ms)
		log.Infof("[router] module %-16s routes=%d methods=%v", name, s.Routes, ms)
	}

	// Also log a compact summary (sorted)
	keys := make([]string, 0, len(modStats))
	for k := range modStats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// var parts []string
	// for _, k := range keys {
	// 	parts = append(parts, k)
	// }
	// log.Infof("[router] modules: %s", strings.Join(parts, ", "))
}
