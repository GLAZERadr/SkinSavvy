package controllers

import (
	"log"
	"strings"

	// "time"

	"github.com/InnoFours/skin-savvy/google/llm"
	// "github.com/InnoFours/skin-savvy/models/entity"
	"github.com/InnoFours/skin-savvy/models/request"
	"github.com/gofiber/fiber/v2"
)

func SkincareRoutine(c *fiber.Ctx) error {
	var req request.GeminiRoutineRecRequest
	if err := c.BodyParser(&req); err != nil {
		log.Fatal("Failed to parse json body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to parse question to gemini",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	geminiResult, err := llm.GeminiRoutineRecommender(req.Products, req.TargetDays, req.UserAge, req.UserSkinProblem)
	if err != nil {
		log.Fatal("Error processing question by gemini: ", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "can't process any question by gemini right now.",
			"status":  fiber.StatusServiceUnavailable,
			"error":   err.Error(),
		})
	}

	raw := geminiResult.Answer.Text
	items := strings.Split(raw, "\n")

	log.Println(items)
	return c.JSON(fiber.Map{
		"message": "successfully processed question by Gemini.",
		"status":  fiber.StatusOK,
		"response": items,
	})
}
