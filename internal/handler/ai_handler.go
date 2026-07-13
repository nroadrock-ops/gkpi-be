package handler

import (
	"strconv"

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
	res, err := h.client.ForwardRequest("/ocr/extract", payload)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "Digitize successful")
}

func (h *AIHandler) AutoTagGaleri(c *fiber.Ctx) error {
	var payload map[string]interface{}
	c.BodyParser(&payload)
	res, err := h.client.ForwardRequest("/galeri/classify", payload)
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
	res, err := h.client.ForwardRequest("/absensi/recognize", payload)
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

func (h *AIHandler) GetMeAbsensi(c *fiber.Ctx) error {
	jemaatIDLocal := c.Locals("jemaat_id")
	if jemaatIDLocal == nil {
		return utils.JSONResponse(c, fiber.StatusForbidden, false, nil, "Only jemaat can access this")
	}
	jemaatID := uint(jemaatIDLocal.(float64))

	// Convert to string for the API path
	jemaatIDStr := strconv.Itoa(int(jemaatID))
	res, err := h.client.ForwardRequest("/absensi/history/" + jemaatIDStr, nil)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, res, "Attendance history retrieved")
}
