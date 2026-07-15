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
	judul := c.FormValue("judul")
	file, err := c.FormFile("file")
	if err != nil {
		// Fallback ke JSON payload (agar tidak breaking jika frontend masih pakai cara lama)
		type Request struct {
			Judul   string `json:"judul"`
			FileURL string `json:"file_url"`
		}
		var req Request
		if err := c.BodyParser(&req); err == nil && req.FileURL != "" {
			galeri, err := h.service.UploadAndSave(req.Judul, req.FileURL)
			if err != nil {
				return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
			}
			return utils.JSONResponse(c, fiber.StatusCreated, true, galeri, "Galeri uploaded (via JSON)")
		}
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "File or file_url is required")
	}

	if judul == "" {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Judul is required")
	}

	// Proses upload file ke Supabase (bucket: gkpi)
	// Asumsi bucket name adalah "gkpi" sesuai dengan URL di Postman
	fileURL, err := utils.UploadToSupabaseStorage(file, "gkpi")
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, "Failed to upload to Supabase: "+err.Error())
	}

	galeri, err := h.service.UploadAndSave(judul, fileURL)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusCreated, true, galeri, "Galeri uploaded successfully")
}

func (h *GaleriHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Galeri deleted")
}
