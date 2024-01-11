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

	firebaseAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalln("Error getting Auth client: ", err)
	}

	conn.AutoMigrate(&entity.User{})

	authService := &auth.AuthService{
		DB			: conn,
		FireAuth	: firebaseAuth,
	}

	server := fiber.New()

	server.Use(logger.New())

	server.Use(middleware.CORSMiddleware())

	server.Use(func(c *fiber.Ctx) error {
		c.Locals("firebaseAuth", firebaseAuth)
		return c.Next()
	})

	routes.SetupEndpoint(server, authService)

	// server.Use(middleware.TokenValidator)

	host := config.ConfigHost()
	port := config.ConfigPort()

	err = server.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}