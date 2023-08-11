package controllers

import (
	"context"

	"github.com/gofiber/fiber/v2"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Data-Alchemist/doculex-api/database"
	"github.com/Data-Alchemist/doculex-api/models/entity"
)


type UserController interface {
	GetAllUserAccount(c *fiber.Ctx) error
}

type userController struct {}

func NewUserController() UserController {
	return &userController{}
}

func(controller *userController) GetAllUserAccount(c *fiber.Ctx) error {
	database.ConnectDB()
	defer database.DisconnectDB()

	client := database.GetDB()
	collection := database.GetCollection(client, "user")

	var users []entity.User

	cursor, err := collection.Find(context.Background(), options.Find())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get all user account",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	if err = cursor.All(context.Background(), &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get all user account",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get all user account",
		"status": fiber.StatusOK,
		"data": users,
	})
}