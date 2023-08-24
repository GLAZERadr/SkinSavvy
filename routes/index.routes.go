package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist/doculex-api/controllers"
)

func SetupEndpoint(r *fiber.App, jwtMiddleware fiber.Handler) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to DoculexAI version 1.0 public API👍",
			"status": fiber.StatusOK,
		})
	})

	//initiate the reservoir variable
	userController := controllers.NewUserController()

	//user endpoint
	//public endpoint
	r.Post("doculex-api/v0.1/post/user/register", userController.CreateUserAccount) //endpoint for user account register
	r.Post("doculex-api/v0.1/post/user/login", userController.UserLoginValidator) //endpoint for user account login

	//protected method
	protected := r.Group("doculex-api/v0.1/protected")
	protected.Use(jwtMiddleware)

	//protected endpoint for use controller
	protected.Get("get/user", userController.GetAllUserAccount) //endpoint to get all user data
	protected.Get("get/user/:id", userController.GetOneUserAccount) //endpoint to get one user data by objectid
	protected.Put("put/user/update/:id", userController.UpdateUserInfo) //endpoint for update user data information by objectcid
	protected.Delete("delete/user/remove/:id", userController.DeleteUserAccount) //endpoint to remove user data by objectid
}