package http

import (
	"freshease/backend/ent"
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
	"freshease/backend/modules/users"
	"freshease/backend/modules/vendors"
	"slices"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RegisterRoutes(app *fiber.App, client *ent.Client) {
	api := app.Group("/api")

	log.Debug("[router] registering modules...")

	// 1) Public: OIDC auth (Google/LINE callbacks)
	if err := authoidc.RegisterModule(api, client); err != nil {
		panic(err)
	}
	genai.RegisterModuleWithEnt(api, client)
	// 2) Secured area (everything below requires Authorization: Bearer <JWT>)
	secured := api.Group("", middleware.RequireAuth())

	// mount protected modules on the secured router
	addresses.RegisterModuleWithEnt(secured, client)
	cart_items.RegisterModuleWithEnt(secured, client)
	carts.RegisterModuleWithEnt(secured, client)
	inventories.RegisterModuleWithEnt(secured, client)
	permissions.RegisterModuleWithEnt(secured, client)
	products.RegisterModuleWithEnt(secured, client)
	product_categories.RegisterModuleWithEnt(secured, client)
	users.RegisterModuleWithEnt(secured, client)
	roles.RegisterModuleWithEnt(secured, client)
	vendors.RegisterModuleWithEnt(secured, client)

	secured.Get("/whoami", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
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
