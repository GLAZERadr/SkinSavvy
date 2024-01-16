package llm

import (
	"fmt"
	"log"

	"github.com/InnoFours/skin-savvy/config"
	gemini "github.com/Limit-LAB/go-gemini"
	"github.com/Limit-LAB/go-gemini/models"
)

type GeminiResult struct {
	Answer Answer `json:"Answer"`
}

type Answer struct {
	Text string `json:"text"`
}

func GeminiClient(promptCompletion string, gender string, age int, outdoor string) (*GeminiResult, error) {
	client := gemini.NewClient(config.ConfigGeminiKey())
	result, err := client.GenerateContent(models.GeminiPro,
		models.NewGenerateContentRequest( 
			models.NewContent(models.RoleUser, models.NewTextPart("You are a dermatologist. Now, you are gonna recommend skincare for your patience skin problem. You have to give 10 list of skincare from Sephora for "+promptCompletion+" problem. You also see other factors such as "+fmt.Sprint(age)+" and "+gender+" so you can determine what skincare product is suitable. You must determine wether your patience is doing outdoor activities commonly or not based on your patience answer."+outdoor+", the patience is regularly engange in outdoor activities. Also, note that you must give some explain of each product you recommend about why the ingredients contain in that product is suitable for your patience. There is no opening and closing statement, just straight to the answer that i request.")),
		),
	)

	if err != nil {
		log.Fatal("Failed connecting with gemini model: ", err.Error())
	}

	response := GeminiResult{
		Answer: Answer{
			Text: *result.Candidates[0].Content.Parts[0].Text,
		},
	}

	return &response, nil
}
