package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"freshease/backend/internal/common/config"
	"freshease/backend/internal/common/db"
	httpserver "freshease/backend/internal/common/http"

	_ "freshease/backend/internal/docs"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title Freshease API
// @version 1.0
// @description API docs for Freshease backend.
// @BasePath /api
func main() {
	cfg := config.Load()

	// Build HTTP app (no routes yet)
	app := httpserver.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}))
	// DB connect + ping with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, closeDB, err := db.NewEntClientPGX(ctx, cfg.DatabaseURL, cfg.Ent.Debug)
	if err != nil {
		log.Fatal("[Fatal] ent client: ", err)
	}
	defer func() { _ = closeDB(context.Background()) }()

	// Migrate
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal("[Fatal] ent schema: ", err)
	}

	// Register routes, grouped under /api
	httpserver.RegisterRoutes(app, client, cfg)

	// Start server
	go func() {
		log.Infof("[HTTP] listening on %s", cfg.HTTPPort)
		if err := app.Listen(cfg.HTTPPort); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for signal
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-sigCtx.Done()
	stop()

	// Graceful shutdown
	shCtx, shCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shCancel()
	_ = app.ShutdownWithContext(shCtx)
}
