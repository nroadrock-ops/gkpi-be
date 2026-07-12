package handler

import (
	"strconv"

	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type GaleriHandler struct {
	service service.GaleriService
}

func NewGaleriHandler(service service.GaleriService) *GaleriHandler {
	return &GaleriHandler{service}
}

func (h *GaleriHandler) GetAll(c *fiber.Ctx) error {
	galeris, err := h.service.GetAll()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, galeris, "Galeri list")
}

func (h *GaleriHandler) Upload(c *fiber.Ctx) error {
	type Request struct {
		Judul   string `json:"judul"`
		FileURL string `json:"file_url"` // Dummy for now, assuming frontend uploaded it directly or via proxy
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid body")
	}

	galeri, err := h.service.UploadAndSave(req.Judul, req.FileURL)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, galeri, "Galeri uploaded")
}

func (h *GaleriHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Galeri deleted")
}
