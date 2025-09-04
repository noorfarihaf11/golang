package handler

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/noorfarihaf11/clean-arc/app/service"
)

type JobHandler struct {
	db *sql.DB
}

func NewJobHandler(db *sql.DB) *JobHandler {
	return &JobHandler{db: db}
}

func (h *JobHandler) GetAll(c *fiber.Ctx) error {
	return service.GetAllJobService(c, h.db)
}

func (h *JobHandler) GetByID(c *fiber.Ctx) error {
	return service.GetJobByIDService(c, h.db)
}

func (h *JobHandler) Create(c *fiber.Ctx) error {
	return service.CreateJobService(c, h.db)
}

func (h *JobHandler) Update(c *fiber.Ctx) error {
	return service.UpdateJobService(c, h.db)
}

func (h *JobHandler) Delete(c *fiber.Ctx) error {
	return service.DeleteJobService(c, h.db)
}
