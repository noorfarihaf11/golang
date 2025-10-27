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
)

func main() {
	config.LoadEnv()

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

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Semua route terpusat di sini
	routes.Routes(app, db)

	log.Printf("Server running on port %s 🚀", port)
	log.Fatal(app.Listen(":" + port))
}
