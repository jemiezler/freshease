package main

import (
	"context"
	"freshease/backend/internal/common/config"
	"freshease/backend/internal/common/http"
	"freshease/backend/internal/common/middleware"
	"os/signal"
	"syscall"
	"time"

	_ "freshease/backend/internal/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2/log"
)

// @title Freshease API
// @version 1.0
// @description API docs for Freshease backend.
// @BasePath /api
// @host localhost:8080
func main() {
	cfg := config.Load()
	app := http.New()
	app.Use(middleware.RequestLogger())
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	// api := app.Group("/api")

	go func() {
		log.Info("listening on", cfg.HTTPAddr)
		if err := app.Listen(cfg.HTTPAddr); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	stop()

	shCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(shCtx)
}
