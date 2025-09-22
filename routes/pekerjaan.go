package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/service"
	"github.com/noorfarihaf11/clean-arc/middleware"
)

func JobRoutes(api fiber.Router, db *sql.DB) {
	job := api.Group("/unair/pekerjaan", middleware.AuthRequired())

	job.Get("/", func(c *fiber.Ctx) error {
		return service.GetJobService(c, db)
	})

	job.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetJobByIDService(c, db)
	})

	job.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetJobsByAlumniIDService(c, db)
	})

	job.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateJobService(c, db)
	})

	job.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateJobService(c, db)
	})

	job.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteJobService(c, db)
	})

	job.Get("/filter/jobmoretwo/:id", func(c *fiber.Ctx) error {
		return service.GetTotalJobAlumniService(c, db)
	})
}
