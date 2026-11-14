package controllers

import (
	"errors"
	"net/http"

	"evermos-hasby/ecommerce/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

type UploadController struct {
	Service interfaces.IUploadService
}

func (h *UploadController) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "File is required"})
	}
	if file.Size > 5<<20 {
		return c.Status(http.StatusRequestEntityTooLarge).JSON(fiber.Map{"message": "File exceeds 5 MB"})
	}
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	name, err := h.Service.Upload(c.UserContext(), file.Filename, src)
	if errors.Is(err, interfaces.ErrUploadTooLarge) {
		return c.Status(http.StatusRequestEntityTooLarge).JSON(fiber.Map{"message": err.Error()})
	}
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{"filename": name})
}
