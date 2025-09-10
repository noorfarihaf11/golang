package service

import (
	"database/sql"
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
        "message": "Berhasil mendapatkan semua data pekerjaan alumni",
        "success": true,
        "pekerjaan alumni":  PekerjaanList,
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
        "message": "Berhasil mendapatkan data pekerjaan",
        "success": true,
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
        "message": "Pekerjaan berhasil ditambahkan",
        "success": true,
        "pekerjaan_alumni":     savedJob,
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
        "message": "pekerjaan berhasil diperbarui",
        "success": true,
        "pekerjaan_alumni": updatedJob,
    })
}

func DeleteJobService(c *fiber.Ctx, db *sql.DB) error {
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

    //  Ambil ID dari URL
    id := c.Params("id")

    //  Hapus alumni dari DB
    err = repository.DeleteJob(db, id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus pekerjaan alumni: " + err.Error(),
            "success": false,
        })
    }

    username := c.Locals("username").(string) 
    log.Printf("Admin %s menghapus pekerjaan ID %s", username, id)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Pekerjaan alumni berhasil dihapus",
        "success": true,
    })
}
