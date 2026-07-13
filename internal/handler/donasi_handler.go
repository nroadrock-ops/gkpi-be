package handler

import (
	"strconv"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type DonasiHandler struct {
	service service.DonasiService
}

func NewDonasiHandler(service service.DonasiService) *DonasiHandler {
	return &DonasiHandler{service}
}

func (h *DonasiHandler) Create(c *fiber.Ctx) error {
	var req domain.Donasi
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}

	donasi, err := h.service.Create(&req)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, donasi, "Donasi initiated")
}

func (h *DonasiHandler) GetStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	donasi, err := h.service.GetStatus(uint(id))
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "Donasi not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, donasi, "Donasi status")
}

func (h *DonasiHandler) Webhook(c *fiber.Ctx) error {
	type WebhookRequest struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"` // settlement, pending, cancel
	}
	var req WebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}

	if err := h.service.ProcessWebhook(req.TransactionID, req.Status); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Webhook processed")
}

func (h *DonasiHandler) GetMeDonasi(c *fiber.Ctx) error {
	jemaatIDLocal := c.Locals("jemaat_id")
	if jemaatIDLocal == nil {
		return utils.JSONResponse(c, fiber.StatusForbidden, false, nil, "Only jemaat can access this")
	}
	jemaatID := uint(jemaatIDLocal.(float64))

	donasis, err := h.service.GetByJemaatID(jemaatID)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, donasis, "Donasi history")
}
