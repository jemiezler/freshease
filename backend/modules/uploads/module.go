package uploads

import (
	"freshease/backend/internal/common/config"

	"github.com/gofiber/fiber/v2"
)

// RegisterModule wires service -> controller and mounts routes.
func RegisterModule(api fiber.Router, cfg config.MinIOConfig) error {
	svc, err := NewService(cfg)
	if err != nil {
		return err
	}
	ctl := NewController(svc)
	Routes(api, ctl)
	return nil
}
