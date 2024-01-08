package main

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/api/option"

	"github.com/InnoFours/skin-savvy/auth"
	"github.com/InnoFours/skin-savvy/config"
	"github.com/InnoFours/skin-savvy/database"
	"github.com/InnoFours/skin-savvy/middleware"
	"github.com/InnoFours/skin-savvy/models/entity"
	"github.com/InnoFours/skin-savvy/routes"
)

func main() {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = location

	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect with database")
	}
	defer conn.Close()

	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalln("Error getting Auth client: ", err)
	}

	conn.AutoMigrate(&entity.User{})

	authService := &auth.AuthService{
		DB			: conn,
		FireAuth	: authClient,
	}

	authToken := middleware.NewFireAuthMiddleware(authService)

	server := fiber.New()

	// server.Use(func(c *fiber.Ctx) error {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"message": "404 Not Found",
	// 		"status": fiber.StatusNotFound,
	// 	})
	// })

	server.Use(logger.New())

	server.Use(middleware.CORSMiddleware())

	routes.SetupEndpoint(server, authService, authToken)

	host := config.ConfigHost()
	port := config.ConfigPort()

	err = server.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}