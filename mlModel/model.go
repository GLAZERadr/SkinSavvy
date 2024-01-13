package mlModel

import (
	"bytes"
	"io"
	"log"

	"github.com/8ff/prettyTimer"
	"github.com/gofiber/fiber/v2"

	"github.com/InnoFours/skin-savvy/mlModel/modelHelper"
)

var blank []float32 

func LoadModel(image io.Reader) (*fiber.Map, error) {
	ts := prettyTimer.NewTimingStats()

	imageData, err := io.ReadAll(image)
	if err != nil {
		log.Fatal("error reading image: ", err.Error())
	}

	imageBuffer := bytes.NewBuffer(imageData)

	reader := bytes.NewReader(imageBuffer.Bytes())

	blank, _, _ = modelHelper.InputPreparation(reader)

	modelSession, err := modelHelper.InitSession(blank)
	if err != nil {
		log.Fatal("error make prediction session: ", err.Error())
	}

	var resultsArray []fiber.Map

	for i := 0; i < 5; i++ {
		reader := bytes.NewReader(imageBuffer.Bytes())
		input, imageWidth, imageHeight := modelHelper.InputPreparation(reader)

		ts.Start()

		output, err := modelHelper.Inference(modelSession, input)
		if err != nil {
			log.Fatal("error processing inference in model session: ", err.Error())
		}

		ts.Finish()

		boxes := modelHelper.OutputProcessing(output, imageWidth, imageHeight)
		for _, box := range boxes {
			objectName := box[4].(string)
			confidence := box[5].(float32)
			x1 := box[0].(float64)
			y1 := box[1].(float64)
			x2 := box[2].(float64)
			y2 := box[3].(float64)

			results := fiber.Map{
				"name":    objectName,
				"confidence": confidence,
				"coords":   []fiber.Map{{"x1": x1, "y1": y1}, 
										{"x2": x2, "y2": y2}},
			}

			resultsArray = append(resultsArray, results)

			// log.Printf("Object: %s Confidence: %.2f Coordinates: (%f, %f), (%f, %f)", objectName, confidence, x1, y1, x2, y2)
		}
	}
	ts.PrintStats()

	return &fiber.Map{"predictions": resultsArray}, nil
}