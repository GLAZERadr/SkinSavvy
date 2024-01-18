package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	
	"github.com/InnoFours/skin-savvy/config"
	"github.com/InnoFours/skin-savvy/middleware"
	"github.com/InnoFours/skin-savvy/routes"
)

func main() {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err)
	}
	time.Local = location

	server := fiber.New()

	server.Use(logger.New())

	server.Use(middleware.CORSMiddleware())

	routes.SetupEndpoint(server)

	host := config.ConfigHost()
	port := config.ConfigPort()

	err = server.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}