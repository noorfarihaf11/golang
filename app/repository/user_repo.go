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

	// ✅ Simpan user baru ke collection users
	_, err := db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("gagal menambahkan user: %v", err)
	}

	// 🧠 Debug untuk memastikan role terbaca
	fmt.Println("DEBUG ROLE:", user.Role)

	role := strings.TrimSpace(strings.ToLower(user.Role))
	switch role {
	case "alumni":
		// ✅ Jika user role = alumni, otomatis buat data di koleksi alumni
		alumni := model.Alumni{
			ID:        primitive.NewObjectID(),
			UserID:    &user.ID, // pakai pointer, karena field opsional
			Nama:      user.Username,
			Email:     user.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = db.Collection("alumni").InsertOne(ctx, alumni)
		if err != nil {
			return nil, fmt.Errorf("gagal menambahkan data alumni: %v", err)
		}

		fmt.Println("✅ Data alumni berhasil dibuat otomatis untuk user:", user.Username)

	default:
		fmt.Println("➡️ Role bukan alumni, tidak dimasukkan ke koleksi alumni.")
	}

	return user, nil
}
