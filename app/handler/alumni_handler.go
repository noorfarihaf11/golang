package handler

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/service"
)

type AlumniHandler struct {
	db *sql.DB
}

func NewAlumniHandler(db *sql.DB) *AlumniHandler {
	return &AlumniHandler{db: db}
}

func (h *AlumniHandler) GetAll(c *fiber.Ctx) error {
	return service.GetAllAlumniService(c, h.db)
}

func (h *AlumniHandler) GetByID(c *fiber.Ctx) error {
	return service.GetAlumniByIDService(c, h.db)
}

func (h *AlumniHandler) Create(c *fiber.Ctx) error {
	return service.CreateAlumniService(c, h.db)
}

func (h *AlumniHandler) Update(c *fiber.Ctx) error {
	return service.UpdateAlumniService(c, h.db)
}

func (h *AlumniHandler) Delete(c *fiber.Ctx) error {
	return service.DeleteAlumniService(c, h.db)
}
