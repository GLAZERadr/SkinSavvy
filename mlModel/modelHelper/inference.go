package modelHelper

import (
	"log"
)

func Inference(modelSession ModelSession, input []float32) ([]float32, error) {
	inputTensor := modelSession.Input.GetData()
	copy(inputTensor, input)
	err := modelSession.Session.Run()
	if err != nil {
		log.Println("error on inference process: ", err.Error())
	}

	return modelSession.Output.GetData(), nil
}