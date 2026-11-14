package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AddressController struct {
	AddressService interfaces.IAddressService
}

func (h *AddressController) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(401)
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	items, err := h.AddressService.GetAddresses(c.UserContext(), userID, c.Query("search"), limit, offset)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	utils.SendResponseHTTP(c, 200, "Success", items)
	return nil
}

func (h *AddressController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var request models.Address
	if err != nil || c.BodyParser(&request) != nil {
		return c.SendStatus(400)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(401)
	}
	if err := h.AddressService.UpdateAddress(c.UserContext(), request, uint(id), userID); err != nil {
		return c.SendStatus(404)
	}
	utils.SendResponseHTTP(c, 200, "Success update address", nil)
	return nil
}

func (h *AddressController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.SendStatus(400)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(401)
	}
	if err := h.AddressService.DeleteAddress(c.UserContext(), uint(id), userID); err != nil {
		return c.SendStatus(404)
	}
	return c.SendStatus(204)
}

func (h *AddressController) CreateAddress(c *fiber.Ctx) error {
	var request models.Address

	if err := c.BodyParser(&request); err != nil {
		customMsg := utils.GetCustomError(err, &request)
		utils.SendResponseHTTP(c, http.StatusBadRequest, customMsg, nil)
		return nil
	}

	userID := c.Locals("user_id").(uint)

	response, err := h.AddressService.CreateAddress(c.UserContext(), request, userID)
	if err != nil {
		utils.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return nil
	}

	utils.SendResponseHTTP(c, http.StatusCreated, "Success create address", response)
	return nil
}
