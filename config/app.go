package config

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/middleware"
	"github.com/noorfarihaf11/clean-arc/routes"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()

	// Middleware global
	app.Use(middleware.LoggerMiddleware)

	// Register semua routes
	routes.AlumniRoutes(app, db)
	routes.JobRoutes(app, db)

	return app
}
