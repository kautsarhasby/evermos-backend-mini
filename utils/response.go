package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendResponseHTTP(c *fiber.Ctx, code int, message string, data interface{}) {
	resp := Response{
		Message: message,
		Data:    data,
	}
	c.Status(code).JSON(resp)
}
