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
		return service.GetAllJobService(c, db)
	})

	job.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetJobByIDService(c, db)
	})

	job.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetJobsByAlumniIDService(c, db)
	})

	job.Post("/", func(c *fiber.Ctx) error {
		return service.CreateJobService(c, db)
	})

	job.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateJobService(c, db)
	})

	job.Delete("/:id", func(c *fiber.Ctx) error {
		return service.DeleteJobService(c, db)
	})

	job.Get("/filter/getbyrole", func(c *fiber.Ctx) error {
		return service.GetJobByRoleService(c, db)
	})

	job.Get("/filter/jobmoretwo/:id", func(c *fiber.Ctx) error {
		return service.GetTotalJobAlumniService(c, db)
	})

	job.Put("/update/:id", func(c *fiber.Ctx) error {
		return service.UpdateJobByRoleService(c, db)
	})

	job.Get("/filter/trash", func(c *fiber.Ctx) error {
		return service.GetTrashService(c, db)
	})

	job.Put("/filter/restore/:id", func(c *fiber.Ctx) error {
		return service.RestoreService(c, db)
	})

	job.Delete("/filter/delete/:id", func(c *fiber.Ctx) error {
		return service.HardDeleteService(c, db)
	})

}
