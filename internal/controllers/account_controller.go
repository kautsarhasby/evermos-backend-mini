package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AccountController struct{ Service interfaces.IAccountService }

func (h *AccountController) List(c *fiber.Ctx) error {
	users, err := h.Service.List(c.UserContext(), c.Locals("user_id").(uint))
	if err != nil {
		return c.SendStatus(http.StatusForbidden)
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success", users)
	return nil
}

func (h *AccountController) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	user, err := h.Service.Get(c.UserContext(), c.Locals("user_id").(uint), uint(id))
	if err != nil || user == nil {
		return c.SendStatus(http.StatusNotFound)
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success", user)
	return nil
}

func (h *AccountController) Create(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if err := h.Service.Create(c.UserContext(), c.Locals("user_id").(uint), &user); err != nil {
		utils.SendResponseHTTP(c, http.StatusBadRequest, err.Error(), nil)
		return nil
	}
	user.Password = ""
	utils.SendResponseHTTP(c, http.StatusCreated, "Success create account", user)
	return nil
}

func (h *AccountController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var request models.User
	if err != nil || c.BodyParser(&request) != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	user, err := h.Service.Update(c.UserContext(), c.Locals("user_id").(uint), uint(id), &request)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success update account", user)
	return nil
}

func (h *AccountController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if err := h.Service.Delete(c.UserContext(), c.Locals("user_id").(uint), uint(id)); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	return c.SendStatus(http.StatusNoContent)
}
