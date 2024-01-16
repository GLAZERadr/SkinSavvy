package modelHelper

import (
	"log"
	"runtime"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

type ModelSession struct {
	Session *ort.AdvancedSession
	Input 	*ort.Tensor[float32]
	Output 	*ort.Tensor[float32]
}

var once sync.Once
var isInit bool

func InitSession(blank []float32) (ModelSession, error) {
	once.Do(func(){
		ort.SetSharedLibraryPath(getSharedLibPath())
		err := ort.InitializeEnvironment()
		if err != nil {
			log.Fatal("error initialize environment for onnx runtime: ", err.Error())
		}
		isInit = true
	})

	if !isInit {
		log.Fatal("error initialize environment for onnx runtime: The onnxruntime has already been initialized")
	}

	inputShape := ort.NewShape(1, 3, 640, 640)
	inputTensor, err := ort.NewTensor(inputShape, blank)

	if err != nil {
		log.Fatal("error creating input tensor: ", err.Error())
	}

	outputShape := ort.NewShape(1, 10, 8400)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		log.Fatal("error creating output tensor: ", err.Error())
	}

	options, err := ort.NewSessionOptions()
	if err != nil {
		log.Fatal("error options onnx runtime: ", err.Error())
	}

	session, err := ort.NewAdvancedSession("D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/mlModel/onnxModel/skin_problem_detection_model.onnx",
		[]string{"images"}, []string{"output0"}, 
		[]ort.ArbitraryTensor{inputTensor}, []ort.ArbitraryTensor{outputTensor}, options)
	if err != nil {
		log.Fatal("error initializing prediction process on init: ", err.Error())
	}
	
	modelSession := ModelSession{
		Session: session,
		Input: inputTensor,
		Output: outputTensor,
	}

	return modelSession, nil
}

func getSharedLibPath() string {
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			return "D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/mlModel/thirdParty/onnxruntime.dll"
		}
	}
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/mlModel/thirdParty/onnxruntime_arm64.dylib"
		}
	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return "D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/mlModel/thirdParty/onnxruntime_arm64.so"
		}
		return "D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/mlModel/thirdParty/onnxruntime.so"
	}
	panic("Unable to find a version of the onnxruntime library supporting this system.")
}