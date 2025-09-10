package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/service"
	"github.com/noorfarihaf11/clean-arc/middleware"
)

func AlumniRoutes(api fiber.Router, db *sql.DB) {
	alumni := api.Group("/unair/alumni", middleware.AuthRequired())

	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})

	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})

	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})

	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})

	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})
}
