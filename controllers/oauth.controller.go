package controllers

import (
	"log"
	"strings"
	"time"

	"github.com/InnoFours/skin-savvy/database"
	"github.com/InnoFours/skin-savvy/middleware"
	"github.com/InnoFours/skin-savvy/models/entity"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
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

	log.Println("Fill user entity")
	userData := entity.User{
		ID			: uuid.New().String(),
		Fullname	: googleUser.Name,
		Email		: strings.ToLower(googleUser.Email),
		Password	: "",
		Photo		: googleUser.Picture,
		CreatedAt	: time.Now(),
		UpdatedAt	: time.Now(),
	}
	log.Println("user: ", userData)

	log.Println("insert to db")

	if err := database.DB.Model(&userData).Where("email = ?", strings.ToLower(googleUser.Email)).Updates(&userData).Error; err != nil {
		log.Fatalf("Error updating user data: %v", err)
	}

	if result := database.DB.Create(&userData); result.Error != nil {
		log.Fatalf("Error creating user data: %v", result.Error)
	}	

	log.Println("take user entity from db")
	var user entity.User
	database.DB.First(&user, "email = ?", strings.ToLower(googleUser.Email))

	c.Cookie(&fiber.Cookie{
		Name	: "userId",
		Value	: user.ID,
		Expires	: time.Now().Add(2 * time.Hour),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message"	: "successfully authenticated with Google account",
		"data"		: user,
		"status"	: fiber.StatusOK,
	})
}