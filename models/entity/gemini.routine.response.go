package entity

import "time"

type RoutineResponse struct {
	Date     		time.Time   `json:"date"`
	Product  		string `json:"product"`
}