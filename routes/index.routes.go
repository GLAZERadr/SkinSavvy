package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/InnoFours/skin-savvy/controllers"
)

func SetupEndpoint(r *fiber.App) {

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to SkinSavvy Beta Version 1.0 Public API👍",
			"status":  fiber.StatusOK,
		})
	})

	skinsavvy := r.Group("/skinsavvy-api/v0.1")

	//public endpoint
	skinsavvy.Post("post/predict", controllers.SkinDetection)
	skinsavvy.Post("post/recommendation", controllers.SkincareRec)
	skinsavvy.Post("sessions/oauth/google", controllers.OauthSignUp)
	skinsavvy.Get("get/all-users", controllers.GetAllUser)
	skinsavvy.Get("get/used-product/:id", controllers.GetUserUsedProduct)
}
