package service

import (
	"database/sql"
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/noorfarihaf11/clean-arc/app/repository"
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
    key := c.Params("key")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
            "success": false,
        })
    }

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

