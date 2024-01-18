package controllers

import (
	"context"
	"log"
	"time"

	"github.com/InnoFours/skin-savvy/database"
	"github.com/InnoFours/skin-savvy/helper"
	"github.com/InnoFours/skin-savvy/models/request"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
)

func AddUsageProduct(c *fiber.Ctx) error {
	var usedProduct request.UserProductUsage

	if err := c.BodyParser(&usedProduct); err != nil {
		log.Println("error parsing json body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	: "error parsing json body",
			"status"	: fiber.StatusBadRequest,
			"error"		: err.Error(),
		})
	}

	client, err := database.FirestoreConnection()
	if err != nil {
		log.Println("error connecting to firestore: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	: "error connecting to firestore",
			"status"	: fiber.StatusBadRequest,
			"error"		: err.Error(),
		})
	}
	defer client.Close()

	productUsedData := request.UserProductUsage {
		ID				: helper.HashProductUsage(usedProduct.ProductBrand, usedProduct.ProductName),
		UserID			: usedProduct.UserID,
		ProductBrand	: usedProduct.ProductBrand,
		ProductName		: usedProduct.ProductName,
		ProductImage	: usedProduct.ProductImage,
		CreatedAt		: time.Now(),
	}

	var allUsersUsedProduct []request.UserProductUsage
    itr := client.Collection("skinsavvy-product-usage").Documents(context.Background())
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

        user := request.UserProductUsage{
            ID				: doc.Data()["ID"].(string),
            UserID			: doc.Data()["UserID"].(string),
            ProductBrand	: doc.Data()["ProductBrand"].(string),
            ProductName		: doc.Data()["ProductName"].(string),
            ProductImage	: doc.Data()["ProductImage"].(string),
            CreatedAt		: doc.Data()["CreatedAt"].(time.Time),
        }
        allUsersUsedProduct = append(allUsersUsedProduct, user)
    }

	var productExists bool
	var existingUsedProduct request.UserProductUsage

	for _, productUsed := range allUsersUsedProduct {
		if (productUsed.UserID == usedProduct.UserID) && (productUsed.ProductBrand == usedProduct.ProductBrand) && (productUsed.ProductName == usedProduct.ProductName) {
			productExists = true
			existingUsedProduct = productUsed
		}
	}

	if productExists {
		_, err = client.Collection("skinsavvy-product-usage").Doc(existingUsedProduct.ID).Set(context.Background(), productUsedData)
		if err != nil {
			log.Println("failed adding user product usage to firestore: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed adding user product usage to firestore",
				"status" : fiber.StatusInternalServerError,
				"error"  : err.Error(),
			})
		}
	} else {
		_, err = client.Collection("skinsavvy-product-usage").Doc(productUsedData.ID).Set(context.Background(), productUsedData)
		if err != nil {
			log.Println("failed adding user product usage to firestore: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed adding user product usage to firestore",
				"status" : fiber.StatusInternalServerError,
				"error"  : err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "successfully added user product usage data",
		"data"   : productUsedData,
		"status" : fiber.StatusOK,
	})
}

func GetUserUsedProduct(c *fiber.Ctx) error {
	userId := c.Params("id")

	client, err := database.FirestoreConnection()
	if err != nil {
		log.Println("error connecting to firestore: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message"	: "error connecting to firestore",
			"status"	: fiber.StatusBadRequest,
			"error"		: err.Error(),
		})
	}
	defer client.Close()

	//search the user using iter
	var allUserUsedProduct []request.UserProductUsage
	itr := client.Collection("skinsavvy-product-usage").Documents(context.Background())
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

        user := request.UserProductUsage{
            ID				: doc.Data()["ID"].(string),
            UserID			: doc.Data()["UserID"].(string),
            ProductBrand	: doc.Data()["ProductBrand"].(string),
            ProductName		: doc.Data()["ProductName"].(string),
            ProductImage	: doc.Data()["ProductImage"].(string),
            CreatedAt		: doc.Data()["CreatedAt"].(time.Time),
        }
        allUserUsedProduct = append(allUserUsedProduct, user)
    }

	var userUsedProduct []request.UserProductUsage
	for _, user := range allUserUsedProduct {
		if userId == user.UserID {
			userUsedProduct = append(userUsedProduct, user)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message"		: "successfully retrieved all product used by user",
		"data"			: userUsedProduct,
		"totalProduct"	: len(userUsedProduct),
	})
}