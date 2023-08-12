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
	r.Get("doculex-api/v0.1/get/user", userController.GetAllUserAccount) //endpoint to get all user data
	r.Get("doculex-api/v0.1/get/user/:id", userController.GetOneUserAccount) //endpoint to get one user data by objectid
	r.Post("doculex-api/v0.1/post/user/register", userController.CreateUserAccount) //endpoint for user account register
	r.Post("doculex-api/v0.1/post/user/login", userController.UserLoginValidator) //endpoint for user account login
	// r.Put("doculex-api/v0.1/put/user/update/:id", userController.UpdateUserInfo) //endpoint for update user data information by objectcid
	// r.Delete("doculex-api/v0.1/delete/user/remove/:id", userController.DeleteUserAccount) //endpoint to remove user data by objectid
}