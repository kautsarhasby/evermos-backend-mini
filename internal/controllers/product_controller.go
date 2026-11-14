package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type ProductController struct{ Service interfaces.IProductService }

func (h *ProductController) Create(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	if err := h.Service.Create(c.UserContext(), &product, userID); err != nil {
		utils.SendResponseHTTP(c, 400, err.Error(), nil)
		return nil
	}
	utils.SendResponseHTTP(c, http.StatusCreated, "Success create product", product)
	return nil
}
func (h *ProductController) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	products, err := h.Service.List(c.UserContext(), userID, c.Query("search"), c.Query("sort"), limit, offset)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	utils.SendResponseHTTP(c, 200, "Success", products)
	return nil
}
func (h *ProductController) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	userID, ok := c.Locals("user_id").(uint)
	if err != nil {
		return c.SendStatus(400)
	}
	if !ok {
		return c.SendStatus(401)
	}
	p, err := h.Service.Get(c.UserContext(), uint(id), userID)
	if err != nil {
		return c.SendStatus(404)
	}
	utils.SendResponseHTTP(c, 200, "Success", p)
	return nil
}
func (h *ProductController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var p models.Product
	if err != nil || c.BodyParser(&p) != nil {
		return c.SendStatus(400)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(401)
	}
	p.ID = uint(id)
	if err := h.Service.Update(c.UserContext(), &p, userID); err != nil {
		return c.SendStatus(404)
	}
	utils.SendResponseHTTP(c, 200, "Success update product", p)
	return nil
}
func (h *ProductController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	userID, ok := c.Locals("user_id").(uint)
	if err != nil {
		return c.SendStatus(400)
	}
	if !ok {
		return c.SendStatus(401)
	}
	if err := h.Service.Delete(c.UserContext(), uint(id), userID); err != nil {
		return c.SendStatus(404)
	}
	return c.SendStatus(204)
}
