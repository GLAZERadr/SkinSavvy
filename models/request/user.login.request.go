package request

type UserLogin struct {
	ID			string  `json:"id"`
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}