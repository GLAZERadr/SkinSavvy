package controllers

import (
	"log"
	"strings"

	"github.com/InnoFours/skin-savvy/llm"
	"github.com/InnoFours/skin-savvy/models/request"
	"github.com/InnoFours/skin-savvy/sephoraScrape"
	"github.com/gofiber/fiber/v2"
)

func SkincareRec(c *fiber.Ctx) error {
	var req request.GeminiRequest

	if err := c.BodyParser(&req); err != nil {
		log.Fatal("Failed to parse json body.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message"	: "failed to parse question to gemini",
			"status"	: fiber.StatusInternalServerError,
			"error"		: err.Error(),
		})
	}

	geminiResult, err := llm.GeminiClient(req.Question, req.Gender, req.Age, req.OutdoorAct)
	if err != nil {
		log.Fatal("Error processing question by gemini: ", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message"	: "can't process any question by gemini right now.",
			"status"	: fiber.StatusServiceUnavailable,
			"error"		: err.Error(),
		})
	}

	raw := geminiResult.Answer.Text
	items := strings.Split(raw, "\n")

	// Remove the empty string at the end (if any)
	if len(items) > 0 && items[len(items)-1] == "" {
		items = items[:len(items)-1]
	}

	// Extract only the item part (remove the number)
	var result []string
	for _, item := range items {
		parts := strings.SplitN(item, ". ", 2)
		if len(parts) == 2 {
			result = append(result, parts[1])
		}
	}

	details, err := sephoraScrape.ProductScraper(result)
	if err != nil {
		log.Fatal("Failed to scrape product", err.Error())
	}

	log.Println(result)

	return c.JSON(fiber.Map{
		"message"	: "successfully processed question by Gemini.",
		"status"	: fiber.StatusOK,
		"response"	: details,
	})
}
