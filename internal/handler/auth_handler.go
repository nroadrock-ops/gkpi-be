package handler

import (
	"strings"
	
	"gkpi-be/internal/domain"
	"gkpi-be/internal/service"
	"gkpi-be/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	type Request struct {
		NamaLengkap  string `json:"nama"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		NomorTelepon string `json:"telepon"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}

	role := "jemaat"
	_, connectToken, connectURL, err := h.service.Register(req.NamaLengkap, req.Email, req.Password, req.NomorTelepon, role)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique constraint") {
			return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Email sudah terdaftar. Silakan login atau gunakan email lain.")
		}
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}

	// Kembalikan token dan link connect Telegram
	return utils.JSONResponse(c, fiber.StatusCreated, true, fiber.Map{
		"connect_token": connectToken,
		"telegram_connect_url": connectURL,
	}, "Registrasi berhasil. Silakan hubungkan akun Telegram Anda.")
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}

	acc, ref, role, err := h.service.VerifyOTP(req.Email, req.Code, "verification")
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, false, nil, err.Error())
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, fiber.Map{
		"access_token":  acc,
		"refresh_token": ref,
		"role":          role,
	}, "OTP Verified and Login successful")
}

func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	type Request struct {
		Email   string `json:"email"`
		Channel string `json:"channel"` // "email" atau "telegram"
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}
	if req.Channel == "" {
		req.Channel = "telegram" // default Telegram
	}

	url, err := h.service.ResendOTP(req.Email, req.Channel, "verification")
	if err != nil {
		if err.Error() == "telegram not connected" {
			return utils.JSONResponse(c, fiber.StatusBadRequest, false, fiber.Map{
				"needs_connection": true,
				"connect_url":      url,
			}, "Telegram belum terhubung. Silakan klik link untuk menghubungkan.")
		}
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, err.Error())
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "OTP resend successful")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}

	acc, ref, role, needsVerif, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, false, nil, "Invalid credentials")
	}

	if needsVerif {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, false, fiber.Map{
			"needs_verification": true,
		}, "Account not verified")
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, fiber.Map{
		"access_token":  acc,
		"refresh_token": ref,
		"role":          role,
	}, "Login successful")
}

func (h *AuthHandler) TelegramStatus(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "token is required")
	}

	connected, err := h.service.CheckTelegramStatus(token)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, err.Error())
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, fiber.Map{"connected": connected}, "Telegram status checked")
}

func (h *AuthHandler) TelegramWebhook(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Jangan block response Telegram (harus kembalikan 200 OK dengan cepat)
	go h.service.ProcessTelegramWebhook(payload)
	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) GetOTPSettings(c *fiber.Ctx) error {
	settings, err := h.service.GetOTPSettings()
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, settings, "OTP settings retrieved")
}

func (h *AuthHandler) UpdateOTPSettings(c *fiber.Ctx) error {
	var req domain.OTPSetting
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}

	// Ambil setting lama, timpa properties-nya
	settings, _ := h.service.GetOTPSettings()
	settings.EmailEnabled = req.EmailEnabled
	settings.TelegramEnabled = req.TelegramEnabled
	settings.OTPExpiryMinutes = req.OTPExpiryMinutes
	settings.MaxAttempts = req.MaxAttempts
	settings.MaxResendPerHour = req.MaxResendPerHour

	if err := h.service.UpdateOTPSettings(settings); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, settings, "OTP settings updated")
}

// Below are the unchanged endpoints
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Token refreshed (mock)")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Logged out successfully")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)
	user, err := h.service.GetMe(uint(userID))
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "User not found")
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, user, "User profile")
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}
	url, err := h.service.ForgotPassword(req.Email)
	if err != nil {
		if err.Error() == "telegram not connected" {
			return utils.JSONResponse(c, fiber.StatusBadRequest, false, fiber.Map{
				"needs_connection": true,
				"connect_url":      url,
			}, "Akun ini belum terhubung ke Telegram. Silakan hubungkan terlebih dahulu.")
		}
		if err.Error() == "email not found" {
			return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "Email tidak terdaftar di sistem kami.")
		}
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, "Failed to process request")
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "OTP telah dikirim ke Telegram Anda")
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		OTP      string `json:"otp"`
		Password string `json:"password"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}
	err := h.service.ResetPassword(req.Email, req.OTP, req.Password)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, err.Error())
	}
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Password berhasil diubah")
}

func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	// Ambil user_id dari JWT token yang sudah diverifikasi middleware
	userIDLocal := c.Locals("user_id")
	if userIDLocal == nil {
		return utils.JSONResponse(c, fiber.StatusForbidden, false, nil, "Unauthorized")
	}
	userID := uint(userIDLocal.(float64))

	type Request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}

	// 1. Pastikan field password HANYA diproses kalau memang dikirim dan tidak kosong
	if req.OldPassword == "" || req.NewPassword == "" {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Password lama dan baru harus diisi")
	}

	// 2. Kirim ke service (di mana hashing dengan bcrypt akan dilakukan)
	if err := h.service.UpdatePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, err.Error())
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Password berhasil diubah")
}
