package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
)

func LoginService(db *mongo.Database, req model.LoginRequest) (string, *model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User

	// cari user berdasarkan username atau email
	filter := bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Username},
		},
	}

	err := db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil, errors.New("username atau password salah")
		}
		return "", nil, err
	}

	// cek password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return "", nil, errors.New("password salah")
	}

	// generate JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", nil, errors.New("gagal generate token")
	}

	return token, &user, nil
}

func RegisterService(c *fiber.Ctx, db *mongo.Database) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal parse request",
			"success": false,
		})
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal enkripsi password",
			"success": false,
		})
	}

	// buat user struct
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role, // default
		CreatedAt:    time.Now(),
	}

	// simpan ke DB via repository
	createdUser, err := repository.RegisterUser(db, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user: " + err.Error(),
			"success": false,
		})
	}

	// generate token JWT
	token, err := utils.GenerateToken(*createdUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token JWT",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"success": true,
		"token":   token,
		"user":    createdUser,
	})
}
