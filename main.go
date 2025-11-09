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
	"github.com/noorfarihaf11/clean-arc/docs" 
	fiberSwagger  "github.com/swaggo/fiber-swagger" 
)

func main() {
	config.LoadEnv()

	docs.SwaggerInfo.Title = "Clean Architecture API"
	docs.SwaggerInfo.Description = "Dokumentasi API untuk proyek Clean Architecture (Fiber + MongoDB)"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	db, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Gagal konek ke MongoDB: %v", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/uploads", "./uploads")
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Semua route terpusat di sini
	routes.Routes(app, db)

	log.Printf("Server running on port %s 🚀", port)
	log.Fatal(app.Listen(":" + port))
}
