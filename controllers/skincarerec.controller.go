package controllers

import (
	"log"

	"github.com/InnoFours/skin-savvy/llm"
	"github.com/InnoFours/skin-savvy/models/request"
	"github.com/gofiber/fiber/v2"
)

func SkincareRec(c * fiber.Ctx) error {
	var req request.GeminiRequest

	if err := c.BodyParser(&req); err != nil {
		log.Fatal("Failed to parse json body.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	: "failed to parse question to gemini",
			"status"	: fiber.StatusInternalServerError,
			"error"		: err.Error(),
		})
	}

	result, err := llm.GeminiClient(req.Question)
	if err != nil {
		log.Fatal("Error processing question by gemini: ", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message"	: "can't process any question by gemini right now.",
			"status"	: fiber.StatusServiceUnavailable,
			"error"		: err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":  "successfully processed question by Gemini.",
		"status":   fiber.StatusOK,
		"response": result,
	})
}