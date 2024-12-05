package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/naufalatha/go-boilerplate/config"
	"github.com/naufalatha/go-boilerplate/handlers"
)

type router struct {
	handler *handlers.Handler
	config  *config.Configuration
}

func InitRouter(handler *handlers.Handler, config *config.Configuration) router {
	return router{
		handler: handler,
		config:  config,
	}
}

func (r router) Route(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Welcome to golang boilerplate service %s", r.config.AppEnv))
	})

	router.Get("/metrics", monitor.New(monitor.Config{Title: "Pippin Metrics Monitoring Page"}))
}

func (r router) JWTRoute(router fiber.Router) {
	router.Get("/check-health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
