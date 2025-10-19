package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() (*mongo.Database, error) {
	// Muat file .env
	_ = godotenv.Load()

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	if uri == "" || dbName == "" {
		log.Fatal("❌ MONGO_URI atau MONGO_DB_NAME belum diset di .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gunakan URI dari .env
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Tes koneksi
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("✅ Berhasil terhubung ke MongoDB!")

	db := client.Database(dbName)
	return db, nil
}
