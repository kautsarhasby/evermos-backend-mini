package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type StoreController struct {
	StoreService interfaces.IStoreService
}

func (h *StoreController) CreateStore(c *fiber.Ctx) error {
	var request models.Store

	if err := c.BodyParser(&request); err != nil {
		customMsg := utils.GetCustomError(err, &request)
		utils.SendResponseHTTP(c, http.StatusBadRequest, customMsg, nil)
		return nil
	}

	userID := c.Locals("user_id").(uint)

	response, err := h.StoreService.CreateStore(c.UserContext(), request, userID)
	if err != nil {
		utils.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return nil
	}

	utils.SendResponseHTTP(c, http.StatusCreated, "Success create store", response)
	return nil
}

func (h *StoreController) GetStore(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	store, err := h.StoreService.GetStore(c.UserContext(), userID)
	if err != nil || store == nil {
		return c.SendStatus(http.StatusNotFound)
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success", store)
	return nil
}

func (h *StoreController) UpdateStore(c *fiber.Ctx) error {
	var request models.Store
	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	store, err := h.StoreService.UpdateStore(c.UserContext(), request, userID)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success update store", store)
	return nil
}
