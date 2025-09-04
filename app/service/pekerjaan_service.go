package service

import (
	"database/sql"
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/noorfarihaf11/clean-arc/app/repository"
    "github.com/noorfarihaf11/clean-arc/app/model"
    "strconv"
)

func GetAllJobService(c *fiber.Ctx, db *sql.DB) error {
    key := c.Get("X-API-KEY") 
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
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
    key := c.Get("X-API-KEY") 
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
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
    key := c.Get("X-API-KEY")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
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
    key := c.Get("X-API-KEY")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
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
    // Validasi API Key
    key := c.Get("X-API-KEY")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
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
    // Validasi API Key
    key := c.Get("X-API-KEY")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "Key tidak valid",
            "success": false,
        })
    }

    //  Ambil ID dari URL
    id := c.Params("id")

    //  Hapus alumni dari DB
    err := repository.DeleteJob(db, id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus pekerjaan alumni: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Pekerjaan alumni berhasil dihapus",
        "success": true,
    })
}
