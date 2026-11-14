package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type TransactionController struct {
	Service interfaces.ITransactionService
}

func (h *TransactionController) Create(c *fiber.Ctx) error {
	var request models.CreateTransactionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(400)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(401)
	}
	tx, err := h.Service.Create(c.UserContext(), userID, request)
	if err != nil {
		utils.SendResponseHTTP(c, 400, err.Error(), nil)
		return nil
	}
	utils.SendResponseHTTP(c, http.StatusCreated, "Success create transaction", tx)
	return nil
}
func (h *TransactionController) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(401)
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	items, err := h.Service.List(c.UserContext(), userID, c.Query("status"), limit, offset)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	utils.SendResponseHTTP(c, 200, "Success", items)
	return nil
}
