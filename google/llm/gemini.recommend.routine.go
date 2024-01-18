package llm

import (
	"fmt"
	"log"
	"time"

	"github.com/InnoFours/skin-savvy/config"
	gemini "github.com/Limit-LAB/go-gemini"
	"github.com/Limit-LAB/go-gemini/models"
)

func GeminiRoutineRecommender(products []string, targetDays int, age int, skinProblem []string) (*GeminiResult, error) {
	client := gemini.NewClient(config.ConfigGeminiKey())
	result, err := client.GenerateContent(models.GeminiPro,
		models.NewGenerateContentRequest(
			models.NewContent(
				models.RoleUser,
				models.NewTextPart(`You are a dermatologist that gives advice about skincare routine to your patients. 
					You usually give advice about the usage of skincare from Sephora. 
					You might get requests about a bunch of products at one time. 
					Your task is to define the best skincare routine for ` +
					fmt.Sprint(targetDays) +
					` days for this list of skincare products from Sephora: ` +
					fmt.Sprint(products) +
					` for ` + fmt.Sprint(age) + ` years old patients with ` +
					fmt.Sprint(skinProblem) + ` problem on their faces. 
					You may also give the right timestamp like this format [Jan 18, 2024 at 5:07:34 PM UTC+8]. 
					Also, make sure that in the skincare routine you recommend, you do not give a routine for other than ` +
					fmt.Sprint(products) +
					`. You may not give other recommendations like meeting a dermatologist, eating a healthy diet, regular exercises, etc. 
					In the recommendation you give, don't skip even a single day from ` +
					fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")) +
					` until the targeted day that is ` +
					fmt.Sprint(time.Now().Add(time.Duration(targetDays*24)*time.Hour).Format("2006-01-02 15:04:05")) +
					`. Note that you must give the skincare routine for the patients and never say to consult a dermatologist or skincare professional for personalized advice tailored to your unique skin needs.
			
					Output format:
					date:product
					
					Note: do not generate format other than this reserved format, don't make an empty new line`)),
			
		),
	)

	if err != nil {
		log.Println("Failed connecting with gemini model: ", err.Error())
	}

	response := GeminiResult{
		Answer: Answer{
			Text: *result.Candidates[0].Content.Parts[0].Text,
		},
	}

	return &response, nil
}
