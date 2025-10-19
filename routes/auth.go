package routes

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/service"
)

func AuthRoutes(api fiber.Router, db *mongo.Database) {
	
	api.Post("/api/login", func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Request body tidak valid",
			})
		}

		token, user, err := service.LoginService(db, req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Login berhasil",
			"data": fiber.Map{
				"user":  user,
				"token": token,
			},
		})
	})

	api.Post("/api/register", func(c *fiber.Ctx) error {
		return service.RegisterService(c, db)
	})
}
