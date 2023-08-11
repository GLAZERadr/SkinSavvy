package main 

import (
	"log"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Data-Alchemist/doculex-api/config"
)

func main() {
	app := fiber.New() //initialize the server

	app.Use(logger.New()) //add logger to track http request 

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	host := config.ConfigHost()
	port := config.ConfigPort()

	fmt.Println("\nServer is running on", host + ":" + port)

	err := app.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}