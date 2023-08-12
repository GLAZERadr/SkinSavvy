package middleware

import (
	"time"
	"strings"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/Data-Alchemist/doculex-api/config"
)

func GenerateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.ConfigJWTkey()))
	if err != nil {
		return "", err
	}

	return tokenString,nil;
}

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.ConfigJWTkey()),
	})
}

//this function validates the jwt token from header
func JWTValidator(c *fiber.Ctx) (jwt.MapClaims, error) {
	tokenString := c.Get("Authorization")

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ConfigJWTkey()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}