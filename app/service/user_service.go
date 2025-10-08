package service

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
)

func LoginService(db *sql.DB, req model.LoginRequest) (string, model.User, error) {
	var user model.User
	var passwordHash string

	err := db.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE username = $1 OR email = $1
	`, req.Username).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		return "", user, errors.New("username atau password salah")
	}

	// cek password
	if !utils.CheckPassword(req.Password, passwordHash) {
		return "", user, errors.New("password salah")
	}

	// generate JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		return "", user, errors.New("gagal generate token")
	}

	return token, user, nil
}

func RegisterService(c *fiber.Ctx, db *sql.DB) error {
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
		Role:         "alumni", // default alumni
	}

	// simpan ke DB via repo
	createdUser, err := repository.RegisterUser(db, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user: " + err.Error(),
			"success": false,
		})
	}

	// generate token
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
