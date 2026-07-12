package handler

import (
	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type RenunganHandler struct {
	service service.RenunganService
}

func NewRenunganHandler(service service.RenunganService) *RenunganHandler {
	return &RenunganHandler{service}
}

func (h *RenunganHandler) GetToday(c *fiber.Ctx) error {
	renungan, err := h.service.GetToday()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "Renungan for today not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, renungan, "Today's renungan")
}

func (h *RenunganHandler) Create(c *fiber.Ctx) error {
	var renungan domain.Renungan
	if err := c.BodyParser(&renungan); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}
	if err := h.service.Create(&renungan); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, renungan, "Renungan created")
}
