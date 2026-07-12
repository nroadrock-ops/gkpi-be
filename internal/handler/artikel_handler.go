package handler

import (
	"strconv"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ArtikelHandler struct {
	service service.ArtikelService
}

func NewArtikelHandler(service service.ArtikelService) *ArtikelHandler {
	return &ArtikelHandler{service}
}

func (h *ArtikelHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	artikels, err := h.service.GetAll(page, limit)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, artikels, "Artikel list")
}

func (h *ArtikelHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	artikel, err := h.service.GetBySlug(slug)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "Artikel not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, artikel, "Artikel detail")
}

func (h *ArtikelHandler) Create(c *fiber.Ctx) error {
	var artikel domain.Artikel
	if err := c.BodyParser(&artikel); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Create(&artikel); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, artikel, "Artikel created")
}

func (h *ArtikelHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var artikel domain.Artikel
	if err := c.BodyParser(&artikel); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Update(uint(id), &artikel); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Artikel updated")
}

func (h *ArtikelHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Artikel deleted")
}
