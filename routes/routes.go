package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, db *sql.DB) {
	api := app.Group("/")

	AuthRoutes(api, db)
	AlumniRoutes(api, db)
	JobRoutes(api, db)
}
