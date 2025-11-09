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

// LoginService godoc
// @Summary Masuk ke dalam sistem
// @Description Mengautentikasi user dan mengembalikan token JWT
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Data login user"
// @Success 200 {object} map[string]interface{} "token dan data user"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/login [post]
func LoginService(db *mongo.Database, req model.LoginRequest) (string, *model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User

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

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return "", nil, errors.New("password salah")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", nil, errors.New("gagal generate token")
	}

	return token, &user, nil
}

// RegisterService godoc
// @Summary Mendaftar ke dalam sistem
// @Description Membuat user baru dan mengembalikan token JWT
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Data registrasi user"
// @Success 200 {object} map[string]interface{} "Token dan data user yang terdaftar"
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/register [post]
func RegisterService(c *fiber.Ctx, db *mongo.Database) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Message: "Gagal parse request",
			Code:    fiber.StatusBadRequest,
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Message: "Gagal enkripsi password",
			Code:    fiber.StatusInternalServerError,
		})
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		CreatedAt:    time.Now(),
	}

	createdUser, err := repository.RegisterUser(db, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Message: "Gagal membuat user: " + err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	token, err := utils.GenerateToken(*createdUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Message: "Gagal membuat token JWT",
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Registrasi berhasil",
		"token":   token,
		"user":    createdUser,
	})
}
