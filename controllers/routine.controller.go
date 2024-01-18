package controllers

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/InnoFours/skin-savvy/google/llm"
	"github.com/InnoFours/skin-savvy/models/entity"
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

	log.Println("gemini result: ", geminiResult)
	raw := geminiResult.Answer.Text
	// Regular expression to extract date and product information
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}): (.+)`)

	// Split the input into lines
	lines := strings.Split(raw, "\n")

	// Iterate through each line and extract date and product
	var routines []entity.RoutineResponse
	for _, line := range lines {
		// Find the match in each line
		match := re.FindStringSubmatch(line)

		// Check if there is a match
		if len(match) == 3 {
			// Extract date and product
			dateStr := match[1]
			product := match[2]

			// Parse the date string
			date, err := time.Parse("2006-01-02 15:04:05", dateStr)
			if err != nil {
				log.Println("Error parsing date:", err)
				continue
			}

			routine := entity.RoutineResponse{
				Date	: date,
				Product	: product,
			}

			routines = append(routines, routine)
		}
	}


	return c.JSON(fiber.Map{
		"message"	: "successfully processed question by Gemini.",
		"status"	: fiber.StatusOK,
		"response"	: routines,
	})
}
