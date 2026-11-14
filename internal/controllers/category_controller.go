package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct{ CategoryService interfaces.ICategoryService }

func (h *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var request models.Category
	if err := c.BodyParser(&request); err != nil || request.CategoryName == "" {
		utils.SendResponseHTTP(c, http.StatusBadRequest, "Category name is required", nil)
		return nil
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	category, err := h.CategoryService.CreateCategory(c.UserContext(), request.CategoryName, userID)
	if err != nil {
		utils.SendResponseHTTP(c, http.StatusForbidden, err.Error(), nil)
		return nil
	}
	utils.SendResponseHTTP(c, http.StatusCreated, "Success create category", category)
	return nil
}

func (h *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	var request struct {
		Name string `json:"name"`
	}
	if err != nil || c.BodyParser(&request) != nil || request.Name == "" {
		return c.SendStatus(http.StatusBadRequest)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	if err := h.CategoryService.UpdateCategory(c.UserContext(), uint(id), request.Name, userID); err != nil {
		utils.SendResponseHTTP(c, http.StatusForbidden, err.Error(), nil)
		return nil
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success update category", nil)
	return nil
}

func (h *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	userID, ok := c.Locals("user_id").(uint)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	if err := h.CategoryService.DeleteCategory(c.UserContext(), uint(id), userID); err != nil {
		utils.SendResponseHTTP(c, http.StatusForbidden, err.Error(), nil)
		return nil
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success delete category", nil)
	return nil
}
