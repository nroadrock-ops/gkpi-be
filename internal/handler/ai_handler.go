package handler

import (
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	client service.AIServiceClient
}

func NewAIHandler(client service.AIServiceClient) *AIHandler {
	return &AIHandler{client}
}

func (h *AIHandler) DigitizeArsip(c *fiber.Ctx) error {
	var payload map[string]interface{}
	c.BodyParser(&payload)
	res, err := h.client.ForwardRequest("/arsip/digitize", payload)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "Digitize successful")
}

func (h *AIHandler) AutoTagGaleri(c *fiber.Ctx) error {
	var payload map[string]interface{}
	c.BodyParser(&payload)
	res, err := h.client.ForwardRequest("/galeri/auto-tag", payload)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "Auto tag successful")
}

func (h *AIHandler) EnrollAbsensi(c *fiber.Ctx) error {
	var payload map[string]interface{}
	c.BodyParser(&payload)
	res, err := h.client.ForwardRequest("/absensi/enroll", payload)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "Enroll successful")
}

func (h *AIHandler) CheckInAbsensi(c *fiber.Ctx) error {
	var payload map[string]interface{}
	c.BodyParser(&payload)
	res, err := h.client.ForwardRequest("/absensi/check-in", payload)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "Check-in successful")
}

func (h *AIHandler) HistoryAbsensi(c *fiber.Ctx) error {
	jemaatID := c.Params("jemaat_id")
	res, err := h.client.ForwardRequest("/absensi/history/" + jemaatID, nil)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "History retrieved")
}
