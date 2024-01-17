package controllers

import (
	"log"
	"time"
	"context"
	
	"github.com/InnoFours/skin-savvy/database"
	"github.com/InnoFours/skin-savvy/models/entity"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func GetAllUser(c *fiber.Ctx) error {
	client, err := database.FirestoreConnection()
	if err != nil {
		log.Fatal("error connecting to firestore", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	: "error connecting to firestore",
			"status"	: fiber.StatusBadRequest,
			"error"		: err.Error(),
		})
	}
	defer client.Close()

	var users []entity.User
	itr := client.Collection("skinsavvy-user").Documents(context.Background())
	for {
		doc, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate the list of posts: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message"	: "failed to get all users",
				"status"	: fiber.StatusInternalServerError,
				"error"		: err.Error(),
			})
		}

		user := entity.User{
			ID			: doc.Data()["ID"].(string),
			Fullname	: doc.Data()["Fullname"].(string),
			Email		: doc.Data()["Email"].(string),
			Password	: doc.Data()["Password"].(string),
			Photo		: doc.Data()["Photo"].(string),
			CreatedAt	: doc.Data()["CreatedAt"].(time.Time),
			UpdatedAt	: doc.Data()["UpdatedAt"].(time.Time),
		}
		users = append(users, user)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message"	: "success retrieving all users",
		"data"		: users,
		"count"		: len(users),
	})
	
}