package service

import (
	"database/sql"
	"fmt"
	"log"
	_ "os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
)

func GetAllJobService(c *fiber.Ctx, db *sql.DB) error {
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

	PekerjaanList, err := repository.GetAllJobs(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Berhasil mendapatkan semua data pekerjaan alumni",
		"success":          true,
		"pekerjaan alumni": PekerjaanList,
	})
}

func GetJobByIDService(c *fiber.Ctx, db *sql.DB) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	// Format harus "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	// Validasi JWT pakai utils
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	p, err := repository.GetJobByID(db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error: " + err.Error(),
			"success": false,
		})
	}

	if p == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pekerjaan tidak ditemukan",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Berhasil mendapatkan data pekerjaan",
		"success":          true,
		"pekerjaan_alumni": p,
	})
}
func GetJobsByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	// Format harus "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	// Validasi JWT pakai utils
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "alumni_id tidak valid",
			"success": false,
		})
	}

	jobs, err := repository.GetJobsByAlumniID(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada pekerjaan untuk alumni ini",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mendapatkan data pekerjaan",
		"success": true,
		"jobs":    jobs,
	})
}
func CreateJobService(c *fiber.Ctx, db *sql.DB) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	// Format harus "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	// Validasi JWT pakai utils
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	var pekerjaan_alumni model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan_alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid",
			"success": false,
		})
	}

	username := c.Locals("username").(string)
	log.Printf("Admin %s menambah pekerjaan baru", username)

	savedJob, err := repository.CreateJob(db, &pekerjaan_alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":          "Pekerjaan berhasil ditambahkan",
		"success":          true,
		"pekerjaan_alumni": savedJob,
	})
}

func UpdateJobService(c *fiber.Ctx, db *sql.DB) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header tidak ada",
			"success": false,
		})
	}

	// Format harus "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Format Authorization salah, gunakan Bearer <token>",
			"success": false,
		})
	}

	// Validasi JWT pakai utils
	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Ambil ID dari URL
	id := c.Params("id")

	// Parse body ke struct Alumni
	var pekerjaan_alumni model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan_alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid",
			"success": false,
		})
	}

	username := c.Locals("username").(string)
	log.Printf("Admin %s mengubah data pekerjaan ID %s", username, id)

	// Update ke DB
	updatedJob, err := repository.UpdateJob(db, id, &pekerjaan_alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "pekerjaan berhasil diperbarui",
		"success":          true,
		"pekerjaan_alumni": updatedJob,
	})
}

func GetTotalJobAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}
	results, err := repository.GetTotalJobAlumni(db, id)
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

func DeleteJobService(c *fiber.Ctx, db *sql.DB) error {
	userIDAny := c.Locals("user_id")
	roleAny := c.Locals("role")

	if userIDAny == nil || roleAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: claims tidak ditemukan",
			"success": false,
		})
	}

	userID := userIDAny.(int)
	role := roleAny.(string)
	jobID := c.Params("id")

	rows, err := repository.SoftDeleteJob(db, jobID, userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus data: " + err.Error(),
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
		"message": "Pekerjaan berhasil dihapus (soft delete)",
		"success": true,
	})
}

func GetJobByRoleService(c *fiber.Ctx, db *sql.DB) error {
	userIDAny := c.Locals("user_id")
	roleAny := c.Locals("role")

	if userIDAny == nil || roleAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: claims tidak ditemukan",
			"success": false,
		})
	}

	userID := userIDAny.(int)
	role := roleAny.(string)

	jobs, err := repository.GetJobByRole(db, userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada pekerjaan ditemukan",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data pekerjaan berhasil diambil",
		"success": true,
		"data":    jobs,
	})
}

func UpdateJobByRoleService(c *fiber.Ctx, db *sql.DB) error {
	// Ambil data dari JWT claims
	userIDAny := c.Locals("user_id")
	roleAny := c.Locals("role")

	if userIDAny == nil || roleAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: claims tidak ditemukan",
			"success": false,
		})
	}

	userID := userIDAny.(int)
	role := roleAny.(string)

	// Ambil parameter id dari URL
	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID pekerjaan tidak diberikan",
			"success": false,
		})
	}

	// Parsing request body ke struct model
	var req model.PekerjaanAlumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal mem-parsing request body: " + err.Error(),
			"success": false,
		})
	}

	// Panggil repository
	updatedJob, err := repository.UpdateJobByRole(db, jobID, userID, role, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memperbarui pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data pekerjaan berhasil diperbarui",
		"success": true,
		"data":    updatedJob,
	})
}

func CreateJobByRoleService(c *fiber.Ctx, db *sql.DB) error {
	userIDAny := c.Locals("user_id")
	roleAny := c.Locals("role")

	if userIDAny == nil || roleAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: claims tidak ditemukan",
		})
	}

	userID := userIDAny.(int)
	role := roleAny.(string)

	var pekerjaan model.PekerjaanAlumni
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gagal parse body: " + err.Error(),
		})
	}

	created, err := repository.CreateJobByRole(db, userID, role, &pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat pekerjaan: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil dibuat",
		"data":    created,
	})
}

func GetTrashService(c *fiber.Ctx, db *sql.DB) error {

	userIDAny := c.Locals("user_id")
	roleAny := c.Locals("role")

	if userIDAny == nil || roleAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: claims tidak ditemukan",
			"success": false,
		})
	}

	userID := userIDAny.(int)
	role := roleAny.(string)

	jobs, err := repository.GetTrash(db, userID, role)
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

func RestoreService(c *fiber.Ctx, db *sql.DB) error {

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
func HardDeleteService(c *fiber.Ctx, db *sql.DB) error {

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
