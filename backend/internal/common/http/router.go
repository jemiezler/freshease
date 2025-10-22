package http

import (
	"freshease/backend/ent"
	"freshease/backend/modules/permissions"
	"freshease/backend/modules/roles"
	"freshease/backend/modules/users"
	"slices"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RegisterRoutes(app *fiber.App, client *ent.Client) {
	api := app.Group("/api")

	log.Debug("[router] registering modules...")

	users.RegisterModuleWithEnt(api, client)
	roles.RegisterModuleWithEnt(api, client)
	permissions.RegisterModuleWithEnt(api, client)

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
