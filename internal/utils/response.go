package utils

import "github.com/gofiber/fiber/v2"

// Response format standar
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func JSONResponse(c *fiber.Ctx, status int, success bool, data interface{}, message string) error {
	return c.Status(status).JSON(Response{
		Success: success,
		Data:    data,
		Message: message,
	})
}
