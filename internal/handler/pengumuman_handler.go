package handler

import (
	"strconv"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type PengumumanHandler struct {
	service service.PengumumanService
}

func NewPengumumanHandler(service service.PengumumanService) *PengumumanHandler {
	return &PengumumanHandler{service}
}

func (h *PengumumanHandler) GetAll(c *fiber.Ctx) error {
	pengumumans, err := h.service.GetAll()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, pengumumans, "Pengumuman list")
}

func (h *PengumumanHandler) Create(c *fiber.Ctx) error {
	var pengumuman domain.Pengumuman
	if err := c.BodyParser(&pengumuman); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Create(&pengumuman); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, pengumuman, "Pengumuman created")
}

func (h *PengumumanHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Pengumuman deleted")
}
