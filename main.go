package main 

import (
	"log"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Data-Alchemist/doculex-api/config"
	"github.com/Data-Alchemist/doculex-api/database"
	"github.com/Data-Alchemist/doculex-api/routes"
)

func main() {
	app := fiber.New() //initialize the server

	database.ConnectDB() //connect to database
	defer database.DisconnectDB() //disconnect from database

	app.Use(logger.New()) //add logger to track http request 

	routes.SetupEndpoint(app)

	//add setup handler for false routes
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "404 Not Found",
			"status": fiber.StatusNotFound,
		})
	})

	host := config.ConfigHost()
	port := config.ConfigPort()

	fmt.Println("\nServer is running on", host + ":" + port)

	err := app.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}