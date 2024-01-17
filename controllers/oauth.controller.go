package controllers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/InnoFours/skin-savvy/database"
	"github.com/InnoFours/skin-savvy/helper"
	"github.com/InnoFours/skin-savvy/middleware"
	"github.com/InnoFours/skin-savvy/models/entity"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func OauthSignUp(c *fiber.Ctx) error {
	code := c.Query("code")

	if code == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	: "you are unauthorized accessing this site",
			"status"	: fiber.StatusUnauthorized,
		})
	}

	tokenRes, err := middleware.GetGoogleOauthToken(code)
	if err != nil {
		log.Fatal("error acquire google oauth token: ", err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message"	: "can't get google oauth token",
			"status"	: fiber.StatusBadGateway,
			"error"		: err.Error(),
		})
	}

	googleUser, err := middleware.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)
	if err != nil {
		log.Fatal("error acquire user details: ", err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message"	: "failed to retrieve user information from Google",
			"status"	: fiber.StatusBadGateway,
			"error"		: err.Error(),
		})
	}

	userData := entity.User{
		ID			: helper.HashEmail(strings.ToLower(googleUser.Email)),
		Fullname	: googleUser.Name,
		Email		: strings.ToLower(googleUser.Email),
		Password	: "",
		Photo		: googleUser.Picture,
		CreatedAt	: time.Now(),
		UpdatedAt	: time.Now() ,
	}

    client, err := database.FirestoreConnection()
    if err != nil {
        log.Fatal("error connecting to firestore", err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "error connecting to firestore",
            "status":  fiber.StatusBadRequest,
            "error":   err.Error(),
        })
    }
    defer client.Close()

    var allUsers []entity.User
    itr := client.Collection("skinsavvy-user").Documents(context.Background())
    for {
        doc, err := itr.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            log.Fatalf("Failed to iterate the list of users: %v", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "failed to get all users",
                "status":  fiber.StatusInternalServerError,
                "error":   err.Error(),
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
        allUsers = append(allUsers, user)
    }

    // Check if the user with the given email exists in the list
    var userExists bool
    var existingUser entity.User

	lowerEmail := strings.ToLower(googleUser.Email)
    for _, user := range allUsers {
        if strings.ToLower(user.Email) == lowerEmail {
            userExists = true
            existingUser = user
            break
        }
    }

    if userExists {
        _, err := client.Collection("skinsavvy-user").Doc(existingUser.ID).Set(context.Background(), userData)
        if err != nil {
            log.Fatal("failed updating user in Firestore: ", err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "failed updating user",
                "status":  fiber.StatusInternalServerError,
                "error":   err.Error(),
            })
        }
    } else {
        _, err := client.Collection("skinsavvy-user").Doc(userData.ID).Set(context.Background(), userData)
        if err != nil {
            log.Fatal("failed adding user to Firestore: ", err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "failed adding new user",
                "status":  fiber.StatusInternalServerError,
                "error":   err.Error(),
            })
        }
    }

	c.Cookie(&fiber.Cookie{
		Name	: "userId",
		Value	: userData.ID,
		Expires	: time.Now().Add(2 * time.Hour),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message"	: "successfully authenticated with Google account",
		"data"		: userData,
		"status"	: fiber.StatusOK,
	})
}