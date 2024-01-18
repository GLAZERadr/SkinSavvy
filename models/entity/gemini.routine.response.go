package entity

type RoutineResponse struct {
	Day      		string   `json:"day"`
	Date     		string   `json:"date"`
	Product  		[]string `json:"product"`
	Explanation   	string   `json:"explanation"`
}