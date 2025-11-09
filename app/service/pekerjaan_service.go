package service

import (
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
)

// GetAllJobService godoc
// @Summary Mendapatkan semua data pekerjaan
// @Description Mengembalikan daftar semua pekerjaan alumni yang tersedia
// @Tags PekerjaanAlumni
// @Accept json
// @Produce json
// @Success 200 {object} model.PekerjaanAlumniResponse "Berhasil mengambil data pekerjaan alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/pekerjaan [get]
func GetAllJobService(c *fiber.Ctx, db *mongo.Database) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}
	if _, err := utils.ValidateToken(token); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
	}

	jobs, err := repository.GetAllJobs(db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if len(jobs) == 0 {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data kosong"})
	}

	return c.JSON(fiber.Map{"success": true, "data": jobs})
}

// GetJobByIDService godoc
// @Summary Mendapatkan pekerjaan berdasarkan ID
// @Description Mengambil data pekerjaan sesuai ID pekerjaan
// @Tags PekerjaanAlumni
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} model.SinglePekerjaanResponse "Berhasil mengambil data pekerjaan alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/pekerjaan/{id} [get]
func GetJobByIDService(c *fiber.Ctx, db *mongo.Database) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}
	if _, err := utils.ValidateToken(token); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
	}

	job, err := repository.GetJobByID(db, c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if job == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Pekerjaan tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": job})
}

// GetJobsByAlumniIDService godoc
// @Summary Mendapatkan pekerjaan berdasarkan ID alumni
// @Description Mengembalikan semua pekerjaan yang dimiliki oleh seorang alumni
// @Tags PekerjaanAlumni
// @Accept json
// @Produce json
// @Param alumni_id path string true "ID Alumni"
// @Success 200 {object} model.SinglePekerjaanResponse "Berhasil mengambil data pekerjaan alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/pekerjaan/alumni/{alumni_id} [get]
func GetJobsByAlumniIDService(c *fiber.Ctx, db *mongo.Database) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}
	if _, err := utils.ValidateToken(token); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
	}

	id := c.Params("alumni_id")
	jobs, err := repository.GetJobsByAlumniID(db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if len(jobs) == 0 {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Tidak ada pekerjaan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": jobs})
}

// CreateJobService godoc
// @Summary Menambahkan pekerjaan baru
// @Description Membuat entri pekerjaan baru untuk alumni
// @Tags PekerjaanAlumni
// @Accept json
// @Produce json
// @Param data body model.PekerjaanAlumni true "Data Pekerjaan"
// @Success 200 {object} model.SinglePekerjaanResponse "Berhasil menambahkan data pekerjaan alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/pekerjaan [post]
func CreateJobService(c *fiber.Ctx, db *mongo.Database) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}
	if _, err := utils.ValidateToken(token); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
	}

	var job model.PekerjaanAlumni
	if err := c.BodyParser(&job); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Body tidak valid"})
	}

	res, err := repository.CreateJob(db, &job)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Berhasil tambah pekerjaan", "data": res})
}

// UpdateJobService godoc
// @Summary Memperbarui data pekerjaan
// @Description Mengubah data pekerjaan berdasarkan ID
// @Tags PekerjaanAlumni
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Param data body model.PekerjaanAlumni true "Data Pekerjaan Baru"
// @Success 200 {object} model.SinglePekerjaanResponse "Berhasil mengupdate data pekerjaan alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/pekerjaan/{id} [put]
func UpdateJobService(c *fiber.Ctx, db *mongo.Database) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}
	if _, err := utils.ValidateToken(token); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
	}

	var job model.PekerjaanAlumni
	if err := c.BodyParser(&job); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Body tidak valid"})
	}

	res, err := repository.UpdateJob(db, c.Params("id"), job)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Berhasil update pekerjaan", "data": res})
}

// DeleteJobService godoc
// @Summary Menghapus pekerjaan (soft delete)
// @Description Menghapus data pekerjaan secara soft delete
// @Tags PekerjaanAlumni
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} model.SinglePekerjaanResponse "Berhasil menghapus data pekerjaan alumni"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /unair/pekerjaan/filter/trash/{id} [put]
func DeleteJobService(c *fiber.Ctx, db *mongo.Database) error {
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

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID pekerjaan wajib diisi",
			"success": false,
		})
	}

	userID := claims.UserID.Hex()
	role := claims.Role
	if role == "" {
		role = "alumni"
	}

	log.Printf("User %s (role: %s) menghapus pekerjaan ID %s", userID, role, id)

	err = repository.SoftDeleteJob(db, id, userID, role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus",
		"success": true,
	})
}

// GetTrashService godoc
// @Summary Mendapatkan daftar data yang ada di trash
// @Description Mengambil semua data pekerjaan yang sudah dihapus (soft delete) berdasarkan role pengguna
// @Tags Trash
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.TrashResponse "Berhasil mengambil data trash"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Tidak ada trash"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil trash"
// @Router /unair/pekerjaan/filter/trash [get]
func GetTrashService(c *fiber.Ctx, db *mongo.Database) error {

	// Ambil user_id dan role dari JWT claims
	userIDAny := c.Locals("user_id")
	roleAny := c.Locals("role")

	if userIDAny == nil || roleAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: claims tidak ditemukan",
			"success": false,
		})
	}

	userIDHex, ok := userIDAny.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "user_id tidak valid",
			"success": false,
		})
	}

	role, ok := roleAny.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "role tidak valid",
			"success": false,
		})
	}

	// Panggil repository
	jobs, err := repository.GetTrash(db, userIDHex, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil trash: " + err.Error(),
			"success": false,
		})
	}

	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada trash",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data trash berhasil diambil",
		"success": true,
		"data":    jobs,
	})
}

// RestoreService godoc
// @Summary Mengembalikan data dari trash
// @Description Melakukan restore terhadap pekerjaan berdasarkan ID
// @Tags Trash
// @Param id path string true "ID pekerjaan"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.TrashResponse "Berhasil restore pekerjaan"
// @Failure 403 {object} map[string]interface{} "Tidak diizinkan menghapus pekerjaan"
// @Failure 500 {object} map[string]interface{} "Gagal restore data"
// @Router /unair/pekerjaan/filter/restore/{id} [get]
func RestoreService(c *fiber.Ctx, db *mongo.Database) error {

	jobID := c.Params("id")

	rows, err := repository.Restore(db, jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal restore data: " + err.Error(),
			"success": false,
		})
	}

	if rows == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak diizinkan menghapus pekerjaan dengan ID %s", jobID),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil direstore",
		"success": true,
	})
}

// HardDeleteService godoc
// @Summary Menghapus data secara permanen
// @Description Menghapus pekerjaan dari trash berdasarkan ID secara permanen
// @Tags Trash
// @Param id path string true "ID pekerjaan"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Pekerjaan berhasil dihapus permanen"
// @Failure 403 {object} map[string]interface{} "Tidak diizinkan menghapus pekerjaan"
// @Failure 500 {object} map[string]interface{} "Gagal delete data"
// @Router /unair/pekerjaan/filter/delete/{id} [delete]
func HardDeleteService(c *fiber.Ctx, db *mongo.Database) error {

	jobID := c.Params("id")

	rows, err := repository.HardDelete(db, jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal delete data: " + err.Error(),
			"success": false,
		})
	}

	if rows == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": fmt.Sprintf("Tidak diizinkan menghapus pekerjaan dengan ID %s", jobID),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus",
		"success": true,
	})
}