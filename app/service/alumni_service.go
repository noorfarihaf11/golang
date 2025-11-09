package service

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllAlumniService godoc
// @Summary Mengambil semua data alumni
// @Description Mengembalikan daftar semua alumni yang terdaftar dalam sistem
// @Tags Alumni
// @Accept json
// @Produce json
// @Success 200 {object} model.AlumniResponse "Berhasil mengambil data alumni"
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/alumni [get]
func GetAllAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	alumniList, err := repository.GetAllAlumni(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mendapatkan semua data alumni",
		"success": true,
		"alumni":  alumniList,
	})
}

// GetAlumniByIDService godoc
// @Summary Mengambil data alumni berdasarkan ID
// @Description Mengembalikan detail satu alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 201 {object} model.SingleAlumniResponse "Berhasil mengambil data alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/alumni/{id} [get]
func GetAlumniByIDService(c *fiber.Ctx, db *mongo.Database) error {
	id := c.Params("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Message: "ID tidak valid",
			Code:    fiber.StatusBadRequest,
		})
	}

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Success: false,
				Message: "Alumni tidak ditemukan",
				Code:    fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Message: "Gagal mengambil data alumni",
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mendapatkan data alumni",
		"alumni":  alumni,
	})
}


// CreateAlumniService godoc
// @Summary Menambahkan data alumni baru
// @Description Menambahkan data alumni baru ke dalam sistem
// @Tags Alumni
// @Accept json
// @Produce json
// @Param alumni body model.Alumni true "Data Alumni"
// @Success 201 {object} model.SingleAlumniResponse "Berhasil menambahkan data alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/alumni [post]
func CreateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid",
			"success": false,
		})
	}

	var userID *primitive.ObjectID
	if claims.Role == "admin" {
		log.Printf("Admin %s menambah alumni baru", claims.Username)
		userID = nil
	} else {
		log.Printf("User %s (alumni) menambah data dirinya sendiri", claims.Username)
		userID = &claims.UserID
	}

	savedAlumni, err := repository.CreateAlumni(db, &alumni, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alumni berhasil ditambahkan",
		"success": true,
		"alumni":  savedAlumni,
	})
}

// UpdateAlumniService godoc
// @Summary Mengubah data alumni
// @Description Mengubah data alumni yang sudah ada dalam sistem
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 201 {object} model.SingleAlumniResponse "Berhasil memperbarui data alumni"
// @Success 200 {object} model.Alumni
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/alumni/{id} [put]
func UpdateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	id := c.Params("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid",
			"success": false,
		})
	}

	updatedAlumni, err := repository.UpdateAlumni(db, id, &alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate alumni: " + err.Error(),
			"success": false,
		})
	}

	username, _ := c.Locals("username").(string)
	log.Printf("User %s mengubah data alumni ID %s", username, id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alumni berhasil diperbarui",
		"success": true,
		"alumni":  updatedAlumni,
	})
}

// DeleteAlumniService godoc
// @Summary Menghapus data alumni
// @Description Menghapus data alumni dari sistem
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 201 {object} model.SingleAlumniResponse "Berhasil menghapus data alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/alumni/{id} [delete]
func DeleteAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	id := c.Params("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	if err := repository.DeleteAlumni(db, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus alumni: " + err.Error(),
			"success": false,
		})
	}

	username, _ := c.Locals("username").(string)
	log.Printf("User %s menghapus data alumni ID %s", username, id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alumni berhasil dihapus",
		"success": true,
	})
}
