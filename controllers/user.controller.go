package controllers

import (
	"fmt"
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
	"github.com/Data-Alchemist/doculex-api/middleware"
)


type UserController interface {
	//Get Request
	GetAllUserAccount(c *fiber.Ctx) error
	GetOneUserAccount(c *fiber.Ctx) error

	//Post Request
	CreateUserAccount(c *fiber.Ctx) error
	UserLoginValidator(c *fiber.Ctx) error

	//Put Request
	UpdateUserInfo(c *fiber.Ctx) error

	//Delete Request
	DeleteUserAccount(c *fiber.Ctx) error
}

type userController struct {}

func NewUserController() UserController {
	return &userController{}
}

func(controller *userController) GetAllUserAccount(c *fiber.Ctx) error {
	claims, err := middleware.JWTValidator(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	:	"failed to validate jwt token",
			"status"	:	fiber.StatusUnauthorized,
			"error"		:	err.Error(),
		})
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get userID from jwt token",
			"status"	:	fiber.StatusInternalServerError,
		})
	}

	fmt.Println(userID)

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
	claims, err := middleware.JWTValidator(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	:	"failed to validate jwt token",
			"status"	:	fiber.StatusUnauthorized,
			"error"		:	err.Error(),
		})
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get userID from jwt token",
			"status"	:	fiber.StatusInternalServerError,
		})
	}

	fmt.Println(userID)

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
	userExist := collection.FindOne(context.Background(), bson.M{"email": register.Email})
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
		Name		:	register.Name,
		Email		:	register.Email,
		Password	:	passwordHashed,
		CreatedAt	:	time.Now(),
		UpdatedAt	:	time.Now(),
	}
	user.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to create new user",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	token, err := middleware.GenerateJWTToken(user.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to generate jwt token",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"successfully create new user",
		"status"	:	fiber.StatusOK,
		"data"		:	user,
		"token"		:	token,
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

	err := collection.FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&user)
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"invalid password",
			"status"	:	fiber.StatusBadRequest,
		})
	}

	token, err := middleware.GenerateJWTToken(user.ID.Hex())
	if err != nil {
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

func(controller *userController) UpdateUserInfo(c *fiber.Ctx) error {
	claims, err := middleware.JWTValidator(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	:	"failed to validate jwt token",
			"status"	:	fiber.StatusUnauthorized,
			"error"		:	err.Error(),
		})
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get userID from jwt token",
			"status"	:	fiber.StatusInternalServerError,
		})
	}

	fmt.Println(userID)

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

	var user request.UserRegister

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	:	"failed to parse user update request",
			"status"	:	fiber.StatusBadRequest,
			"error"		:	err.Error(),
		})
	}

	update := bson.M{"$set":bson.M{}}

	if user.Name != "" {
		update["$set"].(bson.M)["name"] = user.Name
	}

	if user.Email != "" {
		update["$set"].(bson.M)["email"] = user.Email
	}

	if user.Password != "" {
		hashedPasswd, _:= bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		update["$set"].(bson.M)["password"] = string(hashedPasswd)
	}

	if len(update["$set"].(bson.M)) > 0 {
		update["$set"].(bson.M)["updated"] = time.Now()
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to update user account",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"success update user account",
		"status"	:	fiber.StatusOK,
	})
}

func(controller *userController) DeleteUserAccount(c *fiber.Ctx) error {
	claims, err := middleware.JWTValidator(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message"	:	"failed to validate jwt token",
			"status"	:	fiber.StatusUnauthorized,
			"error"		:	err.Error(),
		})
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to get userID from jwt token",
			"status"	:	fiber.StatusInternalServerError,
		})
	}

	fmt.Println(userID)

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

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	:	"failed to delete user account",
			"status"	:	fiber.StatusInternalServerError,
			"error"		:	err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message"	:	"user account not found",
			"status"	:	fiber.StatusNotFound,
		})
	}

	return c.JSON(fiber.Map{
		"message"	:	"success delete user account",
		"status"	:	fiber.StatusOK,
	})
}