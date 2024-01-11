package mlModel

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"math"
	"sort"

	"github.com/InnoFours/skin-savvy/mlModel/images"
	"github.com/gofiber/fiber/v2"

	"github.com/kelseyhightower/envconfig"
	"github.com/nfnt/resize"
	"gorgonia.org/tensor"
	"gorgonia.org/tensor/native"
)

const (
	hSize, wSize = 640, 640
	blockSize     = 32
	gridHeight    = 13
	gridWidth     = 13
	boxesPerCell  = 5
	numClasses    = 6
	envConfPrefix = "yolo"
)

type configuration struct {
	ConfidenceThreshold float64	`envonfig:"confidence_threshold" default:"0.05" required:"true"`
	ClassProbaThreshold float64 `envconfig:"proba_threshold" default:"0.90" required:"true"`
}

func init() {
	err := envconfig.Process(envConfPrefix, &config)
	if err != nil {
		log.Fatal("Error processing env config")
	}
}

var (
	classes = []string{"Acne", "Blackhead", "Dark Circles", "Dry Skin", "Englarged Pores", "Wrinkles"}
	anchors     = []float64{1.08, 1.19, 3.42, 4.41, 6.63, 11.38, 9.42, 5.11, 16.62, 10.52}
	scaleFactor = float32(1) 
	config      configuration
)

func GetInput(imageReader io.Reader) tensor.Tensor {
	img, err := jpeg.Decode(imageReader)
	if err != nil {
		log.Fatal("Cannot decode image")
	}
	
	imgRescaled := image.NewRGBA(image.Rect(0,0,wSize,hSize))
	color := color.RGBA{0, 0, 0, 255}

	draw.Draw(imgRescaled, imgRescaled.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)

	var m image.Image
	if(img.Bounds().Max.X - img.Bounds().Min.X) > (img.Bounds().Max.Y - img.Bounds().Min.Y) {
		scaleFactor = float32(img.Bounds().Max.Y - img.Bounds().Min.Y) / float32(hSize)
		m = resize.Resize(0, hSize, img, resize.Lanczos3)
	} else {
		scaleFactor = float32(img.Bounds().Max.X - img.Bounds().Min.X) / float32(wSize)
		m = resize.Resize(0, wSize, img, resize.Lanczos3)
	}

	switch m.(type) {
	case *image.NRGBA:
		draw.Draw(imgRescaled, imgRescaled.Bounds(), m.(*image.NRGBA), image.ZP, draw.Src)
	case *image.YCbCr:
		draw.Draw(imgRescaled, imgRescaled.Bounds(), m.(*image.YCbCr), image.ZP, draw.Src)
	default:
		log.Fatal("Unhandled image type")
	}

	inputTensor := tensor.New(tensor.WithShape(1,3,hSize,wSize), tensor.Of(tensor.Float32))
	err = images.ImageToBCHW(imgRescaled, inputTensor)
	if err != nil {
		log.Fatal("Couldn't convert image to BCHW")
	}

	return inputTensor
}

type element struct {
	prob  float64
	class string
}

type box struct {
	r          image.Rectangle
	gridcell   []int
	confidence float64
	classes    []element
}

func sigmoid(sum float32) float32 {
	return float32(1.0 / (1.0 + math.Exp(float64(-sum))))
}

func sigmoid64(sum float32) float64 {
	return 1.0 / (1.0 + math.Exp(float64(-sum)))
}

func exp(val float32) float64 {
	return math.Exp(float64(val))
}

func softmax(a []float32) []float64 {
	var sum float64
	output := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		output[i] = math.Exp(float64(a[i]))
		sum += output[i]
	}
	for i := 0; i < len(output); i++ {
		output[i] = output[i] / sum
	}
	return output
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type byProba []element

func (b byProba) Len() int           { return len(b) }
func (b byProba) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byProba) Less(i, j int) bool { return b[i].prob < b[j].prob }

func getOrderedElements(input []float64, classes []string) []element {
	elems := make([]element, len(input))
	for i := 0; i < len(elems); i++ {
		elems[i] = element{
			prob:  input[i],
			class: classes[i],
		}
	}
	sort.Sort(sort.Reverse(byProba(elems)))
	return elems
}

func ProcessOutput(t []tensor.Tensor, err error) *fiber.Map {
	if err != nil {
		log.Fatal("an error has occurred: ", err.Error())
	}

	dense := t[0].(*tensor.Dense)

	dense.Reshape(125, gridHeight, gridWidth)

	data, err := native.Tensor3F32(dense)
	if err != nil {
		log.Fatal("can't convert to native tensor3f32: ", err.Error())
	}

	var boxes = make([]box, gridHeight*gridWidth*boxesPerCell)
	var counter int

	labelCounts := make(map[string]int)

	for cx := 0; cx < gridWidth; cx++ {
		for cy := 0; cy < gridHeight; cy++ {
			for b := 0; b < boxesPerCell; b++ {
				channel := b * (numClasses + 5)
				tx := data[channel][cx][cy]
				ty := data[channel+1][cx][cy]
				tw := data[channel+2][cx][cy]
				th := data[channel+3][cx][cy]
				tc := data[channel+4][cx][cy]
				tclasses := make([]float32, 6)

				for i := 0; i < 6; i++ {
					tclasses[i] = data[channel+5+i][cx][cy]
				}

				tclassesFloat64 := make([]float64, len(tclasses))
				for i, v := range tclasses {
					tclassesFloat64[i] = float64(v)
				}

				classes := getOrderedElements(tclassesFloat64, classes)
				topClass := classes[0].class

				labelCounts[topClass]++

				x := int((float32(cx) + sigmoid(tx)) * blockSize)
				y := int((float32(cy) + sigmoid(ty)) * blockSize)

				w := int(exp(tw) * anchors[2*b] * blockSize)
				h := int(exp(th) * anchors[2*b+1] * blockSize)

				boxes[counter] = box{
					gridcell:   []int{cx, cy},
					r:          image.Rect(max(y-w/2, 0), max(x-h/2, 0), min(y+w/2, wSize), min(x+h/2, hSize)),
					confidence: sigmoid64(tc),
					classes:    classes,
				}
				counter++
			}
		}
	}

	// Calculate percentage for each label
	totalDetections := float64(len(boxes) * boxesPerCell)
	var status []fiber.Map

	for label, count := range labelCounts {
		percentage := (float64(count) / totalDetections) * 100
		percentage = math.Round(percentage*100) / 100 // Round to two decimal places
		result := fiber.Map{
			"confidence level": config.ConfidenceThreshold,
			"label":            label,
			"status":       percentage,
		}
		status = append(status, result)
	}

	// Sort status by descending confidence
	sort.Slice(status, func(i, j int) bool {
		return status[i]["confidence level"].(float64) > status[j]["confidence level"].(float64)
	})

	return &fiber.Map{"predictions": status}
}
