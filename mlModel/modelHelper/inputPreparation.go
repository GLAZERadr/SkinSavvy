package modelHelper

import (
	"log"
	"io"
	"image/jpeg"

	"github.com/nfnt/resize"
)

func InputPreparation(imageBuffer io.Reader) ([]float32, int64, int64){
	imageObject, err := jpeg.Decode(imageBuffer)
	if err != nil {
		log.Println("Cannot decode image")
	}

	imageSize := imageObject.Bounds().Size()
	imageWidth, imageHeight := int64(imageSize.X), int64(imageSize.Y)

	imageObject = resize.Resize(640, 640, imageObject, resize.Lanczos3)

	redChannel := []float32{}
	greenChannel := []float32{}
	blueChannel := []float32{}

	for y := 0; y < 640; y++ {
		for x := 0; x < 640; x++ {
			r, g, b, _ := imageObject.At(x, y).RGBA()
			redChannel = append(redChannel, float32(r/257)/255.0)
			greenChannel = append(greenChannel, float32(g/257)/255.0)
			blueChannel = append(blueChannel, float32(b/257)/255.0)
		}
	}

	inputArray := append(redChannel, greenChannel...)
	inputArray = append(inputArray, blueChannel...)

	return inputArray, imageWidth, imageHeight
}