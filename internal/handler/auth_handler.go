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
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, false, nil, "Invalid request body")
	}

	if req.Role == "" {
		req.Role = "jemaat" // default role
	}

	user, err := h.service.Register(req.Email, req.Password, req.Role)
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

	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusUnauthorized, false, nil, "Invalid credentials")
	}

	return utils.JSONResponse(c, fiber.StatusOK, true, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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
