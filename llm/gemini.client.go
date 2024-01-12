package llm

import (
	"log"

	"github.com/InnoFours/skin-savvy/config"
	gemini "github.com/Limit-LAB/go-gemini"
	"github.com/Limit-LAB/go-gemini/models"

	"github.com/gofiber/fiber/v2"
)

func GeminiClient(promptCompletion string) (*fiber.Map, error) {
	client := gemini.NewClient(config.ConfigGeminiKey())
	result, err := client.GenerateContent(models.GeminiPro, 
		models.NewGenerateContentRequest(
			models.NewContent(models.RoleUser, models.NewTextPart("Give 10 list of skincare for " + promptCompletion + " problem. Note that you just have to gift the name of skincare and you don't have to give an explanation. There is no opening and closing statement, just straight to the answer that i request.")),
		),
	)

	if err != nil {
		log.Fatal("Failed connecting with gemini model: ", err.Error())
	}

	response := fiber.Map{
		"Answer": result.Candidates[0].Content.Parts[0],
	}

	return &response, nil
}