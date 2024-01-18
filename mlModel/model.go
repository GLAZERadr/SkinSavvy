package mlModel

import (
	"bytes"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/InnoFours/skin-savvy/mlModel/modelHelper"
)

var blank []float32

type SkinProblemPercentages struct {
	Categories map[string]float32
	Total      float32
}

type SkinProblem struct {
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
}

func CalculatePercentages(spp SkinProblemPercentages) []SkinProblem {
	results := []SkinProblem{}
	var percentage int = 0

	for category, confidence := range spp.Categories {
		percentage = int((confidence / spp.Total) * 100)

		skinProblem := SkinProblem{
			Name:       category,
			Percentage: percentage,
		}
		results = append(results, skinProblem)
	}

	return results
}

func LoadModel(image io.Reader) (*fiber.Map, error) {
	imageData, err := io.ReadAll(image)
	if err != nil {
		log.Println("error reading image: ", err.Error())
	}

	imageBuffer := bytes.NewBuffer(imageData)

	reader := bytes.NewReader(imageBuffer.Bytes())

	blank, _, _ = modelHelper.InputPreparation(reader)

	modelSession, err := modelHelper.InitSession(blank)
	if err != nil {
		log.Println("error make prediction session: ", err.Error())
	}

	var resultsArray []fiber.Map

	for i := 0; i < 5; i++ {
		reader := bytes.NewReader(imageBuffer.Bytes())
		input, imageWidth, imageHeight := modelHelper.InputPreparation(reader)

		output, err := modelHelper.Inference(modelSession, input)
		if err != nil {
			log.Println("error processing inference in model session: ", err.Error())
		}

		boxes := modelHelper.OutputProcessing(output, imageWidth, imageHeight)
		for _, box := range boxes {
			objectName := box[4].(string)
			confidence := box[5].(float32)
			x1 := box[0].(float64)
			y1 := box[1].(float64)
			x2 := box[2].(float64)
			y2 := box[3].(float64)

			results := fiber.Map{
				"name":       objectName,
				"confidence": confidence,
				"coords": []fiber.Map{{"x1": x1, "y1": y1},
					{"x2": x2, "y2": y2}},
			}

			resultsArray = append(resultsArray, results)
		}
	}

	percentages := SkinProblemPercentages{
		Categories: make(map[string]float32),
		Total:      0,
	}

	for _, result := range resultsArray {
		name := result["name"].(string)
		confidence := result["confidence"].(float32)

		percentages.Categories[name] += confidence
		percentages.Total += confidence
	}

	summary := CalculatePercentages(percentages)

	return &fiber.Map{"summary": summary, "details": resultsArray}, nil
}
