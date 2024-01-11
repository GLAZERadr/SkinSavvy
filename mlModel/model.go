package mlModel

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
)

func LoadModel(image io.Reader) (*fiber.Map, error) {
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)

	//read model files
	b, _ := os.ReadFile("D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/mlModel/onnxModel/skin_problem_detection_model.onnx")
	err := model.UnmarshalBinary(b)
	if err != nil {
		log.Fatal("Error opening model", err.Error())
	}

	//set the input
	err = model.SetInput(0, GetInput(image))
	if err != nil {
		log.Println("Not set properly")
	}

	backend.Run()

    // get output tensors
	model.GetOutputTensors()

	predictions := ProcessOutput(model.GetOutputTensors())
	return predictions, nil
}

