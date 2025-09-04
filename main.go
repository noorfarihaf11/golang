package main

import (
    "log"
    "os"
   	"github.com/noorfarihaf11/clean-arc/database"
	"github.com/noorfarihaf11/clean-arc/config"
)

func main() {
	config.LoadEnv()
	db := database.ConnectDB()
	app := config.NewApp(db)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}