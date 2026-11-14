package controllers

import (
	"evermos-hasby/ecommerce/internal/interfaces"
	"evermos-hasby/ecommerce/internal/models"
	"evermos-hasby/ecommerce/utils"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService interfaces.IUserService
}

func (h *AuthController) Register(c *fiber.Ctx) error {
	var (
		request models.User
	)

	if err := c.BodyParser(&request); err != nil {
		customMsg := utils.GetCustomError(err, &request)
		utils.SendResponseHTTP(c, http.StatusBadRequest, customMsg, nil)
		return nil
	}

	user, err := h.AuthService.Register(c.UserContext(), &request)
	if err != nil {
		utils.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return nil
	}

	utils.SendResponseHTTP(c, http.StatusOK, "Success register account and create "+user.Username+" Store", nil)
	return nil
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	var (
		request models.LoginRequest
	)

	if err := c.BodyParser(&request); err != nil {
		customMsg := utils.GetCustomError(err, &request)
		utils.SendResponseHTTP(c, http.StatusBadRequest, customMsg, nil)
		return nil
	}

	response, err := h.AuthService.Login(c.UserContext(), request)
	if err != nil {
		utils.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return nil
	}

	utils.SendResponseHTTP(c, http.StatusOK, "Success Login", response)
	return nil
}

func (h *AuthController) Refresh(c *fiber.Ctx) error {

	var (
		log      = utils.Logger
		rawToken = c.Get("Authorization")
	)
	refreshToken := strings.TrimPrefix(rawToken, "Bearer ")
	refreshToken = strings.TrimSpace(refreshToken)

	claim := c.Locals("token")
	if claim == nil {
		log.Error("failed to get claim in context")
		utils.SendResponseHTTP(c, http.StatusInternalServerError, "Failed to get claim in context", nil)
		return nil
	}
	tokenClaim, ok := claim.(*models.TokenClaims)
	if !ok {
		log.Error("failed to parse claim in tokenClaim")
		utils.SendResponseHTTP(c, http.StatusInternalServerError, "Failed to parse claim in tokenClaim", nil)
		return nil
	}

	response, err := h.AuthService.Refresh(c.UserContext(), refreshToken, *tokenClaim)
	if err != nil {
		log.Error("Failed on refresh token service")
		utils.SendResponseHTTP(c, http.StatusInternalServerError, "Failed on refresh token service", nil)
		return nil
	}

	utils.SendResponseHTTP(c, http.StatusCreated, "Success refreshed token", response)
	return nil
}

func (h *AuthController) Logout(c *fiber.Ctx) error {
	var (
		log      = utils.Logger
		rawToken = c.Get("Authorization")
	)
	token := strings.TrimPrefix(rawToken, "Bearer ")
	token = strings.TrimSpace(token)

	err := h.AuthService.Logout(c.UserContext(), token)
	if err != nil {
		log.Error("failed to Logout", err)
		utils.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return nil
	}

	utils.SendResponseHTTP(c, http.StatusOK, "Successfully Logout", nil)
	return nil
}

func (h *AuthController) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	user, err := h.AuthService.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success", user)
	return nil
}

func (h *AuthController) UpdateProfile(c *fiber.Ctx) error {
	var request models.User
	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}
	user, err := h.AuthService.UpdateProfile(c.UserContext(), userID, &request)
	if err != nil {
		utils.SendResponseHTTP(c, http.StatusBadRequest, err.Error(), nil)
		return nil
	}
	utils.SendResponseHTTP(c, http.StatusOK, "Success update profile", user)
	return nil
}
