package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist/doculex-api/controllers"
)

func SetupEndpoint(r *fiber.App) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to DoculexAI version 1.0 public API👍",
			"status": fiber.StatusOK,
		})
	})

	//initiate the reservoir variable
	userController := controllers.NewUserController()

	//user endpoint
	r.Get("doculex-api/v0.1/get/user", userController.GetAllUserAccount)
	r.Get("doculex-api/v0.1/get/user/:id", userController.GetOneUserAccount)
}