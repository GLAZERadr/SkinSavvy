package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/InnoFours/skin-savvy/auth"
	"github.com/gofiber/fiber/v2"
)

type FireAuthMiddleware struct {
	authService *auth.AuthService
}

func NewFireAuthMiddleware(authService *auth.AuthService) *FireAuthMiddleware {
	return &FireAuthMiddleware{authService}
}

func(s* FireAuthMiddleware) TokenValidator(c *fiber.Ctx) error {
	header := c.Get("Authorization")

	log.Println("header: ", header)

	if header == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	: "Authorization token is included",
			"status"	: fiber.StatusUnauthorized,
		})
	}

	parts := strings.Split(header, " ")
	log.Println("parts: ", parts)
	if len(parts) < 2 {
		// Handle the case where the token is missing
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid Authorization header",
			"status":  fiber.StatusUnauthorized,
		})
	}
	
	idToken := parts[1]
	log.Println("idToken: ", idToken)
	

	token, err := s.authService.FireAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	: "Invalid token!!!",
			"status"	: fiber.StatusUnauthorized,
			"error"		: err.Error(),
		})
	}

	c.Locals("token", token)
	return c.Next()
}