package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gkpi-be/internal/utils"
)

// Protected verifies the JWT token
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return utils.JSONResponse(c, fiber.StatusUnauthorized, false, nil, "Missing or invalid token")
		}

		tokenString := authHeader[7:]
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return utils.JSONResponse(c, fiber.StatusUnauthorized, false, nil, "Invalid or expired token")
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])
		if claims["jemaat_id"] != nil {
			c.Locals("jemaat_id", claims["jemaat_id"])
		}

		return c.Next()
	}
}

// AdminOnly requires the user role to be admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return utils.JSONResponse(c, fiber.StatusForbidden, false, nil, "Admin access required")
		}
		return c.Next()
	}
}

// RequireRole restricts access to users with a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != requiredRole {
			return utils.JSONResponse(c, fiber.StatusForbidden, false, nil, "Access forbidden for your role")
		}
		return c.Next()
	}
}
