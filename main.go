package main 

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/joho/godotenv"
)

func main() {
	app := fiber.New() //initialize the server

	app.Use(logger.New()) //add logger to track http request 

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	err := app.Listen("localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
}