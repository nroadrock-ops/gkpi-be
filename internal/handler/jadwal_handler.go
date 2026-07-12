package handler

import (
	"strconv"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type JadwalHandler struct {
	service service.JadwalService
}

func NewJadwalHandler(service service.JadwalService) *JadwalHandler {
	return &JadwalHandler{service}
}

func (h *JadwalHandler) GetAll(c *fiber.Ctx) error {
	jadwals, err := h.service.GetAll()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, jadwals, "Jadwal list")
}

func (h *JadwalHandler) Create(c *fiber.Ctx) error {
	var jadwal domain.JadwalIbadah
	if err := c.BodyParser(&jadwal); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Create(&jadwal); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, jadwal, "Jadwal created")
}

func (h *JadwalHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var jadwal domain.JadwalIbadah
	if err := c.BodyParser(&jadwal); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Update(uint(id), &jadwal); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Jadwal updated")
}

func (h *JadwalHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Jadwal deleted")
}
