package service

import (
	"fmt"
	"log"
	_ "os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
)

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


func GetTotalJobAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID := c.Params("id") // langsung string

	results, err := repository.GetTotalJobAlumni(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Alumni tidak lebih dari 1 pekerjaan " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mendapatkan data alumni yang memiliki lebih dari 1 pekerjaan",
		"success": true,
		"alumni":  results,
	})
}

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
