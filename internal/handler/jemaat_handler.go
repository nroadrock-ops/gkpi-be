package handler

import (
	"strconv"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type JemaatHandler struct {
	service service.JemaatService
}

func NewJemaatHandler(service service.JemaatService) *JemaatHandler {
	return &JemaatHandler{service}
}

func (h *JemaatHandler) GetAll(c *fiber.Ctx) error {
	jemaats, err := h.service.GetAll()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, jemaats, "Jemaat list")
}

func (h *JemaatHandler) Create(c *fiber.Ctx) error {
	var jemaat domain.Jemaat
	if err := c.BodyParser(&jemaat); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Create(&jemaat); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, jemaat, "Jemaat created")
}

func (h *JemaatHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	jemaat, err := h.service.GetByID(uint(id))
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "Jemaat not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, jemaat, "Jemaat detail")
}

func (h *JemaatHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var jemaat domain.Jemaat
	if err := c.BodyParser(&jemaat); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Update(uint(id), &jemaat); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Jemaat updated")
}

func (h *JemaatHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Jemaat deleted")
}
