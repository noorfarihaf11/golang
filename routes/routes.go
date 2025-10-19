package routes

import (
    "go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, db *mongo.Database) {
	api := app.Group("/")

	AuthRoutes(api, db)
	AlumniRoutes(api, db)
	JobRoutes(api, db)
}
