package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/noorfarihaf11/clean-arc/config"
	"github.com/noorfarihaf11/clean-arc/database"
	"github.com/noorfarihaf11/clean-arc/routes"

	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/app/service"
)

func main() {
	config.LoadEnv()

	// Koneksi MongoDB
	db, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Gagal konek ke MongoDB: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Static folder untuk file upload
	app.Static("/uploads", "./uploads")

	// Port default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Routes project utama
	routes.Routes(app, db)

	// === Tambahan: fitur upload file ===
	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, "./uploads")

	// Gunakan fileService untuk daftarkan endpoint upload
	routes.SetupFileRoutes(app, fileService)

	log.Printf("Server running on port %s 🚀", port)
	log.Fatal(app.Listen(":" + port))
}
