package controllers

import (
	"github.com/InnoFours/skin-savvy/auth"
	"github.com/InnoFours/skin-savvy/models/entity"
	"github.com/gofiber/fiber/v2"
)

type UserAuthController struct {
	userAuthController *auth.AuthService
}

func NewUserController(userAuthController *auth.AuthService) *UserAuthController {
	return &UserAuthController{userAuthController}
}

func(s *UserAuthController) GetAllUser(c *fiber.Ctx) error {
	var users []entity.User

	if err := s.userAuthController.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message"	: "No user found",
			"status"	: fiber.StatusNoContent,
			"error"		: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data"		:  users,
		"status"	: fiber.StatusOK,
	})
}