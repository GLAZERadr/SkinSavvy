package middleware

import (
	"strings"
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

func TokenValidator(c *fiber.Ctx) error {

	firebaseAuth := c.Locals("firebaseAuth").(*auth.Client)

	authToken := c.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authToken, "Bearer", "", 1))

	log.Println("authToken: ", authToken)

	if idToken == "" {
		log.Println("Required authorization token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	: "Authorization token is included",
			"status"	: fiber.StatusUnauthorized,
		})
	}
	

	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	: "Invalid token!!!",
			"status"	: fiber.StatusUnauthorized,
			"error"		: err.Error(),
		})
	}

	c.Set("UUID", token.UID)
	return c.Next()
}
