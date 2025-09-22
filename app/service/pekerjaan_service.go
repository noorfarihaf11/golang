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
            "alumni": results,
        })
    }

    
func GetJobService(c *fiber.Ctx, db *sql.DB) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")

    offset := (page - 1) * limit

    sortByWhitelist := map[string]bool{"id": true, 
    "alumni_id": true, "nama_perusahaan": true, "posisi_jabatan": true, 
    "bidang_industri": true, "lokasi_kerja": true, 
    "gaji_range": true, "tanggal_mulai_kerja": true, 
    "tanggal_selesai_kerja": true, "status_pekerjaan": true, 
    "deskripsi_pekerjaan": true,     
    }
    if !sortByWhitelist[sortBy] {
        sortBy = "id"
    }
    if strings.ToLower(order) != "desc" {
        order = "asc"
    }

    pekerjaanList, err := repository.GetJobRepo(db, search, sortBy, order, limit, offset)
    if err != nil {
        fmt.Println("Query error:", err) // untuk cek di console
         return c.Status(500).JSON(fiber.Map{
        "error": err.Error(), // kirim error asli ke client
    })
    }

    total, err := repository.CountJobRepo(db, search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to count pekerjaan alumni"})
    }

    response := model.PekerjaanAlumniResponse{
        Data: pekerjaanList,
        Meta: model.MetaInfo{
            Page:   page,
            Limit:  limit,
            Total:  total,
            Pages:  (total + limit - 1) / limit,
            SortBy: sortBy,
            Order:  order,
            Search: search,
        },
    }

    return c.JSON(response)
}


