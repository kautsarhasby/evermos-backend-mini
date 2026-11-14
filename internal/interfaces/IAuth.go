package interfaces

import "github.com/gofiber/fiber/v2"

type IAuthController interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}
