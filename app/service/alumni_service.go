package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
)

func CheckAlumniService(c *fiber.Ctx, db *sql.DB) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Key	tidak	valid",
			"success": false,
		})
	}
	nim := c.FormValue("nim")
	if nim == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NIM	wajib	diisi",
			"success": false,
		})
	}
	alumni, err := repository.CheckAlumniByNim(db, nim)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":  "Mahasiswa	bukan	alumni",
				"success":  true,
				"isAlumni": false,
			})
		}
		 return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "message": "Gagal cek alumni karena " + err.Error(),
        "success": false,
    })
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Berhasil	mendapatkan	data	alumni",
		"success":  true,
		"isAlumni": true,
		"alumni":   alumni,
	})
}

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
    // Ambil Authorization header
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

    // Ambil data alumni
    alumniList, err := repository.GetAllAlumni(db)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan semua data alumni",
        "success": true,
        "alumni":  alumniList,
    })
}

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
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

    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "ID tidak valid",
            "success": false,
        })
    }

    alumni, err := repository.GetAlumniByID(db, id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error: " + err.Error(),
            "success": false,
        })
    }

    if alumni == nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Alumni tidak ditemukan",
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan data alumni",
        "success": true,
        "alumni": alumni,
    })
}

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

    var alumni model.Alumni
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Request body tidak valid",
            "success": false,
        })
    }


    userID := c.Locals("user_id").(int)

    username := c.Locals("username").(string) 
    log.Printf("Admin %s menambah alumni baru", username)

    // insert ke DB
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
        "alumni":  savedAlumni, // sudah ada ID dari RETURNING
    })
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
    // Validasi API Key
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

    // Ambil ID dari URL
    id := c.Params("id")

    // Parse body ke struct Alumni
    var alumni model.Alumni
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Request body tidak valid",
            "success": false,
        })
    }

    // Update ke DB
    updatedAlumni, err := repository.UpdateAlumni(db, id, &alumni)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengupdate alumni: " + err.Error(),
            "success": false,
        })
    }

    username := c.Locals("username").(string) 
    log.Printf("Admin %s mengubah data alumni ID %s", username, id)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Alumni berhasil diperbarui",
        "success": true,
        "alumni": updatedAlumni,
    })
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
    // Validasi API Key
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
    // Ambil ID dari URL
    id := c.Params("id")

    // Hapus alumni dari DB
    err = repository.DeleteAlumni(db, id) 
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus alumni: " + err.Error(),
            "success": false,
        })
    }

    username := c.Locals("username").(string) 
    log.Printf("Admin %s menghapus mahasiswa ID %s", username, id)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Alumni berhasil dihapus",
        "success": true,
    })
}

func GetAlumniBySalaryService(c *fiber.Ctx, db *sql.DB) error {
    alumnis, err := repository.GetAlumniWithHighSalary(db)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengambil data alumni: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan data alumni dengan gaji > 20 juta",
        "success": true,
        "alumni": alumnis,
    })
}

func GetAlumniByYearService(c *fiber.Ctx, db *sql.DB) error {
    alumnis, err := repository.GetAllAlumniByYear(db)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengambil data alumni: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan data alumni dengan tahun lulus 2023",
        "success": true,
        "alumni": alumnis,
    })
}

    func GetAlumniWithYearService(c *fiber.Ctx, db *sql.DB) error {
        alumnis, err := repository.GetAlumniWithYear(db)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Gagal mengambil data alumni: " + err.Error(),
                "success": false,
            })
        }

        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "message": "Berhasil mendapatkan data alumni dengan tahun lulus yang sama dengan bekerjanya",
            "success": true,
            "alumni": alumnis,
        })
    }

func GetAlumniService(c *fiber.Ctx, db *sql.DB) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")

    offset := (page - 1) * limit

    sortByWhitelist := map[string]bool{"id": true, "nama": true, "nim": true, "jurusan": true, "angkatan": true, "tahun_lulus": true}
    if !sortByWhitelist[sortBy] {
        sortBy = "id"
    }
    if strings.ToLower(order) != "desc" {
        order = "asc"
    }


    alumni, err := repository.GetAlumniRepo(db, search, sortBy, order, limit, offset)
     if err != nil {
        fmt.Println("Query error:", err) // untuk cek di console
         return c.Status(500).JSON(fiber.Map{
        "error": err.Error(), // kirim error asli ke client
    })
    }

    total, err := repository.CountAlumniRepo(db, search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to count alumni"})
    }

    response := model.AlumniResponse{
        Data: alumni,
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






