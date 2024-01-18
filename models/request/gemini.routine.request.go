package request 

type GeminiRoutineRecRequest struct {
	Products 		[]string	`json:"products"`
	TargetDays		int 		`json:"targetDays"`
	UserAge			int 		`json:"age"`
	UserSkinProblem	[]string	`json:"skinProblem"`
}