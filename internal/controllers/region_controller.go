package controllers

import (
	"net/http"
	"strings"

	"evermos-hasby/ecommerce/internal/interfaces"

	"github.com/gofiber/fiber/v2"
)

type RegionController struct {
	RegionService interfaces.IRegionService
}

func (h *RegionController) GetListProvincies(c *fiber.Ctx) error {
	data, err := h.RegionService.GetProvinces(c.UserContext())
	return regionResponse(c, data, err)
}

func (h *RegionController) GetListCities(c *fiber.Ctx) error {
	data, err := h.RegionService.GetCitiesByProvince(c.UserContext(), c.Params("prov_id"))
	return regionResponse(c, data, err)
}

func (h *RegionController) GetDetailProvince(c *fiber.Ctx) error {
	data, err := h.RegionService.GetProvinceDetail(c.UserContext(), c.Params("prov_id"))
	return regionResponse(c, data, err)
}

func (h *RegionController) GetDetailCity(c *fiber.Ctx) error {
	data, err := h.RegionService.GetCityDetail(c.UserContext(), c.Params("city_id"))
	return regionResponse(c, data, err)
}

func regionResponse(c *fiber.Ctx, data any, err error) error {
	if err != nil {
		status := http.StatusBadGateway
		message := "Failed to fetch region data"

		if strings.Contains(err.Error(), "cannot be empty") {
			status = http.StatusBadRequest
			message = err.Error()
		}

		return c.Status(status).JSON(fiber.Map{
			"status":  false,
			"message": message,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Success",
		"data":    data,
	})
}
