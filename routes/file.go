package routes

import (
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gofiber/fiber/v2"

	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/app/service"
	"github.com/noorfarihaf11/clean-arc/middleware"
)

func FileRoutes(api fiber.Router, db *mongo.Database) {
	// Group utama dengan middleware login
	files := api.Group("api/files", middleware.AuthRequired())

	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, "./uploads")

	files.Post("/upload/photo/:user_id", middleware.UserAccessMiddleware(), func(c *fiber.Ctx) error {
		return fileService.UploadPhoto(c)
	})
	files.Post("/upload/certificate/:user_id", middleware.UserAccessMiddleware(), func(c *fiber.Ctx) error {
		return fileService.UploadCertificate(c)
	})

	// Endpoint lain
	files.Get("/", func(c *fiber.Ctx) error {
		return fileService.GetAllFiles(c)
	})
	files.Get("/:id", func(c *fiber.Ctx) error {
		return fileService.GetFileByID(c)
	})
	files.Delete("/:id", func(c *fiber.Ctx) error {
		return fileService.DeleteFile(c)
	})
}
