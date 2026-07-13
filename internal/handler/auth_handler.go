package handler

import (
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

	// Otomatis role jemaat
	role := "jemaat"

	user, err := h.service.Register(req.NamaLengkap, req.Email, req.Password, req.NomorTelepon, role)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, err.Error())
	}

	return utils.JSONResponse(c, fiber.StatusCreated, true, user, "Registration successful")
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

	accessToken, refreshToken, role, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, false, nil, "Invalid credentials")
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"role":          role,
	}, "Login successful")
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// TODO: implement refresh token logic (validate refresh token, issue new access token)
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Token refreshed (mock)")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: implement logout logic (e.g. invalidate token in redis/db if needed)
	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "Logged out successfully")
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64) // JWT mapclaims unmarshals numbers to float64
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

	err := h.service.ForgotPassword(req.Email)
	if err != nil {
		if err.Error() == "email not found" {
			// Demi keamanan, lebih baik tetap beri respons OK agar tidak membocorkan keberadaan email
			return utils.JSONResponse(c, fiber.StatusNotFound, false, nil, "Email tidak terdaftar di sistem kami.")
		}
		return utils.JSONResponse(c, fiber.StatusInternalServerError, false, nil, "Failed to process request")
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, nil, "OTP telah dikirim ke email Anda")
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
