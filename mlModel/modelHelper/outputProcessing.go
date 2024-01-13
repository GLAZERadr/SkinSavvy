package modelHelper

import (
	"math"
	"sort"
)

var classes = []string{"Acne", "Blackhead", "Dark Circles", "Dry Skin", "Englarged Pores", "Wrinkles"}

func OutputProcessing(input []float32, imageWidth, imageHeight int64) [][]interface{} {
	boundBox := [][]interface{}{}

	for idx := 0; idx < 8400; idx++ {
		classID, probability := 0, float32(0.0)

		for col := 0; col < 6; col++ {
			currentProb := input[8400 * (col + 4) + idx]

			if currentProb > probability {
				probability = currentProb
				classID = col
			}
		}

		if probability < 0.01 {
			continue
		}

		label := classes[classID]

		xc, yc, w, h := input[idx], input[8400 + idx], input[2 * 8400 + idx], input[3 * 8400 + idx]
		x1 := (xc - w/2) / 640 * float32(imageWidth)
		y1 := (yc - h/2) / 640 * float32(imageHeight)
		x2 := (xc + w/2) / 640 * float32(imageWidth)
		y2 := (yc + h/2) / 640 * float32(imageHeight)

		boundBox = append(boundBox, []interface{}{float64(x1), float64(y1), float64(x2), float64(y2), label, probability})
	}

	sort.Slice(boundBox, func(i, j int) bool {
		return boundBox[i][5].(float32) < boundBox[j][5].(float32)
	})

	results := [][]interface{}{}

	for len(boundBox) > 0 {
		results = append(results, boundBox[0])
		tmp := [][]interface{}{}
		for _, box := range boundBox {
			if iou(boundBox[0], box) < 0.0 {
				tmp = append(tmp, box)
			}
		}

		boundBox = tmp
	}

	return results
}

func iou(box1, box2 []interface{}) float64 {
	// Calculate the area of intersection between the two bounding boxes using the intersection function
	intersectArea := intersection(box1, box2)

	// Calculate the union of the two bounding boxes using the union function
	unionArea := union(box1, box2)

	// The Intersection over Union (IoU) is the ratio of the intersection area to the union area
	return intersectArea / unionArea
}

func union(box1, box2 []interface{}) float64 {
	// Extract coordinates of the first rectangle
	rect1Left, rect1Bottom, rect1Right, rect1Top := box1[0].(float64), box1[1].(float64), box1[2].(float64), box1[3].(float64)

	// Extract coordinates of the second rectangle
	rect2Left, rect2Bottom, rect2Right, rect2Top := box2[0].(float64), box2[1].(float64), box2[2].(float64), box2[3].(float64)

	// Calculate area of the first rectangle
	rect1Area := (rect1Right - rect1Left) * (rect1Top - rect1Bottom)

	// Calculate area of the second rectangle
	rect2Area := (rect2Right - rect2Left) * (rect2Top - rect2Bottom)

	// Use the intersection function to calculate the area of overlap between the two rectangles
	intersectArea := intersection(box1, box2)

	// The union of two rectangles is the sum of their areas minus the area of their overlap
	return rect1Area + rect2Area - intersectArea
}

func intersection(box1, box2 []interface{}) float64 {
	// Extracting the coordinates of the first box
	firstBoxX1, firstBoxY1, firstBoxX2, firstBoxY2 := box1[0].(float64), box1[1].(float64), box1[2].(float64), box1[3].(float64)

	// Extracting the coordinates of the second box
	secondBoxX1, secondBoxY1, secondBoxX2, secondBoxY2 := box2[0].(float64), box2[1].(float64), box2[2].(float64), box2[3].(float64)

	// Calculating the x coordinate of the left side of the intersection
	intersectX1 := math.Max(firstBoxX1, secondBoxX1)

	// Calculating the y coordinate of the bottom side of the intersection
	intersectY1 := math.Max(firstBoxY1, secondBoxY1)

	// Calculating the x coordinate of the right side of the intersection
	intersectX2 := math.Min(firstBoxX2, secondBoxX2)

	// Calculating the y coordinate of the top side of the intersection
	intersectY2 := math.Min(firstBoxY2, secondBoxY2)

	// Calculating and returning the area of the intersection
	return (intersectX2 - intersectX1) * (intersectY2 - intersectY1)
}