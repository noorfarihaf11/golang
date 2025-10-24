package routes

import (
 	"github.com/noorfarihaf11/clean-arc/app/service"
	"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(app *fiber.App, service service.FileService) {
api := app.Group("/api")
files := api.Group("/files")
files.Post("/upload", service.UploadFile)
files.Get("/", service.GetAllFiles)
files.Get("/:id", service.GetFileByID)
files.Delete("/:id", service.DeleteFile)
}
