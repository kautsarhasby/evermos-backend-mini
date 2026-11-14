package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"evermos-hasby/ecommerce/internal/config"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RequireAuth(c *fiber.Ctx) error {
	return requireToken(c, false)
}

func RequireRefreshAuth(c *fiber.Ctx) error {
	return requireToken(c, true)
}

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != role {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.Next()
	}
}

func requireToken(c *fiber.Ctx, refresh bool) error {
	header := c.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	claims, err := utils.ValidateToken(c.UserContext(), token)
	if err != nil || claims.UserID == 0 || claims.ID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}
	if config.Redis == nil {
		return c.SendStatus(fiber.StatusServiceUnavailable)
	}
	data, err := config.Redis.Get(c.UserContext(), fmt.Sprintf("session:%d:%s", claims.UserID, claims.ID)).Bytes()
	if err == redis.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}
	if err != nil {
		return c.SendStatus(fiber.StatusServiceUnavailable)
	}
	var session models.RedisSession
	if json.Unmarshal(data, &session) != nil || session.UserID != claims.UserID || (refresh && session.RefreshToken != token) || (!refresh && session.Token != token) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}
	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)
	c.Locals("token", claims)
	return c.Next()
}
