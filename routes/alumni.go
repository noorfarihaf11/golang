package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/handler"
)

func AlumniRoutes(app *fiber.App, db *sql.DB) {
	h := handler.NewAlumniHandler(db)

	alumni := app.Group("/unair/alumni")
	alumni.Get("/", h.GetAll)
	alumni.Get("/:id", h.GetByID)
	alumni.Post("/", h.Create)
	alumni.Put("/:id", h.Update)
	alumni.Delete("/:id", h.Delete)
}

func JobRoutes(app *fiber.App, db *sql.DB) {
	h := handler.NewJobHandler(db)

	pekerjaan := app.Group("/unair/pekerjaan")
	pekerjaan.Get("/", h.GetAll)
	pekerjaan.Get("/:id", h.GetByID)
	pekerjaan.Post("/", h.Create)
	pekerjaan.Put("/:id", h.Update)
	pekerjaan.Delete("/:id", h.Delete)
}
