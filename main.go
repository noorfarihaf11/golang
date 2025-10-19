package main

import (
    "log"
    "os"
   	"github.com/noorfarihaf11/clean-arc/database"
	"github.com/noorfarihaf11/clean-arc/config"
	"github.com/noorfarihaf11/clean-arc/routes"
)

func main() {
	config.LoadEnv()

	// Koneksi MongoDB
	db, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatalf("Gagal konek ke MongoDB: %v", err)
	}

	app := config.NewApp(db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	routes.Routes(app, db)

	log.Printf("Server running on port %s 🚀", port)
	log.Fatal(app.Listen(":" + port))
}
