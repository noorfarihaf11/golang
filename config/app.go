package config

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/service"
	"github.com/noorfarihaf11/clean-arc/middleware"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()
	app.Use(middleware.LoggerMiddleware)

	app.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})

	app.Get("/alumni/all/:key", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})

	return app
}
