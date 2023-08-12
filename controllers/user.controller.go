package controllers

import (
	"time"
	"context"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"

	"github.com/Data-Alchemist/doculex-api/database"
	"github.com/Data-Alchemist/doculex-api/models/entity"
	"github.com/Data-Alchemist/doculex-api/models/request"
	"github.com/Data-Alchemist/doculex-api/local"
)


type UserController interface {
	//Get Request
	GetAllUserAccount(c *fiber.Ctx) error
	GetOneUserAccount(c *fiber.Ctx) error

	//Post Request
	CreateUserAccount(c *fiber.Ctx) error
	UserLoginValidator(c *fiber.Ctx) error

	//Put Request
	// UpdateUserInfo(c *fiber.Ctx) error

	// //Delete Request
	// DeleteUserAccount(c *fiber.Ctx) error
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
			"message"	:	"failed to get all user account",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	if err = cursor.All(context.Background(), &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get all user account",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"success get all user account",
		"status"	:	fiber.StatusOK,
		"data"		:	users,
	})
}

func(controller *userController) GetOneUserAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"failed to parse user id",
			"status"	:	fiber.StatusBadRequest,
			"error"		:	err.Error(),
		})
	}

	database.ConnectDB()
	defer database.DisconnectDB()

	client := database.GetDB()
	collection := database.GetCollection(client, "user")

	var user entity.User

	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message"	:	"user account not found",
				"status"	:	fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get user account",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"success get user account",
		"status"	:	fiber.StatusOK,
		"data"		:	user,
	})
}

func(controller *userController) CreateUserAccount(c *fiber.Ctx) error {
	var register request.UserRegister

	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"failed to parse user register request",
			"status"	:	fiber.StatusBadRequest,
			"error"		:	err.Error(),
		})
	}

	database.ConnectDB()
	defer database.DisconnectDB()

	client := database.GetDB()
	collection := database.GetCollection(client, "user")

	//check existing email in database
	userExist := collection.FindOne(context.Background(), bson.M{"email": request.Email})
	if userExist.Err() == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"email already exist",
			"status"	:	fiber.StatusBadRequest,
		})
	}

	//encrpyt the password
	passwordHashed, err := local.PasswordHashing(register.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to hash password",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	user := entity.User{
		Name		:	request.Name,
		Email		:	request.Email,
		Password	:	passwordHashed,
		CreatedAt	:	time.Now(),
		UpdatedAt	:	time.Now(),
	}
	user.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to create new user",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"successfully create new user",
		"status"	:	fiber.StatusOK,
		"data"		:	user,
	})
}

func(controller *userController) UserLoginValidator(c *fiber.Ctx) error {
	var login request.UserLogin

	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"failed to parse user login request",
			"status"	:	fiber.StatusBadRequest,
			"error"		:	err.Error(),
		})
	}

	database.ConnectDB()
	defer database.DisconnectDB()

	client := database.GetDB()
	collection := database.GetCollection(client, "user")

	var user entity.User

	err := collection.FindOne(context.Background(), bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message"	:	"user account not found",
				"status"	:	fiber.StatusNotFound,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get user account",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"invalid password",
			"status"	:	fiber.StatusBadRequest,
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"success login",
		"status"	:	fiber.StatusOK,
		"data"		:	user,
	})
}

// func(controller *userController) UpdateUserInfo(c *fiber.Ctx) error {}

// func(controller *userController) DeleteUserAccount(c *fiber.Ctx) error {}