package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/noorfarihaf11/clean-arc/app/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(db *mongo.Database, user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	// ✅ Masukkan user ke collection users
	_, err := db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("gagal menambahkan user: %v", err)
	}

	// 🧠 Debug untuk memastikan role terbaca
	fmt.Println("DEBUG ROLE:", user.Role)

	// ✅ Hanya jika role == "alumni", insert juga ke koleksi alumni
	if strings.TrimSpace(strings.ToLower(user.Role)) == "alumni" {
		alumni := model.Alumni{
			ID:        primitive.NewObjectID(),
			UserID:    user.ID,
			Nama:      user.Username,
			Email:     user.Email,
			CreatedAt: time.Now(),
		}

		_, err = db.Collection("alumni").InsertOne(ctx, alumni)
		if err != nil {
			return nil, fmt.Errorf("gagal menambahkan data alumni: %v", err)
		}
	} else {
		fmt.Println("➡️ Role bukan alumni, tidak dimasukkan ke tabel alumni.")
	}

	return user, nil
}
