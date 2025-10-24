package service

import (
	_ "fmt"
	"log"
	_ "os"
	_ "strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/model"
	"github.com/noorfarihaf11/clean-arc/app/repository"
	"github.com/noorfarihaf11/clean-arc/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllAlumniService(c *fiber.Ctx, db *mongo.Database) error {
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
	username := claims.Username

	// Cek role user dari token
	if claims.Role == "admin" {
		log.Printf("Admin %s menambah alumni baru", username)
		userID = nil
	} else {
		log.Printf("User %s (alumni) menambah data dirinya sendiri", username)
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

func UpdateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
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

func DeleteAlumniService(c *fiber.Ctx, db *mongo.Database) error {
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

// func GetAlumniBySalaryService(c *fiber.Ctx, db *mongo.Database) error {
//     alumnis, err := repository.GetAlumniWithHighSalary(db)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "message": "Gagal mengambil data alumni: " + err.Error(),
//             "success": false,
//         })
//     }

//     return c.Status(fiber.StatusOK).JSON(fiber.Map{
//         "message": "Berhasil mendapatkan data alumni dengan gaji > 20 juta",
//         "success": true,
//         "alumni": alumnis,
//     })
// }

// func GetAlumniByYearService(c *fiber.Ctx, db *mongo.Database) error {
//     alumnis, err := repository.GetAllAlumniByYear(db)
//     if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "message": "Gagal mengambil data alumni: " + err.Error(),
//             "success": false,
//         })
//     }

//     return c.Status(fiber.StatusOK).JSON(fiber.Map{
//         "message": "Berhasil mendapatkan data alumni dengan tahun lulus 2023",
//         "success": true,
//         "alumni": alumnis,
//     })
// }

//     func GetAlumniWithYearService(c *fiber.Ctx, db *mongo.Database) error {
//         alumnis, err := repository.GetAlumniWithYear(db)
//         if err != nil {
//             return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//                 "message": "Gagal mengambil data alumni: " + err.Error(),
//                 "success": false,
//             })
//         }

//         return c.Status(fiber.StatusOK).JSON(fiber.Map{
//             "message": "Berhasil mendapatkan data alumni dengan tahun lulus yang sama dengan bekerjanya",
//             "success": true,
//             "alumni": alumnis,
//         })
//     }

// func GetAlumniService(c *fiber.Ctx, db *mongo.Database) error {
//     page, _ := strconv.Atoi(c.Query("page", "1"))
//     limit, _ := strconv.Atoi(c.Query("limit", "10"))
//     sortBy := c.Query("sortBy", "id")
//     order := c.Query("order", "asc")
//     search := c.Query("search", "")

//     offset := (page - 1) * limit

//     sortByWhitelist := map[string]bool{"id": true, "nama": true, "nim": true, "jurusan": true, "angkatan": true, "tahun_lulus": true}
//     if !sortByWhitelist[sortBy] {
//         sortBy = "id"
//     }
//     if strings.ToLower(order) != "desc" {
//         order = "asc"
//     }


//     alumni, err := repository.GetAlumniRepo(db, search, sortBy, order, limit, offset)
//      if err != nil {
//         fmt.Println("Query error:", err) // untuk cek di console
//          return c.Status(500).JSON(fiber.Map{
//         "error": err.Error(), // kirim error asli ke client
//     })
//     }

//     total, err := repository.CountAlumniRepo(db, search)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{"error": "Failed to count alumni"})
//     }

//     response := model.AlumniResponse{
//         Data: alumni,
//         Meta: model.MetaInfo{
//             Page:   page,
//             Limit:  limit,
//             Total:  total,
//             Pages:  (total + limit - 1) / limit,
//             SortBy: sortBy,
//             Order:  order,
//             Search: search,
//         },
//     }

//     return c.JSON(response)
// }






