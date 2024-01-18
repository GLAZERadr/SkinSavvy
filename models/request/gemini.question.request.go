package request

type GeminiProductRecRequest struct {
	Question	string	`json:"question"`
	Gender		string	`json:"gender"`
	Age			int		`json:"age"`
	OutdoorAct	string	`json:"outdoor_activities"`
}