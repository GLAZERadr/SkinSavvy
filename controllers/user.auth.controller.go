package controllers

import (
	"context"
	"log"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/InnoFours/skin-savvy/auth"
	"github.com/InnoFours/skin-savvy/helper"
	"github.com/InnoFours/skin-savvy/models/entity"
	"github.com/InnoFours/skin-savvy/models/request"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct {
	authService *auth.AuthService
}

func NewAuthController(authService *auth.AuthService) *AuthController {
	return &AuthController{authService}
}

func(s *AuthController) UserRegister(c *fiber.Ctx) error {
	var register request.UserRegister

	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"failed to parse user login request",
			"status"	:	fiber.StatusBadRequest,
			"error"		:	err.Error(),
		})
	}

	var user entity.User

	if err := s.authService.DB.Raw("SELECT id, email, password FROM user WHERE email=?", register.Email).Scan(&user); err != nil {
		if s.authService.DB.RecordNotFound() {
			log.Fatal("Failed get user by email")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message"	: "Fa.iled to get user by email from database.",
				"status"	: fiber.StatusInternalServerError,
			})
		}
	}

	if user.ID != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":    "Email already exists.",
			"status":     fiber.StatusBadRequest,
		})
	}

	uid := uuid.New().String()

	passwordHashed, err := helper.PasswordHashing(register.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to hash password",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	if err := s.authService.DB.Exec("INSERT INTO users(id, email, fullname, age, password) VALUES (?,?,?,?,?)", uid, register.Email, register.Fullname, register.Age, passwordHashed).Error; err != nil {
		log.Fatal("Failed creating user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" 	: "Failed create new account.",
			"data"	    : register,
			"status"	: fiber.StatusInternalServerError,
		})
	}

	user = entity.User {
		ID       	: uid,
		Fullname 	: register.Fullname,
		Email     	: register.Email,
		Age			: register.Age,
		Password	: passwordHashed,
	}

	token, err := s.authService.FireAuth.CustomToken(context.Background(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to generate jwt token",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message"	:	"successfully create new user",
		"status"	:	fiber.StatusOK,
		"data"		:	user,
		"token"		:	token,
	})
}

func(s *AuthController) Login(c *fiber.Ctx) error {
	var login request.UserLogin

	if err := c.BodyParser(&login); err != nil {
		log.Fatal("Failed to parse json body.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"failed to parse user login request",
			"status"	:	fiber.StatusBadRequest,
			"error"		:	err.Error(),
		})
	}

	var user entity.User

	// s.authService.DB.LogMode(true)

	err := s.authService.DB.Where("email = ?", login.Email).First(&user).Error

	// s.authService.DB.LogMode(false)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Fatal("Error account not found:", err.Error())
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message"	: "user not found",
				"status"	: fiber.StatusNotFound,
			})
		}

		log.Fatal("Failed getting user account:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":  "Failed getting user account",
			"status":   fiber.StatusInternalServerError,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		log.Fatal("Bad request, invalid password.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"invalid password",
			"status"	:	fiber.StatusBadRequest,
		})
	}

	token, err := s.authService.FireAuth.CustomToken(context.Background(), user.ID)
	if err != nil {
		log.Fatal("Couldn't generate token.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to generate jwt token",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"success login",
		"status"	:	fiber.StatusOK,
		"data"		:	user,
		"token"		:	token,
	})
}