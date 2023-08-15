package controllers

import (
	"bytes"
	"time"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/Data-Alchemist/doculex-api/config"
)

type S3Controller interface {
	UploadFile(c *fiber.Ctx) error
	// DeleteFile(c *fiber.Ctx) error
	// GetFilesList() ([]*s3.Object,error)
}

type s3Controller struct {}

func NewS3Controller() S3Controller {
	return &s3Controller{}
}

func (controller *s3Controller) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error track any file",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	sess, err := config.ConfigAWSSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error making aws session",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	svc := s3.New(sess)
	fileBytes := make([]byte, file.Size)
	if _, err := file.Open(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error read file",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	UploadedFileName := fmt.Sprintf("%s_%s", timeStamp, file.Filename)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket : aws.String("doculex-storage"),
		Key : aws.String(UploadedFileName),
		Body : bytes.NewReader(fileBytes),
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error upload file to bucket",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success upload file to bucket",
		"status": fiber.StatusOK,
		"file_name": UploadedFileName,
	})
}