package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckController struct {
	HealthCheckService interfaces.IHealthCheckServices
}

func (h *HealthCheckController) HealthCheck(c *fiber.Ctx) {
	msg, err := h.HealthCheckService.HealthCheckServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError)
		return
	}

	utils.SendResponseHTTP(c, http.StatusOK, msg, nil)
}
